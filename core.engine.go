package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// engine is the ebiten.Game adapter. Each frame it polls input, advances
// the active script (if any), then lets every registered Widget tick
// (reverse registration order, so the top widget claims input first).
// Whatever survives is offered to the scene as a hotspot click.
type engine struct {
	g *Game
}

func (e *engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return e.g.Width, e.g.Height
}

func (e *engine) Update() error {
	g := e.g
	dt := 1.0 / 60.0

	// A screen change asked for last frame lands here, before anything is
	// ticked, so nothing is torn down underneath a runner mid-Tick.
	g.applyPending()
	if g.quitting {
		return ebiten.Termination
	}

	g.clock += dt
	g.Input.poll()
	g.transition.update(dt)
	g.pumpQueue()

	if g.transition.active && g.transition.out {
		return nil
	}

	// Active script consumes the frame — every widget but the awake ones,
	// which tick first so that a click or a key of theirs is not also read
	// as "skip this line".
	world := g.OnWorldStage()

	if g.scriptRunner != nil {
		g.tickAwakeWidgets(dt)
		ctx := g.scriptCtx
		ctx.DT = dt
		if g.currentScene != "" {
			ctx.Scene = g.SceneManager.MustGet(g.currentScene)
		}
		st := g.scriptRunner.Tick(ctx)
		if st != StatusRunning {
			g.scriptRunner = nil
			g.scriptCtx = nil
		}
		if world {
			g.tickCharacters(dt)
		}
		return nil
	}

	uictx := &UICtx{Game: g, DT: dt}

	// Clear per-frame transient state that widgets / scene-resolution write.
	g.SetHoverLabel("")

	// Scene triggers run before widgets/input — a firing trigger queues
	// an action that consumes the next frame's script slot, so this
	// frame's UI work still happens, but the player's clicks won't
	// race a cutscene that's about to start. Off the world Stage there
	// is no scene to trigger anything.
	if world {
		g.tickTriggers()
	}

	// Top-down input: the widget on the highest layer gets the click
	// first, then the one below it, etc. A widget signals "I took it"
	// via g.Input.ConsumeLeft / ConsumeRight.
	cctx := g.makeCtx()
	for _, w := range reversedWidgets(g.WidgetManager) {
		if !WidgetVisible(cctx, w) {
			continue
		}
		w.Tick(uictx)
	}

	// Anything still unconsumed is a scene-world interaction — and only the
	// world Stage has one, so a click that falls through a menu reaches
	// nothing rather than the alley behind it.
	if world {
		e.handleSceneInput()
		g.tickCharacters(dt)
	}
	return nil
}

// tickAwakeWidgets ticks the widgets that declared themselves awake, top
// layer down, on a frame a script would otherwise have to itself.
func (g *Game) tickAwakeWidgets(dt float64) {
	uictx := &UICtx{Game: g, DT: dt}
	cctx := g.makeCtx()
	for _, w := range reversedWidgets(g.WidgetManager) {
		if WidgetAwake(w) && WidgetVisible(cctx, w) {
			w.Tick(uictx)
		}
	}
}

// handleSceneInput is the "click on hotspot with active verb" flow, running
// only when no widget has claimed the input. With Game.WalkOnClick, a click
// that lands on no hotspot is a walk: Game.Player goes there, routed
// through the walkboxes.
func (e *engine) handleSceneInput() {
	g := e.g
	mp := g.Input.Point()
	// hover label resolution: pointer over a hotspot in the current scene
	if h := e.hotspotAt(mp); h != nil {
		g.SetHoverLabel(hotspotLabel(h).String())
	}

	if g.Input.RightClicked() {
		g.Input.ConsumeRight()
		if g.Inventory.Selected() != "" {
			g.Inventory.Select("")
		} else {
			g.SetSelectedVerb("look")
		}
	}

	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	h := e.hotspotAt(mp)
	if h == nil {
		// Empty floor: with WalkOnClick, the player walks there. This
		// goes straight to walkCharacter rather than through the queue as
		// a Walk action, because a queued walk would be a script, and a
		// script owns the frame — the next click could not redirect it
		// until it arrived.
		if g.WalkOnClick && g.Player != nil && g.SceneArea().Contains(mp) {
			g.walkCharacter(g.Player.Name, mp)
		}
		return
	}
	g.InvokeVerb(h, g.SelectedVerb())
}

// InvokeVerb runs the action bound to verb on hotspot h, honoring the
// SCUMM use-with-item flow when an inventory item is selected: the
// hotspot's own OnUseWith first, then the item's, then UseWithFail.
//
// It is the one path from "the player picked a verb and a target" to an
// action on the queue, and it is exported because a target does not have
// to be in the picture. A HUD widget that owns something clickable — a
// device in the corner, a slot on a panel — hands its own Hotspot here
// and gets the same logging, the same use-with resolution and the same
// nothing-happens line as a painted hotspot, rather than re-deciding all
// of that in its Tick. HotspotPanel is that widget.
func (g *Game) InvokeVerb(h *Hotspot, verb string) {
	if h == nil {
		return
	}
	sel, held := g.ItemManager.Get(g.Inventory.Selected())
	target := hotspotLabel(h)
	if held {
		logged := "> " + Format(HintVerbOn, verbLabel(g, "use"), g.ItemLabel(sel.Name), target).String()
		if h.OnUseWith != nil {
			if a, ok := h.OnUseWith[sel]; ok && a != nil {
				g.LogAction(logged)
				g.queueAction(a, "useWith hotspot")
				g.Inventory.Select("")
				return
			}
		}
		if h.OnUseWithAny != nil {
			g.LogAction(logged)
			g.queueAction(h.OnUseWithAny, "useWith any")
			g.Inventory.Select("")
			return
		}
		if g.UseWithFail != nil {
			g.LogAction(logged)
			g.queueAction(g.UseWithFail(sel, h), "useWith fail")
		} else {
			g.FlashLine(MsgNotThat)
		}
		g.Inventory.Select("")
		return
	}
	a := h.handler(verb)
	if a == nil {
		if v, ok := g.VerbManager.Get(verb); ok && v.Default != nil {
			a = v.Default
		}
	}
	if a == nil {
		g.FlashLine(MsgNothing)
		return
	}
	g.LogAction("> " + Format(HintVerb, verbLabel(g, verb), target).String())
	g.queueAction(a, "hotspot "+verb+" "+h.Name)
}

// hotspotLabel is what to call a hotspot on screen: its Label, or its Name
// when nobody gave it one.
func hotspotLabel(h *Hotspot) Text {
	if h.Label.Empty() {
		return plain(h.Name)
	}
	return h.Label
}

func verbLabel(g *Game, name string) Text {
	if v, ok := g.VerbManager.Get(name); ok {
		return v.Label
	}
	return plain(name)
}

func (e *engine) hotspotAt(p Point) *Hotspot { return e.g.HotspotAt(p) }

func (e *engine) Draw(screen *ebiten.Image) {
	g := e.g
	th := g.Theme()

	// The backdrop is what shows wherever nothing is painted over it: a
	// theme's own colour, or black.
	if th.SceneBackdrop != nil {
		screen.Fill(th.SceneBackdrop)
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 255})
	}

	world := g.OnWorldStage()

	// The picture: the current scene on the world Stage, its layers behind
	// and in front of its cast, and whatever the Stage itself carries
	// anywhere else. A scene with no layers and a Stage with no Background leave the
	// backdrop showing.
	switch {
	case world && g.currentScene != "":
		drawScene(screen, g, g.SceneManager.MustGet(g.currentScene))
	case !world && g.currentStage != nil:
		drawFullScreen(screen, g, g.currentStage.Background)
	}

	// widgets from the bottom layer up
	uictx := &UICtx{Game: g, DT: 1.0 / 60.0}
	cctx := g.makeCtx()
	for _, w := range orderedWidgets(g.WidgetManager) {
		if !WidgetVisible(cctx, w) {
			continue
		}
		w.Draw(screen, uictx)
	}

	// transition overlay on top of everything except cursor (cursor is the
	// last widget, so it has already drawn above the transition? no — the
	// fade should cover the cursor too; draw it after the widgets).
	g.transition.draw(screen, g.Width, g.Height)
}

// drawFullScreen paints a named image over the whole window, which is what
// both a scene background and a Stage background are.
func drawFullScreen(dst *ebiten.Image, g *Game, asset *Asset) {
	if asset == nil {
		return
	}
	img := g.loaded.image(asset)
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w == 0 || h == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.Width)/float64(w), float64(g.Height)/float64(h))
	dst.DrawImage(img, op)
}
