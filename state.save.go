package inkwell

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// saveVersion is bumped when the on-disk schema changes. Load returns
// ErrSaveVersion if the file is from a future version.
const saveVersion = 2

// saveFile is the JSON-serialisable snapshot of a game session. Only the
// runtime mutable state lands here — managers, themes and assets are
// reconstructed by the domain's Build() on every launch.
type saveFile struct {
	Version       int                  `json:"version"`
	Title         string               `json:"title"`
	SavedAt       time.Time            `json:"saved_at"`
	CurrentScene  string               `json:"current_scene"`
	PreviousScene string               `json:"previous_scene"`
	SelectedVerb  string               `json:"selected_verb"`
	ActiveTheme   string               `json:"active_theme"`
	Characters    map[string]savedChar `json:"characters"`
	Inventory     savedInventory       `json:"inventory"`
	State         savedState           `json:"state"`
}

type savedChar struct {
	Pos    Point `json:"pos"`
	Target Point `json:"target"`
	Moving bool  `json:"moving"`
}

type savedInventory struct {
	Items    []string `json:"items"`
	Selected string   `json:"selected"`
}

type savedState struct {
	Flags   map[string]bool `json:"flags"`
	Vars    map[string]any  `json:"vars"`
	Visited map[string]int  `json:"visited"`
	Talked  map[string]int  `json:"talked"`
}

const slotsVersion = 1

const slotsKey = "slots"

// slotsFile is the one document every save slot lives in: a slot is a key,
// "slot1", and its value carries what a load screen lists — when, where —
// in the clear, beside the session itself as a base64 string of the
// saveFile's JSON. The metadata is out here so that listing four slots
// does not mean decoding four sessions, and the session is opaque so that
// the document reads as a table of slots rather than as four games.
type slotsFile struct {
	Version int                  `json:"version"`
	Slots   map[string]slotEntry `json:"slots"`
}

type slotEntry struct {
	SavedAt time.Time `json:"saved_at"`
	Scene   string    `json:"scene"`
	Title   string    `json:"title,omitempty"`
	State   string    `json:"state"`
}

func slotKey(slot int) string { return fmt.Sprintf("slot%d", slot) }

func (g *Game) readSlots() (slotsFile, error) {
	doc := slotsFile{Version: slotsVersion, Slots: make(map[string]slotEntry)}
	data, err := g.store().Read(slotsKey)
	if errors.Is(err, ErrNotStored) {
		return doc, nil
	}
	if err != nil {
		return doc, err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return doc, fmt.Errorf("inkwell: slots unmarshal: %w", err)
	}
	if doc.Slots == nil {
		doc.Slots = make(map[string]slotEntry)
	}
	return doc, nil
}

func (g *Game) writeSlots(doc slotsFile) error {
	doc.Version = slotsVersion
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("inkwell: slots marshal: %w", err)
	}
	return g.store().Write(slotsKey, data)
}

// Save serialises the current session into slot N of the slots document.
// The in-flight script and dialog are dropped — saves capture state at
// idle/ready boundaries.
func (g *Game) Save(slot int) error {
	doc, err := g.readSlots()
	if err != nil {
		return err
	}
	sf := g.buildSave()
	data, err := json.Marshal(sf)
	if err != nil {
		return fmt.Errorf("inkwell: save marshal: %w", err)
	}
	doc.Slots[slotKey(slot)] = slotEntry{
		SavedAt: sf.SavedAt,
		Scene:   sf.CurrentScene,
		Title:   sf.Title,
		State:   base64.StdEncoding.EncodeToString(data),
	}
	return g.writeSlots(doc)
}

// Load replaces the runtime state with the session in slot N. Managers
// (items, scenes, characters, …) are left intact. References in the save
// (scene name, theme, character names) are validated against the current
// managers before any state is touched.
func (g *Game) Load(slot int) error {
	doc, err := g.readSlots()
	if err != nil {
		return err
	}
	entry, ok := doc.Slots[slotKey(slot)]
	if !ok {
		return fmt.Errorf("inkwell: load: slot %d is empty", slot)
	}
	data, err := base64.StdEncoding.DecodeString(entry.State)
	if err != nil {
		return fmt.Errorf("inkwell: load decode: %w", err)
	}
	var sf saveFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return fmt.Errorf("inkwell: load unmarshal: %w", err)
	}
	if sf.Version > saveVersion {
		return fmt.Errorf("inkwell: save version %d unsupported (max %d)", sf.Version, saveVersion)
	}
	return g.applySave(&sf)
}

func (g *Game) buildSave() saveFile {
	sf := saveFile{
		Version:       saveVersion,
		Title:         g.Title,
		SavedAt:       time.Now(),
		CurrentScene:  g.currentScene,
		PreviousScene: g.previousScene,
		SelectedVerb:  g.selectedVerb,
		ActiveTheme:   g.activeTheme,
		Characters:    make(map[string]savedChar, len(g.chars)),
		Inventory: savedInventory{
			Items:    g.Inventory.Items(),
			Selected: g.Inventory.Selected(),
		},
		State: savedState{
			Flags:   copyMap(g.State.flags),
			Vars:    copyMap(g.State.vars),
			Visited: copyMap(g.State.visited),
			Talked:  copyMap(g.State.talked),
		},
	}
	for name, c := range g.chars {
		sf.Characters[name] = savedChar{
			Pos:    c.pos,
			Target: c.target,
			Moving: c.moving,
		}
	}
	return sf
}

func (g *Game) applySave(sf *saveFile) error {
	// Validate references first so a corrupt save can't leave Game in a
	// half-restored state.
	if sf.CurrentScene != "" && !g.SceneManager.Has(sf.CurrentScene) {
		return fmt.Errorf("%w: %q (in save)", ErrUnknownScene, sf.CurrentScene)
	}
	if sf.ActiveTheme != "" && !g.ThemeManager.Has(sf.ActiveTheme) {
		return fmt.Errorf("inkwell: unknown theme %q in save", sf.ActiveTheme)
	}

	g.scriptRunner = nil
	g.scriptCtx = nil
	g.dialog = nil
	g.ClearSpeech()
	g.endCard = ""
	g.flash = ""
	g.flashTimer = 0

	g.currentScene = sf.CurrentScene
	g.previousScene = sf.PreviousScene
	if sf.SelectedVerb != "" {
		g.selectedVerb = sf.SelectedVerb
	}
	if sf.ActiveTheme != "" {
		g.activeTheme = sf.ActiveTheme
	}

	st := NewState()
	for k, v := range sf.State.Flags {
		if v {
			st.flags[k] = true
		}
	}
	for k, v := range sf.State.Vars {
		st.vars[k] = v
	}
	for k, v := range sf.State.Visited {
		st.visited[k] = v
	}
	for k, v := range sf.State.Talked {
		st.talked[k] = v
	}
	g.State = st

	inv := NewInventory()
	for _, it := range sf.Inventory.Items {
		inv.items = append(inv.items, it)
	}
	if sf.Inventory.Selected != "" && inv.Has(sf.Inventory.Selected) {
		inv.selected = sf.Inventory.Selected
	}
	g.Inventory = inv

	g.chars = make(map[string]*runtimeChar)
	for name, ch := range sf.Characters {
		if !g.CharacterManager.Has(name) {
			continue
		}
		def := g.CharacterManager.MustGet(name)
		g.chars[name] = &runtimeChar{
			def:    def,
			pos:    ch.Pos,
			target: ch.Target,
			moving: ch.Moving,
		}
	}
	// Backfill any actors required by the current scene that the save
	// missed (e.g. character added after the save was taken).
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		for _, a := range s.Actors {
			if _, ok := g.chars[a.Character.Name]; ok {
				continue
			}
			g.chars[a.Character.Name] = &runtimeChar{def: a.Character, pos: a.At}
		}
		// Resume the scene's music; PlayMusic no-ops on the same track.
		if s.Music != nil {
			g.Audio.PlayMusic(s.Music)
		}
	}

	g.resetTriggers()
	return nil
}

// copyMap snapshots a state map into the save file, so a later mutation of
// the live game cannot reach through into a marshalled save. An empty map is
// left out of the JSON entirely.
func copyMap[K comparable, V any](m map[K]V) map[K]V {
	if len(m) == 0 {
		return nil
	}
	out := make(map[K]V, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// SlotInfo describes one slot on disk, for a load or save screen to list.
// An empty slot answers too — a screen showing four slots wants four rows,
// three of which say "empty".
type SlotInfo struct {
	Slot    int
	Exists  bool
	Scene   string
	SavedAt time.Time
}

// Slots lists slots 1 to count from the document's metadata, without
// decoding a session. A slot nobody has written comes back as one that does
// not exist, because a screen listing saves is not the place to report a
// missing file.
func (g *Game) Slots(count int) []SlotInfo {
	doc, err := g.readSlots()
	if err != nil {
		logf("slots: %v", err)
	}
	slots := make([]SlotInfo, 0, count)
	for slot := 1; slot <= count; slot++ {
		info := SlotInfo{Slot: slot}
		if entry, ok := doc.Slots[slotKey(slot)]; ok {
			info.Exists = true
			info.Scene = entry.Scene
			info.SavedAt = entry.SavedAt
		}
		slots = append(slots, info)
	}
	return slots
}
