package inkwell

import (
	"fmt"
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
	UIManager        *UIManager
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

	State     *State
	Inventory *Inventory
	Audio     *AudioPlayer
	Camera    *Camera
	Input     *Input

	startID     string
	onStart     Action
	activeTheme string

	// runtime
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
	hoverLabel   string
	flash        string
	flashTimer   float64
	endCard      string
	speech       speechState
	dialog       *runtimeDialog
	activeDialog string

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

// LogKind classifies a chat-log entry — UI widgets style them differently.
type LogKind int

const (
	LogAction   LogKind = iota // player command ("> look at door")
	LogResponse                // in-world reply / Say output
	LogSystem                  // game/system message
)

// LogMessage is a single line in the chat-log buffer.
type LogMessage struct {
	Speaker string // empty for actions / system
	Text    string
	Kind    LogKind
}

// speechState is what the SpeechBubble widget renders. Mutated by Say.
type speechState struct {
	Active  bool
	Speaker string
	Text    string
}

// runtimeDialog holds the in-flight Dialogue state. Mutated by the
// startDialogue/gotoDialogueNode/endDialogue plumbing and read by the
// DialogBox widget.
type runtimeDialog struct {
	Dialogue   *Dialogue
	Node       *DialogueNode
	LineIdx    int
	ChoiceHits []Rectangle
}

// NewGame initializes a game with empty entity managers, the SCUMM-style
// verb set, all preset themes, and "classic-scumm" selected. Widgets are
// NOT auto-registered — call RegisterDefaultUI(g) explicitly.
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
		UIManager:        NewManager[Widget](),
		ThemeManager:     NewManager[Theme](),
		State:            NewState(),
		Inventory:        NewInventory(),
		Audio:            NewAudioPlayer(),
		Camera:           NewCamera(),
		Input:            newInput(),
		loaded:           newLoadedAssets(w, h),
		chars:            make(map[string]*runtimeChar),
		transition:       &transition{},
		selectedVerb:     "look",
	}
	g.Audio.attach(g.AssetManager)
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
func (g *Game) OnStart(a Action) *Game { g.onStart = a; return g }

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
func (g *Game) FlashLine(text string) {
	g.flash = text
	g.flashTimer = 2.0
}

// LogAction appends a player-command line (rendered with the prompt color
// by ChatLog widgets).
func (g *Game) LogAction(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogAction})
}

// LogResponse appends a reply line (rendered with the response color).
func (g *Game) LogResponse(speaker, text string) {
	g.appendLog(LogMessage{Speaker: speaker, Text: text, Kind: LogResponse})
}

// LogSystem appends a system/meta line (rendered with the response color).
func (g *Game) LogSystem(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogSystem})
}

func (g *Game) appendLog(m LogMessage) {
	g.messages = append(g.messages, m)
	if g.MaxLogLines > 0 && len(g.messages) > g.MaxLogLines {
		g.messages = g.messages[len(g.messages)-g.MaxLogLines:]
	}
}

// Messages returns the current log buffer (read-only).
func (g *Game) Messages() []LogMessage { return g.messages }

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
	s := g.SceneManager.MustGet(g.currentScene)
	for _, a := range s.Actors {
		if a.CharacterName == name {
			return true
		}
	}
	return false
}

// DrawText is a thin wrapper that lets widgets accept a Theme-supplied color
// today and switch to text/v2 later without changing call sites.
func (g *Game) DrawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	drawText(dst, s, x, y, c)
}

// HoverHint composes the "verb + target" string used by the StatusLine.
func (g *Game) HoverHint() string {
	verbLabel := ""
	if v, ok := g.VerbManager.Get(g.selectedVerb); ok {
		verbLabel = v.Label
	}
	if sel := g.Inventory.Selected(); sel != "" {
		if g.hoverLabel != "" {
			return verbLabel + " " + sel + " ezen: " + g.hoverLabel
		}
		return verbLabel + " " + sel
	}
	if g.hoverLabel == "" {
		return ""
	}
	return verbLabel + " " + g.hoverLabel
}

// ----- validation -------------------------------------------------------

func (g *Game) Validate() error {
	if !g.SceneManager.Has(g.startID) {
		return fmt.Errorf("%w: %q", ErrNoStartScene, g.startID)
	}
	for _, name := range g.SceneManager.Names() {
		s := g.SceneManager.MustGet(name)
		if s.Background == "" {
			return fmt.Errorf("%w: scene %q", ErrSceneMissingBackground, name)
		}
		if s.Background != "" && !g.AssetManager.Has(s.Background) {
			return fmt.Errorf("%w: scene %q background %q", ErrUnknownAsset, name, s.Background)
		}
		if s.Music != "" && !g.AssetManager.Has(s.Music) {
			return fmt.Errorf("%w: scene %q music %q", ErrUnknownAsset, name, s.Music)
		}
		for _, a := range s.Actors {
			if !g.CharacterManager.Has(a.CharacterName) {
				return fmt.Errorf("%w: scene %q actor %q", ErrUnknownCharacter, name, a.CharacterName)
			}
		}
		for _, e := range s.Exits {
			if e.To != "" && !g.SceneManager.Has(e.To) {
				return fmt.Errorf("%w: scene %q exit to %q", ErrUnknownScene, name, e.To)
			}
		}
	}
	if g.activeTheme == "" || !g.ThemeManager.Has(g.activeTheme) {
		return fmt.Errorf("inkwell: no active theme (got %q)", g.activeTheme)
	}
	return nil
}

// ----- runtime: characters & scenes -------------------------------------

type runtimeChar struct {
	def Character
	pos Point
	// path is the queue of remaining waypoints, head-first. Empty when
	// idle. The last entry equals target.
	path []Point
	// target is the final destination — kept separately so saves and the
	// post-Load resume in tickCharacters can rebuild a missing path.
	target Point
	moving bool

	// animation playback state — see actor.animation.go for the rules
	// that drive currentClip selection.
	currentClip string
	clipElapsed float64
	clipFrame   int
}

func (g *Game) makeCtx() *Ctx {
	c := &Ctx{Game: g, DT: 1.0 / 60.0}
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		c.Scene = &s
	}
	return c
}

func (g *Game) changeScene(name string) {
	if !g.SceneManager.Has(name) {
		logf("changeScene: unknown %q", name)
		return
	}
	prev := g.currentScene
	g.transition.start(func() {
		if prev != "" && prev != name {
			g.previousScene = prev
		}
		if prev != "" {
			old := g.SceneManager.MustGet(prev)
			if old.OnLeave != nil {
				g.queueAction(old.OnLeave, "OnLeave")
			}
		}
		g.currentScene = name
		g.State.NoteVisit(name)
		g.placeActors(name)
		g.resetTriggers()
		s := g.SceneManager.MustGet(name)
		if s.Music != "" {
			g.Audio.PlayMusic(s.Music)
		}
		if s.OnEnter != nil {
			g.queueAction(s.OnEnter, "OnEnter")
		}
	})
}

func (g *Game) placeActors(sceneName string) {
	s := g.SceneManager.MustGet(sceneName)
	for _, a := range s.Actors {
		def := g.CharacterManager.MustGet(a.CharacterName)
		rc, ok := g.chars[def.Name]
		if !ok {
			rc = &runtimeChar{def: def, pos: a.At}
			g.chars[def.Name] = rc
		} else {
			rc.def = def
			rc.pos = a.At
			rc.moving = false
		}
	}
}

func (g *Game) runtimeChar(name string) (*runtimeChar, bool) {
	c, ok := g.chars[name]
	return c, ok
}

func (g *Game) walkCharacter(name string, to Point) {
	c, ok := g.chars[name]
	if !ok {
		return
	}
	var boxes []Polygon
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		boxes = s.Walkboxes
	}
	path := pathfind(c.pos, to, boxes)
	if len(path) == 0 {
		c.moving = false
		return
	}
	c.path = path
	c.target = path[len(path)-1]
	c.moving = true
}

func (g *Game) characterMoving(name string) bool {
	c, ok := g.chars[name]
	return ok && c.moving
}

func (g *Game) tickCharacters(dt float64) {
	for _, c := range g.chars {
		g.tickAnimation(c, dt)
		if !c.moving {
			continue
		}
		// Resume after Load: moving with no live path means we lost the
		// route but still know the final destination; rebuild a trivial
		// straight-line path so the character finishes its walk.
		if len(c.path) == 0 {
			if c.pos == c.target {
				c.moving = false
				continue
			}
			c.path = []Point{c.target}
		}
		next := c.path[0]
		dx := next.X - c.pos.X
		dy := next.Y - c.pos.Y
		d := c.pos.Dist(next)
		speed := c.def.Speed
		if speed <= 0 {
			speed = 60
		}
		step := speed * dt
		if d <= step {
			c.pos = next
			c.path = c.path[1:]
			if len(c.path) == 0 {
				c.moving = false
			}
			continue
		}
		c.pos.X += dx / d * step
		c.pos.Y += dy / d * step
	}
}

// ----- script & dialog plumbing -----------------------------------------

func (g *Game) queueAction(a Action, label string) {
	if a == nil {
		return
	}
	if g.scriptRunner != nil {
		logf("queueAction: %s ignored, runner busy", label)
		return
	}
	g.scriptRunner = a.Start()
	g.scriptCtx = g.makeCtx()
	logf("queueAction %s", label)
}

func (g *Game) startDialogue(name string) {
	d, ok := g.DialogueManager.Get(name)
	if !ok {
		logf("startDialogue: unknown %q", name)
		return
	}
	startNode := d.Start
	if startNode == "" && len(d.Nodes) > 0 {
		startNode = d.Nodes[0].Name
	}
	node, ok := d.Node(startNode)
	if !ok {
		return
	}
	g.activeDialog = name
	g.dialog = &runtimeDialog{Dialogue: &d, Node: &node, LineIdx: 0}
}

func (g *Game) endDialogue() {
	g.activeDialog = ""
	g.dialog = nil
}

func (g *Game) gotoDialogueNode(name string) {
	if g.dialog == nil || g.activeDialog == "" {
		return
	}
	d := g.DialogueManager.MustGet(g.activeDialog)
	node, ok := d.Node(name)
	if !ok {
		logf("gotoNode: unknown %q in %q", name, g.activeDialog)
		return
	}
	g.dialog.Node = &node
	g.dialog.LineIdx = 0
	g.dialog.ChoiceHits = nil
}

func (g *Game) dialogueActive() bool { return g.dialog != nil }

func (g *Game) showEndCard(text string) { g.endCard = text }
