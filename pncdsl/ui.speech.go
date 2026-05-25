package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// speechBubble holds the currently-spoken line. The Say action sets it and
// the engine draws it above the speaker (or centered at screen top if the
// speaker is unknown).
type speechBubble struct {
	speaker string
	lines   []string
	pos     Point
	active  bool
}

func (s *speechBubble) set(g *Game, speaker, text string) {
	s.speaker = speaker
	s.lines = wrapText(text, 200)
	s.active = true
	// position above the speaker, if we know where they are
	if c, ok := g.runtimeChar(speaker); ok {
		s.pos = Point{X: c.pos.X, Y: c.pos.Y - c.def.H - 4}
	} else {
		s.pos = Point{X: 160, Y: 14}
	}
}

func (s *speechBubble) clear() {
	s.active = false
	s.speaker = ""
	s.lines = nil
}

func (s *speechBubble) tick(dt float64) { _ = dt }

func (s *speechBubble) draw(dst *ebiten.Image, g *Game) {
	if !s.active || len(s.lines) == 0 {
		return
	}
	maxW := 0
	for _, ln := range s.lines {
		if w := textWidth(ln); w > maxW {
			maxW = w
		}
	}
	pad := 3
	w := maxW + pad*2
	h := len(s.lines)*16 + pad*2
	x := int(s.pos.X) - w/2
	y := int(s.pos.Y) - h
	if x < 2 {
		x = 2
	}
	if x+w > 318 {
		x = 318 - w
	}
	if y < 2 {
		y = 2
	}
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), color.RGBA{0, 0, 0, 200}, false)
	speechColor := color.RGBA{255, 255, 255, 255}
	if c, ok := g.CharacterManager.Get(s.speaker); ok {
		if sc, ok := c.SpeechColor.(color.RGBA); ok {
			speechColor = sc
		}
	}
	for k, ln := range s.lines {
		drawText(dst, ln, x+pad, y+pad+k*16-2, speechColor)
	}
}
