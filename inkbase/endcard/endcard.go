package endcard

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/input"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Colors is the palette an [EndCard] draws with: BG fills the screen, Text
// is the message.
type Colors struct {
	BG   color.Color
	Text color.Color
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	WidgetInput() *input.Input
	EndCardText() string
	EndCardColors() Colors
}

// EndCard is a fullscreen overlay with the host's end text centred on it,
// shown whenever that text is not empty. While it is up it swallows every
// click, so nothing underneath reacts.
type EndCard struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
}

// GetName returns the end card's Name, the key it is registered under.
func (e *EndCard) GetName() string { return e.Name }

// VisibleWhen returns When, the condition that shows the end card; nil shows
// it always.
func (e *EndCard) VisibleWhen() inkbase.Condition { return e.When }

// Layer puts the end card on [widget.LayerCurtain], over everything but the cursor.
func (e *EndCard) Layer() widget.Layer { return widget.LayerCurtain }

// Tick consumes this frame's left and right clicks while the end text is
// set.
func (e *EndCard) Tick(ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	if h.EndCardText() == "" {
		return
	}
	in := h.WidgetInput()
	if in.LeftClicked() {
		in.ConsumeLeft()
	}
	if in.RightClicked() {
		in.ConsumeRight()
	}
}

// Draw fills the screen and prints the end text centred on it; with no end
// text it draws nothing.
func (e *EndCard) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	text := h.EndCardText()
	if text == "" {
		return
	}
	colors := h.EndCardColors()
	w, sh := h.ScreenSize()
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(sh), colors.BG, false)
	tw := font.TextWidth(text)
	h.DrawText(dst, text, (w-tw)/2, sh/2-8, colors.Text)
}
