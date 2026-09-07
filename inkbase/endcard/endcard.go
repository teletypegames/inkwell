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

type EndCard struct {
	Name string
	When inkbase.Condition
}

func (e *EndCard) GetName() string                { return e.Name }
func (e *EndCard) VisibleWhen() inkbase.Condition { return e.When }
func (e *EndCard) Layer() widget.Layer            { return widget.LayerCurtain }

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
