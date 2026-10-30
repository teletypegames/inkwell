package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// StatusLine renders a single line: either the active flash message (set
// by FlashLine) or the hover-hint (verb + target under the cursor).
type StatusLine struct {
	Name  string
	When  Condition
	Y     int
	Align Align
	// ScreenWidth, if 0, uses Game.Width.
	ScreenWidth int
}

// GetName returns the name the status line is registered under.
func (s *StatusLine) GetName() string { return s.Name }

// VisibleWhen returns the condition that switches the line on; nil means
// always.
func (s *StatusLine) VisibleWhen() Condition { return s.When }

// Layer puts the line on [LayerHUD].
func (s *StatusLine) Layer() Layer { return LayerHUD }

// Tick counts down the flash message's timer and clears the message when it
// runs out. The game's flash only expires while a StatusLine is ticking.
func (s *StatusLine) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	if g.flashTimer > 0 {
		g.flashTimer -= ctx.DT
		if g.flashTimer <= 0 {
			g.flash = ""
		}
	}
}

// Draw prints the flash message, or the hover hint when there is none, at
// height Y, placed by Align across ScreenWidth. It draws nothing when both
// are empty.
func (s *StatusLine) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
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
	tw := TextWidth(text)
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
