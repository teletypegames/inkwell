package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The scene's cast, painted back to front: a character standing lower on
// the picture is nearer the camera, so it is drawn last. A character with
// no usable sprite falls back to one of the two placeholder shapes, which
// is what keeps a game playable before any art exists.

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
