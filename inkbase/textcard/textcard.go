package textcard

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

type Colors struct {
	BG   color.Color
	Text color.Color
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	TextCardColors() Colors
}

type TextCard struct {
	Name       string
	When       inkbase.Condition
	Lines      []inkbase.Text
	Align      inkbase.Align
	LineHeight int
	Backdrop   bool
}

func (c *TextCard) GetName() string                { return c.Name }
func (c *TextCard) VisibleWhen() inkbase.Condition { return c.When }
func (c *TextCard) Layer() widget.Layer            { return widget.LayerHUD }
func (c *TextCard) Tick(ctx *widget.UICtx)         {}

func (c *TextCard) lineHeight() int {
	if c.LineHeight > 0 {
		return c.LineHeight
	}
	return font.GlyphH + 2
}

func (c *TextCard) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	colors := h.TextCardColors()
	w, sh := h.ScreenSize()
	if c.Backdrop {
		vector.DrawFilledRect(dst, 0, 0, float32(w), float32(sh), colors.BG, false)
	}
	if len(c.Lines) == 0 {
		return
	}
	lh := c.lineHeight()
	widest := 0
	rendered := make([]string, len(c.Lines))
	for i, line := range c.Lines {
		rendered[i] = line.String()
		if tw := font.TextWidth(rendered[i]); tw > widest {
			widest = tw
		}
	}
	top := (sh - len(rendered)*lh) / 2
	left := (w - widest) / 2
	for i, line := range rendered {
		x := left
		switch c.Align {
		case inkbase.AlignCenter:
			x = (w - font.TextWidth(line)) / 2
		case inkbase.AlignRight:
			x = left + widest - font.TextWidth(line)
		}
		h.DrawText(dst, line, x, top+i*lh, colors.Text)
	}
}
