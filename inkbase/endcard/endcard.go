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

// Host is what an [EndCard] needs from the game it is ticked and drawn for,
// the UICtx's Game. On a Game that does not implement it the card does
// nothing.
type Host interface {
	widget.Screen
	// WidgetInput is the frame's consume-on-use input; the card swallows
	// its clicks.
	WidgetInput() *input.Input
	// EndCardText is the end text; "" means the card is down.
	EndCardText() string
	// EndCardColors is the palette the card draws with.
	EndCardColors() Colors
}

// EndCard is a fullscreen overlay with the host's end text centred on it,
// up whenever that text is not empty — in inkwell, once the ShowEnd action
// fires. While it is up it consumes every click so nothing underneath
// reacts.
type EndCard struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
}

// GetName returns the name the card is registered under.
func (e *EndCard) GetName() string { return e.Name }

// VisibleWhen returns the condition that switches the card on; nil means
// always.
func (e *EndCard) VisibleWhen() inkbase.Condition { return e.When }

// Layer puts the card on [widget.LayerCurtain], over the whole game.
func (e *EndCard) Layer() widget.Layer { return widget.LayerCurtain }

// Tick consumes every left and right click while the end card is up.
func (e *EndCard) Tick(ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
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

// Draw fills the screen and prints the end text centred on it, while the
// end card is up.
func (e *EndCard) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
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
