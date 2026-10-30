package textcard

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Colors is the palette a [TextCard] draws with: BG for the backdrop, Text
// for the lines.
type Colors struct {
	BG   color.Color
	Text color.Color
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	TextCardColors() Colors
}

// TextCard is a block of lines centred on the screen: a title card, a
// credits page, a notice. It draws nothing on a host that does not provide
// the screen size, text drawing and [Colors].
type TextCard struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Lines are printed one under another, each resolved into the player's
	// language; an empty Text is a blank line.
	Lines []inkbase.Text
	// Align is how each line sits against the widest one; the block itself
	// is always centred, and AlignCenter centres each line on the screen.
	Align inkbase.Align
	// LineHeight is a line's height in pixels; 0 means the glyph height + 2.
	LineHeight int
	// Backdrop fills the whole screen with [Colors].BG behind the block.
	Backdrop bool
}

// GetName returns the card's Name, the key it is registered under.
func (c *TextCard) GetName() string { return c.Name }

// VisibleWhen returns When, the condition that shows the card; nil shows
// it always.
func (c *TextCard) VisibleWhen() inkbase.Condition { return c.When }

// Layer puts the card on [widget.LayerHUD].
func (c *TextCard) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing: the card has no state of its own to advance.
func (c *TextCard) Tick(ctx *widget.UICtx) {}

func (c *TextCard) lineHeight() int {
	if c.LineHeight > 0 {
		return c.LineHeight
	}
	return font.GlyphH + 2
}

// Draw paints the backdrop when asked for and the lines, centred as a block
// both ways.
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
