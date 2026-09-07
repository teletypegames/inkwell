package statusline

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

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

type StatusLine struct {
	Name        string
	When        inkbase.Condition
	Y           int
	Align       inkbase.Align
	ScreenWidth int
}

func (s *StatusLine) GetName() string                { return s.Name }
func (s *StatusLine) VisibleWhen() inkbase.Condition { return s.When }
func (s *StatusLine) Layer() widget.Layer            { return widget.LayerHUD }

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
