// Package font provides the built-in glyph atlas and text rendering
// functions for inkbase widgets.
package font

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

// The built-in font is a fixed 6×16 cell.
const (
	GlyphW = 6
	GlyphH = 16
)

//go:embed font.png
var fontPNG []byte

const glyphCellW = GlyphW + 1
const glyphCols = 32

var (
	fontAtlas  *ebiten.Image
	glyphCache = map[rune]*ebiten.Image{}
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(fontPNG))
	if err != nil {
		panic("inkbase/font: font.png: " + err.Error())
	}
	fontAtlas = ebiten.NewImageFromImage(img)
}

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

// DrawText renders s in colour c, glyph by glyph out of the atlas.
func DrawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	if s == "" {
		return
	}
	if c == nil {
		c = color.White
	}
	r, g, b, a := c.RGBA()
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(
		float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff,
	)
	px, py := x+1, y
	for _, ru := range s {
		if ru == '\n' {
			px, py = x+1, py+GlyphH
			continue
		}
		op.GeoM.Reset()
		op.GeoM.Translate(float64(px), float64(py))
		dst.DrawImage(glyph(ru), op)
		px += GlyphW
	}
}

// TextWidth is the pixel width of s in the built-in font.
func TextWidth(s string) int { return utf8.RuneCountInString(s) * GlyphW }

// WrapText splits s at word boundaries into lines no wider than maxPx.
func WrapText(s string, maxPx int) []string {
	if maxPx <= GlyphW || TextWidth(s) <= maxPx {
		return []string{s}
	}
	var out []string
	line, word := "", ""
	flush := func() {
		if line != "" {
			out = append(out, line)
			line = ""
		}
	}
	emit := func() {
		if word == "" {
			return
		}
		switch {
		case line == "":
			line = word
		case TextWidth(line+" "+word) <= maxPx:
			line += " " + word
		default:
			flush()
			line = word
		}
		word = ""
	}
	for _, r := range s {
		switch r {
		case ' ':
			emit()
		case '\n':
			emit()
			flush()
		default:
			word += string(r)
		}
	}
	emit()
	flush()
	if len(out) == 0 {
		return []string{s}
	}
	return out
}

// ClipText cuts s to maxPx pixels, ending in an ellipsis when it had to cut.
func ClipText(s string, maxPx int) string {
	r := []rune(s)
	max := maxPx / GlyphW
	switch {
	case max <= 0:
		return ""
	case len(r) <= max:
		return s
	case max == 1:
		return string(r[:1])
	default:
		return string(r[:max-1]) + "…"
	}
}
