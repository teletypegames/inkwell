package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Cursor renders the mouse cursor. When the player has an inventory item
// selected, the item's sprite is drawn at the cursor instead of the
// default crosshair. It sits on LayerCursor, so it draws on top.
type Cursor struct {
	Name string
	When Condition
}

// GetName returns the name the cursor is registered under.
func (c *Cursor) GetName() string { return c.Name }

// VisibleWhen returns the condition that switches the cursor on; nil means
// always.
func (c *Cursor) VisibleWhen() Condition { return c.When }

// Layer puts the cursor on [LayerCursor], above everything else.
func (c *Cursor) Layer() Layer { return LayerCursor }

// Tick does nothing; the cursor follows the pointer when it draws.
func (c *Cursor) Tick(ctx *UICtx) {}

// Draw paints the held item's sprite, 16 pixels square and centred on the
// pointer, or, with nothing held or no sprite, a small crosshair in
// Theme.CursorColor.
func (c *Cursor) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	x, y := g.Input.Pos()
	if sel := g.Inventory.Selected(); sel != "" {
		if it, ok := g.ItemManager.Get(sel); ok {
			if drawSprite(dst, g, it.Sprite, itemFrame(g, it.Sprite), FacingRight, float64(x-8), float64(y-8), 16, 16) {
				return
			}
		}
	}
	col := g.Theme().CursorColor
	vector.StrokeLine(dst, float32(x-3), float32(y), float32(x+4), float32(y), 1, col, false)
	vector.StrokeLine(dst, float32(x), float32(y-3), float32(x), float32(y+4), 1, col, false)
}
