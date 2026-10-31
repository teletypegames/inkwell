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

// Host is what a [TextCard] needs from the game it is drawn for, the
// UICtx's Game. On a Game that does not implement it the card draws nothing.
type Host interface {
	widget.Screen
	// TextCardColors is the palette the card draws with.
	TextCardColors() Colors
}

// TextCard is a block of lines centred on the screen — a credits roll that
// does not roll, a title card, an apology. It is the plainest thing a Stage
// that is not the world can be made of, and it needs no tree, no rows and no
// navigation to be one.
type TextCard struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Lines are printed one under another, each resolved into the player's
	// language. An empty Text is a blank line, which is how a block is
	// spaced.
	Lines []inkbase.Text
	// Align is how each line sits against the others; the block itself is
	// always centred on the screen.
	Align inkbase.Align
	// LineHeight is a line's height in pixels; 0 gives the glyph height + 2.
	LineHeight int
	// Backdrop fills the screen behind the block with [Colors].BG. False
	// leaves whatever the Stage put there.
	Backdrop bool
}

// GetName returns the name the card is registered under.
func (c *TextCard) GetName() string { return c.Name }

// VisibleWhen returns the condition that switches the card on; nil means
// always.
func (c *TextCard) VisibleWhen() inkbase.Condition { return c.When }

// Layer puts the card on [widget.LayerHUD].
func (c *TextCard) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing; the card only displays.
func (c *TextCard) Tick(ctx *widget.UICtx) {}

func (c *TextCard) lineHeight() int {
	if c.LineHeight > 0 {
		return c.LineHeight
	}
	return font.GlyphH + 2
}

// Draw paints the Backdrop when asked, then the Lines as a block centred on
// the screen, each placed by Align within the width of the widest.
func (c *TextCard) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
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
