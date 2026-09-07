package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TextCard is a block of lines centred on the screen — a credits roll that
// does not roll, a title card, an apology. It is the plainest thing a Stage
// that is not the world can be made of, and it needs no tree, no rows and no
// navigation to be one.
type TextCard struct {
	Name string
	When Condition
	// Lines are printed one under another, each resolved into the player's
	// language. An empty Text is a blank line, which is how a block is
	// spaced.
	Lines []Text
	// Align is how each line sits against the others; the block itself is
	// always centred on the screen.
	Align Align
	// LineHeight is a line's height; 0 gives GlyphH + 2.
	LineHeight int
	// Backdrop fills the screen behind the block. Nil leaves whatever the
	// Stage put there.
	Backdrop bool
}

func (c *TextCard) GetName() string        { return c.Name }
func (c *TextCard) VisibleWhen() Condition { return c.When }
func (c *TextCard) Layer() Layer           { return LayerHUD }
func (c *TextCard) Tick(ctx *UICtx)        {}

func (c *TextCard) lineHeight() int {
	if c.LineHeight > 0 {
		return c.LineHeight
	}
	return GlyphH + 2
}

func (c *TextCard) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	th := g.Theme()
	if c.Backdrop {
		vector.DrawFilledRect(dst, 0, 0, float32(g.Width), float32(g.Height), th.EndCardBG, false)
	}
	if len(c.Lines) == 0 {
		return
	}
	lh := c.lineHeight()
	widest := 0
	rendered := make([]string, len(c.Lines))
	for i, line := range c.Lines {
		rendered[i] = line.String()
		if w := TextWidth(rendered[i]); w > widest {
			widest = w
		}
	}
	top := (g.Height - len(rendered)*lh) / 2
	left := (g.Width - widest) / 2
	for i, line := range rendered {
		x := left
		switch c.Align {
		case AlignCenter:
			x = (g.Width - TextWidth(line)) / 2
		case AlignRight:
			x = left + widest - TextWidth(line)
		}
		g.DrawText(dst, line, x, top+i*lh, th.EndCardText)
	}
}
