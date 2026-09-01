package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game is the root aggregate. It carries every Manager plus the runtime
// state. A domain package builds one via NewGame, registers entities,
// optionally calls RegisterDefaultUI / UseTheme, then hands it to Run.
type Game struct {
	Title         string
	Width, Height int

	ItemManager      *ItemManager
	SceneManager     *SceneManager
	CharacterManager *CharacterManager
	DialogueManager  *DialogueManager
	ScriptManager    *ScriptManager
	AssetManager     *AssetManager
	VerbManager      *VerbManager
	WidgetManager    *WidgetManager
	ThemeManager     *ThemeManager

	// SceneRect is the part of the window the picture occupies. Zero means
	// the whole window. A game with a HUD along the foot sets it, so that
	// generated exit strips (see scene.exit.go) land inside the painting.
	SceneRect Rectangle

	// ExitLook and ExitTake supply the default look and take responses for
	// every hotspot generated from Scene.Exits. Nil means no response.
	// An Exit can override either one.
	ExitLook func(Exit) Action
	ExitTake func(Exit) Action

	// UseWithFail supplies the response when the player uses an item on a
	// hotspot nobody paired it with. Nil falls back to a flash line.
	UseWithFail func(item string, h *Hotspot) Action

	// Player names the character a scene gets when it lists no actors of
	// its own, placed at that character's Start. Empty leaves such a
	// scene unpeopled.
	Player string

	// Walkboxes is the floor a scene walks on when it declares none of
	// its own. Nil means no routing in those scenes.
	Walkboxes []Polygon

	// WindowScale multiplies the internal resolution when Run sizes the
	// window. 0 means 4.
	WindowScale int

	State     *State
	Inventory *Inventory
	Audio     *AudioPlayer
	Input     *Input

	startID     string
	onStart     string
	onFinale    string
	activeTheme string

	// runtime
	queue         []Action
	loaded        *loadedAssets
	currentScene  string
	previousScene string
	exitHotspots  map[string][]Hotspot
	chars         map[string]*runtimeChar
	scriptRunner  Runner
	scriptCtx     *Ctx
	transition    *transition
	selectedVerb  string

	// UI runtime state (read by widgets, written by actions / engine)
	hoverLabel     string
	flash          string
	flashTimer     float64
	endCard        string
	speech         speechState
	dialogue       *runtimeDialogue
	activeDialogue string

	// in-game message log for ChatLog widgets; ring buffer behavior.
	messages    []LogMessage
	MaxLogLines int // 0 = unlimited (memory grows); set per-game.

	// SaveDir overrides the on-disk directory used by Save/Load.
	// Empty falls back to "saves" relative to the working directory.
	SaveDir string

	// triggerStates tracks per-trigger rising-edge / fired bookkeeping
	// for the current scene. Cleared on scene change and on Load.
	triggerStates map[string]*triggerState
}

// speechState is what the SpeechBubble widget renders. Mutated by Say.
type speechState struct {
	Active  bool
	Speaker string
	Text    string
}

// NewGame initializes a game with empty entity managers, the SCUMM-style
// verb set, all preset themes, and "classic-scumm" selected. Widgets are
// NOT auto-registered — call RegisterDefaultUI(g) explicitly.
//
// A domain is free to replace any of the entity managers with one of its
// own before Run — the registries are ordinary *Manager values, and Run
// wires the engine to whichever ones the Game holds by then.
func NewGame(title string, w, h int) *Game {
	g := &Game{
		Title:            title,
		Width:            w,
		Height:           h,
		ItemManager:      NewManager[Item](),
		SceneManager:     NewManager[Scene](),
		CharacterManager: NewManager[Character](),
		DialogueManager:  NewManager[Dialogue](),
		ScriptManager:    NewManager[Script](),
		AssetManager:     NewManager[Asset](),
		VerbManager:      NewManager[Verb](),
		WidgetManager:    NewManager[Widget](),
		ThemeManager:     NewManager[Theme](),
		State:            NewState(),
		Inventory:        NewInventory(),
		Audio:            NewAudioPlayer(),
		Input:            newInput(),
		loaded:           newLoadedAssets(w, h),
		chars:            make(map[string]*runtimeChar),
		transition:       &transition{},
		selectedVerb:     "look",
	}
	for _, v := range defaultVerbs() {
		g.VerbManager.Register(v)
	}
	RegisterPresetThemes(g)
	g.UseTheme("classic-scumm")
	return g
}

func (g *Game) StartAt(name string) *Game { g.startID = name; return g }

// CurrentScene is the scene the player is in, empty before the first one is
// entered. PreviousScene is the one before it, which is where Back() leads.
func (g *Game) CurrentScene() string  { return g.currentScene }
func (g *Game) PreviousScene() string { return g.previousScene }

// SceneArea is SceneRect, or the whole window when it was never set.
func (g *Game) SceneArea() Rectangle {
	if g.SceneRect.W > 0 && g.SceneRect.H > 0 {
		return g.SceneRect
	}
	return Rect(0, 0, float64(g.Width), float64(g.Height))
}

// SceneHotspots is everything clickable in a scene: its own hotspots first,
// then one per Exit. Authored hotspots come first because the engine takes the
// first area that contains the click, so a painted thing beats the edge strip
// an exit sits on wherever the two overlap.
//
// The expansion is cached, so the pointers HotspotAt hands out stay valid.
func (g *Game) SceneHotspots(name string) []Hotspot {
	if hs, ok := g.exitHotspots[name]; ok {
		return hs
	}
	s, ok := g.SceneManager.Get(name)
	if !ok {
		return nil
	}
	hs := make([]Hotspot, 0, len(s.Hotspots)+len(s.Exits))
	hs = append(hs, s.Hotspots...)
	for _, e := range s.Exits {
		hs = append(hs, e.hotspot(g))
	}
	if g.exitHotspots == nil {
		g.exitHotspots = make(map[string][]Hotspot)
	}
	g.exitHotspots[name] = hs
	return hs
}

// OnStart names the script queued after the start scene's OnEnter, once,
// when the game boots.
func (g *Game) OnStart(script string) *Game { g.onStart = script; return g }

// OnFinale names the game's closing script, queued directly after the start
// script — which is what a "boot straight into the ending" debug flag wants.
// Leave it unset for an ordinary run.
func (g *Game) OnFinale(script string) *Game { g.onFinale = script; return g }

// ----- theme + UI conveniences ------------------------------------------

func (g *Game) Theme() Theme             { return g.ThemeManager.MustGet(g.activeTheme) }
func (g *Game) UseTheme(name string)     { g.activeTheme = name }
func (g *Game) SelectedVerb() string     { return g.selectedVerb }
func (g *Game) SetSelectedVerb(s string) { g.selectedVerb = s }
func (g *Game) HoverLabel() string       { return g.hoverLabel }
func (g *Game) SetHoverLabel(s string)   { g.hoverLabel = s }

// SetSpeech / ClearSpeech are called by the Say action; widgets render
// whatever the current state says.
func (g *Game) SetSpeech(speaker, text string) {
	g.speech = speechState{Active: true, Speaker: speaker, Text: text}
}
func (g *Game) ClearSpeech() { g.speech = speechState{} }

// FlashLine shows a status-line message for ~2 seconds.
func (g *Game) FlashLine(text Text) {
	g.flash = text.String()
	g.flashTimer = 2.0
}

// HotspotAt returns the topmost hotspot in the current scene whose Area
// contains p, or nil if none. Exposed so widgets (e.g., a verb coin
// that only opens over hotspots) can do their own hit-tests in their
// Tick — the engine's own resolution runs only after every widget had a
// chance, so widgets can't rely on it.
func (g *Game) HotspotAt(p Point) *Hotspot {
	if g.currentScene == "" {
		return nil
	}
	hs := g.SceneHotspots(g.currentScene)
	for i := range hs {
		h := &hs[i]
		if h.Area != nil && h.Area.Contains(p) {
			return h
		}
	}
	return nil
}

// CharacterInScene reports whether a registered character is listed as
// an actor in the current scene. Used by the CharacterPanel widget so
// it auto-hides when the character isn't in view.
func (g *Game) CharacterInScene(name string) bool {
	if g.currentScene == "" {
		return false
	}
	for _, a := range g.sceneActors(g.SceneManager.MustGet(g.currentScene)) {
		if a.CharacterName == name {
			return true
		}
	}
	return false
}

// VoiceOf reports how a character reaches the player. An unregistered
// name speaks in a bubble, which is what a bare Say always did.
func (g *Game) VoiceOf(name string) Voice {
	if c, ok := g.CharacterManager.Get(name); ok {
		return c.Voice
	}
	return VoiceBubble
}

// CharacterLabel is what the UI should print for a character: the Label
// if the entity carries one, the name otherwise.
func (g *Game) CharacterLabel(name string) Text {
	if c, ok := g.CharacterManager.Get(name); ok {
		return c.DisplayName()
	}
	return plain(name)
}

// ItemLabel is what the UI should print for an item: its Description if it
// carries one, its name otherwise. The status line and the chat log both
// name the item the player is holding, and a bare registry key is not a
// thing anybody wrote — least of all in the player's language.
func (g *Game) ItemLabel(name string) Text {
	if it, ok := g.ItemManager.Get(name); ok && !it.Description.Empty() {
		return it.Description
	}
	return plain(name)
}

// DrawText is a thin wrapper that lets widgets accept a Theme-supplied color
// today and switch to text/v2 later without changing call sites.
func (g *Game) DrawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	drawText(dst, s, x, y, c)
}

// HoverHint composes the "verb + target" string used by the StatusLine.
func (g *Game) HoverHint() string {
	var verb Text
	if v, ok := g.VerbManager.Get(g.selectedVerb); ok {
		verb = v.Label
	}
	if sel := g.Inventory.Selected(); sel != "" {
		item := g.ItemLabel(sel)
		if g.hoverLabel != "" {
			return Format(HintVerbOn, verb, item, g.hoverLabel).String()
		}
		return Format(HintVerb, verb, item).String()
	}
	if g.hoverLabel == "" {
		return ""
	}
	return Format(HintVerb, verb, g.hoverLabel).String()
}

func (g *Game) makeCtx() *Ctx {
	c := &Ctx{Game: g, DT: 1.0 / 60.0}
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		c.Scene = &s
	}
	return c
}

func (g *Game) showEndCard(text string) { g.endCard = text }
