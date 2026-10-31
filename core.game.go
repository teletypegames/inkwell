package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase/input"
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
	SpeakerManager   *SpeakerManager
	DialogManager    *DialogManager
	ScriptManager    *ScriptManager
	AssetManager     *AssetManager
	VerbManager      *VerbManager
	WidgetManager    *WidgetManager
	ThemeManager     *ThemeManager
	StageManager     *StageManager
	// AchievementManager holds what the player can be recognised for. Its
	// entries are evaluated on the world stage and reported through the
	// [Extension]s — see [Achievement] and [Extension].
	AchievementManager *AchievementManager

	// Extensions live beside the game for the whole run: started before the
	// settings are read, ticked first thing every frame, stopped when the
	// loop returns. A platform layer such as Steam is one. Nil is a game
	// with none.
	Extensions []Extension

	// SceneRect is the part of the window the picture occupies. Zero means
	// the whole window. A game with a HUD along the foot sets it, so that
	// generated exit strips (see [Exit]) land inside the painting.
	SceneRect Rectangle

	// ExitLook and ExitTake supply the default look and take responses for
	// every hotspot generated from Scene.Exits. Nil means no response.
	// An Exit can override either one.
	ExitLook func(Exit) Action
	ExitTake func(Exit) Action

	// UseWithFail supplies the response when the player uses an item on a
	// hotspot nobody paired it with. Nil falls back to a flash line.
	UseWithFail func(item *Item, h *Hotspot) Action

	// ----- what the game boots as -------------------------------------
	//
	// These four say where the game opens, not what it does afterwards.
	// New reads them; StartAt, OnStart, OnFinale and UseTheme set the
	// same values for a caller that would rather build a game in steps.

	// StartStage is the screen the game opens on — its main menu, or the
	// world itself for a run that goes straight in.
	StartStage *Stage
	// StartScene is the scene a session opens in. It is entered when the
	// world starts, which is at boot for a game whose StartStage is the
	// world, and on Restart for one that opens on a menu.
	StartScene *Scene
	// StartScript is queued after that scene's OnEnter, once.
	StartScript *Script
	// FinaleScript is the closing script, queued straight after
	// StartScript — which is what a "boot into the ending" flag wants.
	// Empty for an ordinary run.
	FinaleScript *Script
	// StartTheme is the theme selected at boot. Empty takes
	// "classic-scumm", so a game that never mentions a theme has one.
	StartTheme *Theme
	// Flags are set in State before the first frame. A game whose debug
	// build or command line switches something on says so here rather
	// than reaching into State after construction.
	Flags []string

	// SystemLang is the language the machine or the platform reports — the
	// Steam client's, an OS locale — and the one a first run speaks when
	// nothing stored says otherwise. Run's order is: the language set
	// before Run, then the stored setting, then this, then Fallback. An
	// Extension's Start is where a platform writes it.
	SystemLang string

	// Player names the character the player controls: the one a scene gets
	// when it lists no actors of its own, placed at that character's Start,
	// and the one WalkOnClick moves. Empty leaves such a scene unpeopled.
	Player *Character

	// WalkOnClick makes a left click on empty floor walk Player there,
	// routed through the walkboxes. Off by default: a game whose hotspots
	// answer from any distance has no reason for its character to move,
	// and one that does says so here.
	WalkOnClick bool

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

	activeTheme   string
	currentStage  *Stage
	previousStage *Stage
	started       bool
	quitting      bool
	pending       pendingChange

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
	hoverLabel string
	flash      string
	flashTimer float64
	endCard    string
	speech     speechState
	talking    string
	clock      float64
	dialog     *runtimeDialog

	// in-game message log for ChatLog widgets; ring buffer behavior.
	messages    []LogMessage
	MaxLogLines int // 0 = unlimited (memory grows); set per-game.

	// App names the game to the machine it runs on: the directory its
	// settings and saves live in natively, the localStorage prefix in a
	// browser. Empty is "inkwell".
	App string
	// Store is where settings and saves are kept. Nil takes the platform's
	// default for App; a test hands in a MemoryStore.
	Store Store

	settings Settings

	// triggerStates tracks per-trigger rising-edge / fired bookkeeping
	// for the current scene. Cleared on scene change and on Load.
	triggerStates map[string]*triggerState

	// unlocked is every Achievement reported to the Extensions this process,
	// so a condition that stays true is handed over once.
	unlocked map[string]bool
}

// speechState is what the SpeechBubble widget renders. Mutated by Say.
type speechState struct {
	Active  bool
	Speaker string
	Text    string
}

// New builds a game from a declared one. Everything a game is rather than
// does — its size, its managers, where it opens, the responses it hands the
// engine — is a field, so the whole of it can be one literal at the call
// site instead of a page of assignments the reader has to hold in their
// head:
//
//	game := inkwell.New(inkwell.Game{
//		Title:        "Real World",
//		Width:        640,
//		Height:       380,
//		StartScene:   "alley",
//		StartScript:  "opening",
//		Player:       "paul",
//		SceneManager: scene.SceneManager,
//		…
//	})
//
// What the spec leaves out, New fills in: an empty manager for every
// category the caller did not hand over, the runtime halves (State,
// Inventory, Audio, Input), the SCUMM verb set when no verbs were
// registered, and the preset themes — added to whatever the caller's
// ThemeManager already holds rather than replacing them. Widgets are NOT
// auto-registered; Run falls back to the SCUMM preset when none were.
//
// A manager handed over here is used in place, not copied, so the engine
// and the domain share one registry per category instead of two to keep
// in step.
func New(spec Game) *Game {
	g := &spec
	if g.ItemManager == nil {
		g.ItemManager = NewManager[*Item]()
	}
	if g.SceneManager == nil {
		g.SceneManager = NewManager[*Scene]()
	}
	if g.CharacterManager == nil {
		g.CharacterManager = NewManager[*Character]()
	}
	if g.SpeakerManager == nil {
		g.SpeakerManager = NewManager[*Speaker]()
	}
	if g.DialogManager == nil {
		g.DialogManager = NewManager[*Dialog]()
	}
	if g.ScriptManager == nil {
		g.ScriptManager = NewManager[*Script]()
	}
	if g.AssetManager == nil {
		g.AssetManager = NewManager[*Asset]()
	}
	if g.VerbManager == nil {
		g.VerbManager = NewManager[*Verb]()
	}
	if g.WidgetManager == nil {
		g.WidgetManager = NewManager[Widget]()
	}
	if g.ThemeManager == nil {
		g.ThemeManager = NewManager[*Theme]()
	}
	if g.StageManager == nil {
		g.StageManager = NewManager[*Stage]()
	}
	if g.AchievementManager == nil {
		g.AchievementManager = NewManager[*Achievement]()
	}

	g.State = NewState()
	g.Inventory = NewInventory()
	g.Audio = NewAudioPlayer()
	g.Input = input.New()
	g.loaded = newLoadedAssets(g.Width, g.Height)
	g.chars = make(map[string]*runtimeChar)
	g.unlocked = make(map[string]bool)
	g.transition = &transition{}
	g.selectedVerb = "look"
	g.pending = noChange
	g.settings.Voice = true

	if g.VerbManager.Len() == 0 {
		for _, v := range defaultVerbs() {
			g.VerbManager.Register(v)
		}
	}
	RegisterPresetThemes(g)

	g.activeTheme = "classic-scumm"
	if g.StartTheme != nil {
		g.activeTheme = g.StartTheme.Name
	}
	for _, f := range g.Flags {
		g.State.SetFlag(f)
	}
	return g
}

// NewGame is New for a game that would rather be built in steps: it
// declares nothing but the window, and the caller assigns the rest.
func NewGame(title string, w, h int) *Game {
	return New(Game{
		Title:  title,
		Width:  w,
		Height: h,
	})
}

// StartAt sets StartScene, the scene a new session opens in, and hands the
// game back for chaining.
func (g *Game) StartAt(scene *Scene) *Game { g.StartScene = scene; return g }

// CurrentScene is the name of the scene the player is in, empty before the
// first one is entered.
func (g *Game) CurrentScene() string { return g.currentScene }

// PreviousScene is the name of the scene before the current one, which is
// where Back leads; empty if there was none.
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
func (g *Game) OnStart(script *Script) *Game { g.StartScript = script; return g }

// OnFinale names the game's closing script, queued directly after the start
// script — which is what a "boot straight into the ending" debug flag wants.
// Leave it unset for an ordinary run.
func (g *Game) OnFinale(script *Script) *Game { g.FinaleScript = script; return g }

// ----- theme + UI conveniences ------------------------------------------

// Theme returns the active Theme. When the active name is not registered —
// a settings file or a save written by another build of the game — it
// falls back to classic-scumm, then to the first registered theme. It
// panics only when the ThemeManager is empty, which [Game.Validate]
// rejects at start.
func (g *Game) Theme() *Theme {
	if t, ok := g.ThemeManager.Get(g.activeTheme); ok {
		return t
	}
	if t, ok := g.ThemeManager.Get("classic-scumm"); ok {
		return t
	}
	if names := g.ThemeManager.Names(); len(names) > 0 {
		return g.ThemeManager.MustGet(names[0])
	}
	panic("inkwell: no theme registered")
}

// UseTheme makes theme the active one, by its Name; it is not saved to the
// settings (the UseTheme action does that). A nil theme, or one not
// registered in the ThemeManager, leaves the active theme as it is and
// reports false.
func (g *Game) UseTheme(theme *Theme) bool {
	if theme == nil || !g.ThemeManager.Has(theme.Name) {
		return false
	}
	g.activeTheme = theme.Name
	return true
}

// SelectedVerb returns the name of the verb the player has picked, or ""
// for none.
func (g *Game) SelectedVerb() string { return g.selectedVerb }

// SetSelectedVerb picks the verb with the name s; "" picks none.
func (g *Game) SetSelectedVerb(s string) { g.selectedVerb = s }

// HoverLabel returns the label of what the pointer is over, as the engine
// set it this frame, for a widget to show.
func (g *Game) HoverLabel() string { return g.hoverLabel }

// SetHoverLabel sets the label HoverLabel returns.
func (g *Game) SetHoverLabel(s string) { g.hoverLabel = s }

// SetSpeech puts up text as the line speaker is saying; the Say action
// calls it, and the speech widgets draw whatever it holds.
func (g *Game) SetSpeech(speaker, text string) {
	g.speech = speechState{Active: true, Speaker: speaker, Text: text}
}

// ClearSpeech takes the line SetSpeech put up down again.
func (g *Game) ClearSpeech() { g.speech = speechState{} }

// itemFrame is the frame an item's sprite shows now: items animate on the
// game clock, because an inventory slot keeps no state of its own.
func itemFrame(g *Game, s *Sprite) int {
	if s == nil {
		return 0
	}
	return s.FrameAt(g.clock)
}

// FlashLine puts up a status-line message for about two seconds of game
// time, replacing the one up now. The game loop counts it down every frame
// and clears it, so it expires whether or not a [StatusLine] is showing it;
// [Game.FlashText] reads it.
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
func (g *Game) CharacterInScene(character *Character) bool {
	if g.currentScene == "" || character == nil {
		return false
	}
	for _, a := range g.sceneActors(g.SceneManager.MustGet(g.currentScene)) {
		if a.Character == character {
			return true
		}
	}
	return false
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
	return &Ctx{Game: g, DT: 1.0 / 60.0, State: g.State}
}

func (g *Game) showEndCard(text string) { g.endCard = text }

// ----- stages -----------------------------------------------------------

// pendingChange is a screen change asked for by an action and applied at the
// top of the next frame. Replacing State, the inventory and every character
// underneath the action that asked for it is not something that action should
// have to survive, so nothing here happens where it is called.
type pendingChange struct {
	stage   *Stage
	restart bool
	load    int
}

var noChange = pendingChange{load: -1}

// CurrentStage is the screen showing, nil before the first frame.
func (g *Game) CurrentStage() *Stage { return g.currentStage }

// PreviousStage is the screen shown before the current one, which is where
// a "back" row leads ([GoBack]); nil if there was none.
func (g *Game) PreviousStage() *Stage { return g.previousStage }

// Started reports whether a session exists: the world has been begun or
// loaded. A "Resume" or "Save" row is off until it does.
func (g *Game) Started() bool { return g.started }

// OnWorldStage reports whether the screen showing is the game. Off it the
// scene is not painted, characters do not walk and triggers do not fire.
func (g *Game) OnWorldStage() bool {
	return g.currentStage != nil && g.currentStage.World
}

// GoToStage shows another screen, at the top of the next frame.
func (g *Game) GoToStage(stage *Stage) { g.pending.stage = stage }

// Restart begins a new session, at the top of the next frame.
func (g *Game) Restart() { g.pending.restart = true }

// LoadLater restores a saved session, at the top of the next frame. Load
// itself is immediate, and is what this ends up calling.
func (g *Game) LoadLater(slot int) { g.pending.load = slot }

// Quit closes the window at the end of the frame.
func (g *Game) Quit() { g.quitting = true }

// applyPending is the one place a screen changes. It runs before anything
// else in a frame, so the world it tears down is not one anybody is standing
// in the middle of. A Stage that leaves through a Transition defers the
// whole change to the covered frame — the restart or load asked for in the
// same breath included, so the menu still showing during the fade is a menu
// over the world it was over.
func (g *Game) applyPending() {
	change := g.pending
	if change == noChange {
		return
	}
	g.pending = noChange
	g.settings.Voice = true

	if spec := g.leavingTransition(change.stage); spec != nil {
		g.transition.start(spec, func() { g.applyChange(change) })
		return
	}
	g.applyChange(change)
}

// leavingTransition is the Transition the current Stage goes out through,
// when the change is to another Stage and the Stage names one.
func (g *Game) leavingTransition(to *Stage) Transition {
	if to == nil || to == g.currentStage || g.currentStage == nil {
		return nil
	}
	return g.currentStage.Transition
}

func (g *Game) applyChange(change pendingChange) {
	if change.restart || change.load >= 0 {
		g.queue = nil
		g.scriptRunner = nil
		g.scriptCtx = nil
	}
	switch {
	case change.load >= 0:
		if err := g.Load(change.load); err != nil {
			logf("load slot %d: %v", change.load, err)
			break
		}
		g.started = true
	case change.restart:
		g.startWorld()
	}
	if change.stage != nil {
		g.enterStage(change.stage)
	}
}

// enterStage runs the outgoing screen's OnLeave and the incoming one's
// OnEnter, both through the queue, so either may take as many frames as it
// likes.
func (g *Game) enterStage(stage *Stage) {
	if stage == g.currentStage {
		return
	}
	if previous := g.currentStage; previous != nil && previous.OnLeave != nil {
		g.queueAction(previous.OnLeave, "stage leave "+previous.Name)
	}
	g.previousStage = g.currentStage
	g.currentStage = stage
	g.arriveStage(stage)
}

// arriveStage queues what a Stage does once it is showing: its OnEnter, and
// after that its Hold, on the same queue, so the clock starts when OnEnter
// is done and the move lands like any other change, at the top of a frame
// and through the Stage's Transition.
func (g *Game) arriveStage(stage *Stage) {
	if stage == nil {
		return
	}
	if stage.Music != nil {
		g.Audio.PlayMusic(stage.Music)
	}
	if stage.OnEnter != nil {
		g.queueAction(stage.OnEnter, "stage enter "+stage.Name)
	}
	if stage.Hold > 0 && stage.Next != nil {
		g.queueAction(&holdAction{stage: stage}, "stage hold "+stage.Name)
	}
}

// startWorld begins a session: everything a save would have replaced, put
// back the way a new game leaves it, and the start scene entered.
func (g *Game) startWorld() {
	g.State = NewState()
	for _, flag := range g.Flags {
		g.State.SetFlag(flag)
	}
	g.Inventory = NewInventory()
	g.chars = make(map[string]*runtimeChar)
	g.dialog = nil
	g.ClearSpeech()
	g.endCard = ""
	g.flash = ""
	g.flashTimer = 0
	g.selectedVerb = "look"
	g.previousScene = ""
	g.currentScene = g.StartScene.Name
	g.State.NoteVisit(g.currentScene)
	g.placeActors(g.currentScene)
	g.resetTriggers()
	g.started = true

	scene := g.StartScene
	if scene.Music != nil {
		g.Audio.PlayMusic(scene.Music)
	}
	var opening []Action
	if scene.OnEnter != nil {
		opening = append(opening, scene.OnEnter)
	}
	if g.StartScript != nil {
		opening = append(opening, RunScript(g.StartScript))
	}
	if g.FinaleScript != nil {
		opening = append(opening, RunScript(g.FinaleScript))
	}
	if len(opening) > 0 {
		g.queueAction(Seq(opening...), "world start")
	}
}
