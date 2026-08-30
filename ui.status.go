package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// StatusLine renders a single line: either the active flash message (set
// by FlashLine) or the hover-hint (verb + target under the cursor).
type StatusLine struct {
	Name  string
	Y     int
	Align Align
	// ScreenWidth, if 0, uses Game.Width.
	ScreenWidth int
}

func (s *StatusLine) GetName() string { return s.Name }
func (s *StatusLine) Layer() Layer    { return LayerHUD }

func (s *StatusLine) Tick(ctx *UICtx) {
	g := ctx.Game
	if g.flashTimer > 0 {
		g.flashTimer -= ctx.DT
		if g.flashTimer <= 0 {
			g.flash = ""
		}
	}
}

func (s *StatusLine) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	th := g.Theme()
	w := s.ScreenWidth
	if w == 0 {
		w = g.Width
	}
	text := g.flash
	col := th.FlashText
	if text == "" {
		text = g.HoverHint()
		col = th.StatusText
	}
	if text == "" {
		return
	}
	tw := textWidth(text)
	var x int
	switch s.Align {
	case AlignLeft:
		x = 2
	case AlignRight:
		x = w - tw - 2
	default:
		x = (w - tw) / 2
	}
	g.DrawText(dst, text, x, s.Y, col)
}
