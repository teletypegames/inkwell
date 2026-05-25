package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// engine is the ebiten.Game adapter. It owns the per-frame Update/Draw flow
// and delegates everything stateful to *Game.
type engine struct {
	g *Game
}

func (e *engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return e.g.Width, e.g.Height
}

func (e *engine) Update() error {
	g := e.g
	dt := 1.0 / 60.0
	g.input.poll()
	g.UI.tick(dt)
	g.transition.update(dt)

	if g.transition.active && g.transition.out {
		// during fade-out, freeze input
		return nil
	}

	if g.UI.endCard != "" {
		return nil
	}

	if g.scriptRunner != nil {
		ctx := g.scriptCtx
		ctx.DT = dt
		// keep Scene reference fresh in case the script changed it
		if g.currentScene != "" {
			s := g.SceneManager.MustGet(g.currentScene)
			ctx.Scene = &s
		}
		s := g.scriptRunner.Tick(ctx)
		if s != StatusRunning {
			g.scriptRunner = nil
			g.scriptCtx = nil
		}
		g.tickCharacters(dt)
		return nil
	}

	// dialog click handling
	if g.dialog != nil {
		if g.input.LeftClicked() {
			g.input.ConsumeLeft()
			ctx := g.makeCtx()
			picked := g.dialog.handleClick(ctx, g.input.Point())
			if picked != nil {
				// record once-tracking
				if picked.Once {
					g.State.NoteTalked(g.dialog.node.Name + ":" + picked.Text)
				}
				// run the choice's actions as a synthetic script
				if len(picked.Actions) > 0 {
					g.queueAction(Seq(picked.Actions...), "choice")
				}
			}
		}
		g.tickCharacters(dt)
		return nil
	}

	// free interaction
	e.handleFreeInput()
	g.tickCharacters(dt)
	return nil
}

func (e *engine) handleFreeInput() {
	g := e.g
	g.UI.rebuildVerbButtons()
	mp := g.input.Point()

	// hover label resolution
	g.UI.hoverLabel = ""
	if mp.Y < uiPanelY-4 {
		if h := e.hotspotAt(mp); h != nil {
			if h.Label != "" {
				g.UI.hoverLabel = h.Label
			} else {
				g.UI.hoverLabel = h.Name
			}
		}
	} else {
		// hovering UI: show item description when over an inv slot
		if hit := g.UI.hitTest(mp); len(hit) > 4 && hit[:4] == "inv:" {
			name := hit[4:]
			if it, ok := g.ItemManager.Get(name); ok && it.Description != "" {
				g.UI.hoverLabel = it.Description
			} else {
				g.UI.hoverLabel = name
			}
		}
	}

	if g.input.RightClicked() {
		g.input.ConsumeRight()
		// right click: deselect item / reset verb
		if g.Inventory.Selected() != "" {
			g.Inventory.Select("")
		} else {
			g.selectedVerb = "look"
		}
	}

	if !g.input.LeftClicked() {
		return
	}

	// UI hits first
	if hit := g.UI.hitTest(mp); hit != "" {
		g.input.ConsumeLeft()
		switch hit[:4] {
		case "verb":
			g.selectedVerb = hit[5:]
		case "inv:":
			name := hit[4:]
			if g.selectedVerb == "look" {
				if it, ok := g.ItemManager.Get(name); ok && it.Description != "" {
					g.queueAction(Say("player", it.Description), "inv-look")
					return
				}
			}
			if g.Inventory.Selected() == name {
				g.Inventory.Select("")
			} else {
				g.Inventory.Select(name)
			}
		}
		return
	}

	// game-world click
	g.input.ConsumeLeft()
	h := e.hotspotAt(mp)
	if h == nil {
		return
	}
	sel := g.Inventory.Selected()
	if sel != "" {
		// use-with: first check hotspot's OnUseWith, then item's
		if h.OnUseWith != nil {
			if a, ok := h.OnUseWith[sel]; ok && a != nil {
				g.queueAction(a, "useWith hotspot")
				g.Inventory.Select("")
				return
			}
		}
		if it, ok := g.ItemManager.Get(sel); ok && it.OnUseWith != nil {
			if a, ok := it.OnUseWith[h.Name]; ok && a != nil {
				g.queueAction(a, "useWith item")
				g.Inventory.Select("")
				return
			}
		}
		g.UI.FlashLine("Nem ehhez.")
		g.Inventory.Select("")
		return
	}
	a := h.handler(g.selectedVerb)
	if a == nil {
		if v, ok := g.VerbManager.Get(g.selectedVerb); ok && v.Default != nil {
			a = v.Default
		}
	}
	if a == nil {
		g.UI.FlashLine("Semmi említésre méltó.")
		return
	}
	g.queueAction(a, "hotspot "+g.selectedVerb+" "+h.Name)
}

func (e *engine) hotspotAt(p Point) *Hotspot {
	g := e.g
	if g.currentScene == "" {
		return nil
	}
	s := g.SceneManager.MustGet(g.currentScene)
	for i := range s.Hotspots {
		h := &s.Hotspots[i]
		if h.Area != nil && h.Area.Contains(p) {
			return h
		}
	}
	return nil
}

func (e *engine) Draw(screen *ebiten.Image) {
	g := e.g
	// scene background
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		img := g.loaded.image(g.AssetManager, s.Background)
		op := &ebiten.DrawImageOptions{}
		bw, bh := img.Bounds().Dx(), img.Bounds().Dy()
		if bw > 0 && bh > 0 {
			op.GeoM.Scale(float64(g.Width)/float64(bw), float64(g.Height)/float64(bh))
		}
		screen.DrawImage(img, op)

		// hotspot debug outlines
		if DebugLog {
			for _, h := range s.Hotspots {
				if r, ok := h.Area.(Rectangle); ok {
					vector.StrokeRect(screen, float32(r.X), float32(r.Y), float32(r.W), float32(r.H), 1, color.RGBA{255, 255, 0, 200}, false)
				}
			}
		}
	} else {
		screen.Fill(color.RGBA{0, 0, 0, 255})
	}

	// characters (y-sorted)
	for _, c := range sortedChars(g) {
		drawCharacter(screen, g, c)
	}

	// speech bubble
	g.UI.speech.draw(screen, g)

	// dialog box (if active)
	if g.dialog != nil {
		g.dialog.draw(screen, g)
	} else {
		g.UI.drawHUD(screen)
	}

	// transition overlay
	g.transition.draw(screen, g.Width, g.Height)

	// end card
	if g.UI.endCard != "" {
		vector.DrawFilledRect(screen, 0, 0, float32(g.Width), float32(g.Height), color.RGBA{0, 0, 0, 230}, false)
		w := textWidth(g.UI.endCard)
		drawText(screen, g.UI.endCard, (g.Width-w)/2, g.Height/2-8, color.White)
	}

	// cursor on top
	drawCursor(screen, g)
}

func sortedChars(g *Game) []*runtimeChar {
	out := make([]*runtimeChar, 0, len(g.chars))
	for _, c := range g.chars {
		out = append(out, c)
	}
	// insertion sort by Y (small N)
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
		w = 14
	}
	if h == 0 {
		h = 26
	}
	if c.def.Sprite != "" {
		img := g.loaded.image(g.AssetManager, c.def.Sprite)
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		if sw > 0 && sh > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(w/float64(sw), h/float64(sh))
			op.GeoM.Translate(c.pos.X-w/2, c.pos.Y-h)
			dst.DrawImage(img, op)
			return
		}
	}
	col := color.RGBA{200, 200, 200, 255}
	if rgba, ok := c.def.SpeechColor.(color.RGBA); ok {
		col = rgba
	}
	vector.DrawFilledRect(dst, float32(c.pos.X-w/2), float32(c.pos.Y-h), float32(w), float32(h), col, false)
}
