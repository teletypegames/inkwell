package pncdsl

import (
	"fmt"
)

// Game is the root aggregate. It carries every Manager plus the runtime
// state. A domain package builds one via NewGame, registers entities, then
// calls Run (or hands it to pncdsl.Run).
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

	State     *State
	Inventory *Inventory
	UI        *UI
	Audio     *AudioPlayer
	Camera    *Camera

	startID string
	onStart Action

	// runtime
	loaded       *loadedAssets
	input        *Input
	currentScene string
	chars        map[string]*runtimeChar
	scriptRunner Runner
	scriptCtx    *Ctx // for click consumption
	transition   *transition
	dialog       *dialogBox
	activeDialog string
	selectedVerb string
}

// NewGame initializes a game with empty managers and the SCUMM-style verb set.
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
		State:            NewState(),
		Inventory:        NewInventory(),
		Audio:            NewAudioPlayer(),
		Camera:           NewCamera(),
		loaded:           newLoadedAssets(w, h),
		input:            newInput(),
		chars:            make(map[string]*runtimeChar),
		transition:       &transition{},
		selectedVerb:     "look",
	}
	g.UI = newUI(g)
	for _, v := range defaultVerbs() {
		g.VerbManager.Register(v)
	}
	return g
}

func (g *Game) StartAt(name string) *Game { g.startID = name; return g }
func (g *Game) OnStart(a Action) *Game    { g.onStart = a; return g }

// Validate cross-checks name references between managers.
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
	}
	return nil
}

// ----- runtime helpers ---------------------------------------------------

type runtimeChar struct {
	def    Character
	pos    Point
	target Point
	moving bool
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
		if prev != "" {
			old := g.SceneManager.MustGet(prev)
			if old.OnLeave != nil {
				g.queueAction(old.OnLeave, "OnLeave")
			}
		}
		g.currentScene = name
		g.State.NoteVisit(name)
		g.placeActors(name)
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
	c.target = to
	c.moving = true
}

func (g *Game) characterMoving(name string) bool {
	c, ok := g.chars[name]
	return ok && c.moving
}

func (g *Game) tickCharacters(dt float64) {
	for _, c := range g.chars {
		if !c.moving {
			continue
		}
		dx := c.target.X - c.pos.X
		dy := c.target.Y - c.pos.Y
		d := c.pos.Dist(c.target)
		speed := c.def.Speed
		if speed <= 0 {
			speed = 60
		}
		step := speed * dt
		if d <= step {
			c.pos = c.target
			c.moving = false
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
	g.dialog = newDialogBox(&d, &node)
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
	g.dialog.node = &node
	g.dialog.lineIdx = 0
	g.dialog.choiceHits = nil
}

func (g *Game) dialogueActive() bool { return g.dialog != nil }

func (g *Game) showEndCard(text string) { g.UI.endCard = text }
