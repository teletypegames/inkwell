package inkwell

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// saveVersion is bumped when the on-disk schema changes. Load returns
// ErrSaveVersion if the file is from a future version.
const saveVersion = 1

// saveFile is the JSON-serialisable snapshot of a game session. Only the
// runtime mutable state lands here — managers, themes and assets are
// reconstructed by the domain's Build() on every launch.
type saveFile struct {
	Version       int                  `json:"version"`
	Title         string               `json:"title"`
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

// SaveDir overrides the on-disk directory used by Save/Load. When empty,
// the library falls back to "saves" relative to the working directory.
func (g *Game) saveDir() string {
	if g.SaveDir != "" {
		return g.SaveDir
	}
	return "saves"
}

func (g *Game) slotPath(slot int) string {
	return filepath.Join(g.saveDir(), fmt.Sprintf("slot%d.json", slot))
}

// Save serialises the current session into saves/slot<N>.json (or
// SaveDir/slot<N>.json if set). The in-flight script and dialogue are
// dropped — saves capture state at idle/ready boundaries.
func (g *Game) Save(slot int) error {
	if err := os.MkdirAll(g.saveDir(), 0o755); err != nil {
		return fmt.Errorf("inkwell: save mkdir: %w", err)
	}
	sf := g.buildSave()
	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return fmt.Errorf("inkwell: save marshal: %w", err)
	}
	if err := os.WriteFile(g.slotPath(slot), data, 0o644); err != nil {
		return fmt.Errorf("inkwell: save write: %w", err)
	}
	return nil
}

// Load replaces the runtime state with the contents of saves/slot<N>.json.
// Managers (items, scenes, characters, …) are left intact — Build() owns
// those. References in the save (scene name, theme, character names) are
// validated against the current managers before any state is touched.
func (g *Game) Load(slot int) error {
	data, err := os.ReadFile(g.slotPath(slot))
	if err != nil {
		return fmt.Errorf("inkwell: load read: %w", err)
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
	g.dialogue = nil
	g.activeDialogue = ""
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
			if _, ok := g.chars[a.CharacterName]; ok {
				continue
			}
			if !g.CharacterManager.Has(a.CharacterName) {
				continue
			}
			def := g.CharacterManager.MustGet(a.CharacterName)
			g.chars[a.CharacterName] = &runtimeChar{def: def, pos: a.At}
		}
		// Resume the scene's music; PlayMusic no-ops on the same track.
		if s.Music != "" {
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
