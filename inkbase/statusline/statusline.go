package statusline

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Colors is the palette a [StatusLine] draws with: Flash for a flash
// message, Status for the hover hint.
type Colors struct {
	Flash  color.Color
	Status color.Color
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	FlashText() (string, float64)
	SetFlash(text string, timer float64)
	HoverHint() string
	StatusLineColors() Colors
}

// StatusLine renders a single line of text: the host's flash message while
// one is up, otherwise the hover hint for what is under the cursor. It does
// nothing on a host that does not provide them, the screen size, text
// drawing and [Colors].
type StatusLine struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Y is the line's top in pixels.
	Y int
	// Align places the line: left or right with a 2 pixel margin, centred
	// otherwise.
	Align inkbase.Align
	// ScreenWidth is the width Align works against; 0 uses the screen's.
	ScreenWidth int
}

// GetName returns the status line's Name, the key it is registered under.
func (s *StatusLine) GetName() string { return s.Name }

// VisibleWhen returns When, the condition that shows the status line; nil shows
// it always.
func (s *StatusLine) VisibleWhen() inkbase.Condition { return s.When }

// Layer puts the status line on [widget.LayerHUD].
func (s *StatusLine) Layer() widget.Layer { return widget.LayerHUD }

// Tick counts the host's flash timer down by the frame's time and clears
// the flash message when it runs out.
func (s *StatusLine) Tick(ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	text, timer := h.FlashText()
	if timer > 0 {
		timer -= ctx.DT
		if timer <= 0 {
			text = ""
			timer = 0
		}
		h.SetFlash(text, timer)
	}
}

// Draw prints the flash message, or the hover hint when there is none, at
// Y; with neither it draws nothing.
func (s *StatusLine) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	colors := h.StatusLineColors()
	w := s.ScreenWidth
	if w == 0 {
		w, _ = h.ScreenSize()
	}
	text, _ := h.FlashText()
	col := colors.Flash
	if text == "" {
		text = h.HoverHint()
		col = colors.Status
	}
	if text == "" {
		return
	}
	tw := font.TextWidth(text)
	var x int
	switch s.Align {
	case inkbase.AlignLeft:
		x = 2
	case inkbase.AlignRight:
		x = w - tw - 2
	default:
		x = (w - tw) / 2
	}
	h.DrawText(dst, text, x, s.Y, col)
}
