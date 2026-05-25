package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawCursor renders a simple crosshair at the mouse position. When an
// inventory item is selected, the item sprite is drawn instead.
func drawCursor(dst *ebiten.Image, g *Game) {
	x, y := g.input.Pos()
	if sel := g.Inventory.Selected(); sel != "" {
		if it, ok := g.ItemManager.Get(sel); ok {
			img := g.loaded.image(g.AssetManager, it.Sprite)
			sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
			if sw > 0 && sh > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(16/float64(sw), 16/float64(sh))
				op.GeoM.Translate(float64(x-8), float64(y-8))
				dst.DrawImage(img, op)
				return
			}
		}
	}
	c := color.RGBA{255, 255, 255, 255}
	vector.StrokeLine(dst, float32(x-3), float32(y), float32(x+4), float32(y), 1, c, false)
	vector.StrokeLine(dst, float32(x), float32(y-3), float32(x), float32(y+4), 1, c, false)
}
