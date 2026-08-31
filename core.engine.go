package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	g.Input.poll()
	g.transition.update(dt)
	g.pumpQueue()

	if g.transition.active && g.transition.out {
		return nil
	}

	// Active script consumes the frame — every widget but the awake ones,
	// which tick first so that a click or a key of theirs is not also read
	// as "skip this line".
	if g.scriptRunner != nil {
		g.tickAwakeWidgets(dt)
		ctx := g.scriptCtx
		ctx.DT = dt
		if g.currentScene != "" {
			s := g.SceneManager.MustGet(g.currentScene)
			ctx.Scene = &s
		}
		st := g.scriptRunner.Tick(ctx)
		if st != StatusRunning {
			g.scriptRunner = nil
			g.scriptCtx = nil
		}
		g.tickCharacters(dt)
		return nil
	}

	uictx := &UICtx{Game: g, DT: dt}

	// Clear per-frame transient state that widgets / scene-resolution write.
	g.SetHoverLabel("")

	// Scene triggers run before widgets/input — a firing trigger queues
	// an action that consumes the next frame's script slot, so this
	// frame's UI work still happens, but the player's clicks won't
	// race a cutscene that's about to start.
	g.tickTriggers()

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

	// Anything still unconsumed is a scene-world interaction.
	e.handleSceneInput()
	g.tickCharacters(dt)
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

// handleSceneInput is the legacy "click on hotspot with active verb" flow,
// running only when no widget has claimed the input.
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
		return
	}
	g.invokeHotspotVerb(h, g.SelectedVerb())
}

// invokeHotspotVerb runs the action bound to verb on hotspot h, honoring
// the SCUMM use-with-item flow when an inventory item is selected.
// Shared by left-click hotspot resolution and the verb-coin widget so
// both paths log, flash and queue identically.
func (g *Game) invokeHotspotVerb(h *Hotspot, verb string) {
	if h == nil {
		return
	}
	sel := g.Inventory.Selected()
	target := hotspotLabel(h)
	if sel != "" {
		logged := "> " + Format(HintVerbOn, verbLabel(g, "use"), g.ItemLabel(sel), target).String()
		if h.OnUseWith != nil {
			if a, ok := h.OnUseWith[sel]; ok && a != nil {
				g.LogAction(logged)
				g.queueAction(a, "useWith hotspot")
				g.Inventory.Select("")
				return
			}
		}
		if it, ok := g.ItemManager.Get(sel); ok && it.OnUseWith != nil {
			if a, ok := it.OnUseWith[h.Name]; ok && a != nil {
				g.LogAction(logged)
				g.queueAction(a, "useWith item")
				g.Inventory.Select("")
				return
			}
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

	// optional backdrop (helps when the background image doesn't fill the screen)
	if th.SceneBackdrop != nil {
		screen.Fill(th.SceneBackdrop)
	}

	// scene background
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		img := g.loaded.image(g.AssetManager, s.Background)
		bw, bh := img.Bounds().Dx(), img.Bounds().Dy()
		if bw > 0 && bh > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(g.Width)/float64(bw), float64(g.Height)/float64(bh))
			screen.DrawImage(img, op)
		}
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 255})
	}

	// characters (y-sorted)
	for _, c := range sortedChars(g) {
		drawCharacter(screen, g, c)
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

func sortedChars(g *Game) []*runtimeChar {
	out := make([]*runtimeChar, 0, len(g.chars))
	for _, c := range g.chars {
		out = append(out, c)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].pos.Y < out[j-1].pos.Y; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func drawCharacter(dst *ebiten.Image, g *Game, c *runtimeChar) {
	w := c.def.W
	h := c.def.H
	if w == 0 {
		w = 24
	}
	if h == 0 {
		h = 56
	}
	if c.def.Sprite != "" && !g.loaded.isPlaceholder(c.def.Sprite) {
		img := g.loaded.image(g.AssetManager, c.def.Sprite)
		// Active animation clip → blit only the current frame's sub-rect
		// from the sprite sheet, scaled into the character's W×H box.
		if r, ok := c.currentFrameRect(); ok {
			ib := img.Bounds()
			clamped := r.Intersect(ib)
			fw, fh := clamped.Dx(), clamped.Dy()
			if fw > 0 && fh > 0 {
				sub := img.SubImage(clamped).(*ebiten.Image)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(w/float64(fw), h/float64(fh))
				op.GeoM.Translate(c.pos.X-w/2, c.pos.Y-h)
				dst.DrawImage(sub, op)
				return
			}
		}
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		if sw > 0 && sh > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(w/float64(sw), h/float64(sh))
			op.GeoM.Translate(c.pos.X-w/2, c.pos.Y-h)
			dst.DrawImage(img, op)
			return
		}
	}
	if c.def.Sprite != "" {
		_ = g.loaded.image(g.AssetManager, c.def.Sprite)
	}
	x := c.pos.X - w/2
	y := c.pos.Y - h
	bodyCol := color.RGBA{200, 200, 200, 255}
	if rgba, ok := c.def.SpeechColor.(color.RGBA); ok {
		bodyCol = rgba
	}
	if w > h {
		drawQuadrupedPlaceholder(dst, x, y, w, h, bodyCol)
	} else {
		drawHumanoidPlaceholder(dst, x, y, w, h, bodyCol)
	}
}

func drawHumanoidPlaceholder(dst *ebiten.Image, x, y, w, h float64, body color.RGBA) {
	pants := color.RGBA{40, 50, 90, 255}
	skin := color.RGBA{245, 200, 155, 255}
	outline := color.RGBA{20, 20, 30, 255}

	headR := h * 0.14
	headCY := y + headR + 1

	pantsTop := y + h*0.62
	pantsH := h - (pantsTop - y) - 1
	legW := w * 0.36
	vector.DrawFilledRect(dst, float32(x+w*0.08), float32(pantsTop), float32(legW), float32(pantsH), pants, false)
	vector.DrawFilledRect(dst, float32(x+w*0.56), float32(pantsTop), float32(legW), float32(pantsH), pants, false)

	torsoTop := headCY + headR
	torsoH := pantsTop - torsoTop
	vector.DrawFilledRect(dst, float32(x), float32(torsoTop), float32(w), float32(torsoH), body, false)
	vector.StrokeRect(dst, float32(x), float32(torsoTop), float32(w), float32(torsoH), 1, outline, false)

	vector.DrawFilledCircle(dst, float32(x+w/2), float32(headCY), float32(headR), skin, true)

	footH := 2.0
	vector.DrawFilledRect(dst, float32(x+w*0.05), float32(y+h-footH), float32(legW+2), float32(footH), outline, false)
	vector.DrawFilledRect(dst, float32(x+w*0.55), float32(y+h-footH), float32(legW+2), float32(footH), outline, false)
}

func drawQuadrupedPlaceholder(dst *ebiten.Image, x, y, w, h float64, body color.RGBA) {
	outline := color.RGBA{20, 20, 30, 255}
	dark := color.RGBA{R: body.R / 2, G: body.G / 2, B: body.B / 2, A: 255}

	bodyY := y + h*0.30
	bodyH := h * 0.65
	vector.DrawFilledRect(dst, float32(x), float32(bodyY), float32(w*0.78), float32(bodyH), body, false)

	headR := h * 0.30
	headCX := x + w - headR
	headCY := y + h*0.50
	vector.DrawFilledCircle(dst, float32(headCX), float32(headCY), float32(headR), body, true)

	earW := w * 0.07
	earH := h * 0.30
	vector.DrawFilledRect(dst, float32(headCX-headR*0.7), float32(y), float32(earW), float32(earH), body, false)
	vector.DrawFilledRect(dst, float32(headCX+headR*0.4), float32(y), float32(earW), float32(earH), body, false)

	vector.DrawFilledCircle(dst, float32(headCX+headR*0.25), float32(headCY-1), 1.2, outline, true)

	vector.DrawFilledRect(dst, float32(x-w*0.04), float32(bodyY+1), float32(w*0.08), 2, dark, false)

	legW := w * 0.07
	vector.DrawFilledRect(dst, float32(x+w*0.10), float32(y+h-3), float32(legW), 3, dark, false)
	vector.DrawFilledRect(dst, float32(x+w*0.55), float32(y+h-3), float32(legW), 3, dark, false)
}
