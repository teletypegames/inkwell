package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// SpeechBubble renders the line set by the Say action above the speaking
// character (or at FallbackY if the speaker has no on-screen position).
type SpeechBubble struct {
	Name      string
	When      Condition
	MaxWidth  int
	Padding   int
	OffsetY   int
	FallbackY int
}

func (s *SpeechBubble) GetName() string        { return s.Name }
func (s *SpeechBubble) VisibleWhen() Condition { return s.When }
func (s *SpeechBubble) Layer() Layer           { return LayerSpeech }
func (s *SpeechBubble) Tick(ctx *UICtx)        {}

func (s *SpeechBubble) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	sp := g.speech
	if !sp.Active {
		return
	}
	maxW := s.MaxWidth
	if maxW <= 0 {
		maxW = 200
	}
	pad := s.Padding
	if pad <= 0 {
		pad = 3
	}
	lines := WrapText(sp.Text, maxW)
	w := 0
	for _, ln := range lines {
		if t := TextWidth(ln); t > w {
			w = t
		}
	}
	w += pad * 2
	h := len(lines)*16 + pad*2

	// position above speaker if known
	cx, cy := g.Width/2, s.FallbackY
	if cy == 0 {
		cy = 14
	}
	if c, ok := g.runtimeChar(sp.Speaker); ok {
		cx = int(c.pos.X)
		cy = int(c.pos.Y-c.def.H) - s.OffsetY
	}
	x := cx - w/2
	y := cy - h
	if x < 2 {
		x = 2
	}
	if x+w > g.Width-2 {
		x = g.Width - 2 - w
	}
	if y < 2 {
		y = 2
	}

	th := g.Theme()
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), th.SpeechBubbleBG, false)

	textCol := th.SpeechDefaultText
	if c, ok := g.CharacterManager.Get(sp.Speaker); ok {
		if rgba, ok := c.SpeechColor.(color.RGBA); ok && rgba.A > 0 {
			textCol = rgba
		}
	}
	for i, ln := range lines {
		g.DrawText(dst, ln, x+pad, y+pad+i*16-2, textCol)
	}
}
