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

// Host is what a [StatusLine] needs from the game it is drawn for, the
// UICtx's Game. On a Game that does not implement it the line draws nothing.
type Host interface {
	widget.Screen
	// FlashText is the flash message up now, or "" when there is none. The
	// host owns it and its expiry; the line only shows it.
	FlashText() string
	// HoverHint is what the line says when there is no flash: in inkwell
	// the verb and the target under the cursor.
	HoverHint() string
	// StatusLineColors is the palette the line draws with.
	StatusLineColors() Colors
}

// StatusLine renders a single line: either the host's flash message while
// one is up, or the hover hint (verb + target under the cursor). It only
// displays: the flash expires on the host's clock whether a StatusLine is
// registered, visible or not.
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

// GetName returns the name the status line is registered under.
func (s *StatusLine) GetName() string { return s.Name }

// VisibleWhen returns the condition that switches the line on; nil means
// always.
func (s *StatusLine) VisibleWhen() inkbase.Condition { return s.When }

// Layer puts the line on [widget.LayerHUD].
func (s *StatusLine) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing; the line reads the host when it draws.
func (s *StatusLine) Tick(ctx *widget.UICtx) {}

// Draw prints the flash message, or the hover hint when there is none, at
// height Y, placed by Align across ScreenWidth. It draws nothing when both
// are empty.
func (s *StatusLine) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
	if !ok {
		return
	}
	colors := h.StatusLineColors()
	w := s.ScreenWidth
	if w == 0 {
		w, _ = h.ScreenSize()
	}
	text := h.FlashText()
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
