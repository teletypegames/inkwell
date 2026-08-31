package inkwell

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:generate go -C internal/fontgen run .

// font.png is the engine's own glyph atlas, baked by internal/fontgen. It
// covers Latin-1 and Latin Extended-A, which is what puts ő and ű on screen
// — ebitenutil's debug font, which the engine drew through before, stops at
// U+00FF and renders them, and the em dashes and ellipses this engine has
// always emitted, as blanks.
//
//go:embed font.png
var fontPNG []byte

// glyphCellW is one glyph plus the column its baked drop shadow falls into.
// Cells are blitted at GlyphW pitch, so that column overlaps the next glyph
// and composites under it.
const glyphCellW = GlyphW + 1

const glyphCols = 32

var (
	fontAtlas  *ebiten.Image
	glyphCache = map[rune]*ebiten.Image{}
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(fontPNG))
	if err != nil {
		panic("inkwell: font.png: " + err.Error())
	}
	fontAtlas = ebiten.NewImageFromImage(img)
}

// glyphIndex is a rune's cell in the atlas, or glyphNotdef when the atlas
// does not cover it.
func glyphIndex(r rune) int {
	base := 0
	for _, rg := range glyphRanges {
		if r >= rg[0] && r <= rg[1] {
			return base + int(r-rg[0])
		}
		base += int(rg[1]-rg[0]) + 1
	}
	return glyphNotdef
}

func glyph(r rune) *ebiten.Image {
	if img, ok := glyphCache[r]; ok {
		return img
	}
	i := glyphIndex(r)
	x, y := (i%glyphCols)*glyphCellW, (i/glyphCols)*GlyphH
	img := fontAtlas.SubImage(image.Rect(x, y, x+glyphCellW, y+GlyphH)).(*ebiten.Image)
	glyphCache[r] = img
	return img
}
