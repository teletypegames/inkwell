package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Cursor renders the mouse cursor. When the player has an inventory item
// selected, the item's sprite is drawn at the cursor instead of the
// default crosshair. Registered last so it draws on top.
type Cursor struct {
	Name string
}

func (c *Cursor) GetName() string { return c.Name }
func (c *Cursor) Layer() Layer    { return LayerCursor }
func (c *Cursor) Tick(ctx *UICtx) {}

func (c *Cursor) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	x, y := g.Input.Pos()
	if sel := g.Inventory.Selected(); sel != "" {
		if it, ok := g.ItemManager.Get(sel); ok {
			img := g.loaded.image(g.AssetManager, it.Sprite)
			sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
			if sw > 0 && sh > 0 && !g.loaded.isPlaceholder(it.Sprite) {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(16/float64(sw), 16/float64(sh))
				op.GeoM.Translate(float64(x-8), float64(y-8))
				dst.DrawImage(img, op)
				return
			}
		}
	}
	col := g.Theme().CursorColor
	vector.StrokeLine(dst, float32(x-3), float32(y), float32(x+4), float32(y), 1, col, false)
	vector.StrokeLine(dst, float32(x), float32(y-3), float32(x), float32(y+4), 1, col, false)
}
