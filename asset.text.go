package inkwell

import (
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

// The built-in debug font is a fixed 6×16 cell. Every layout measurement
// in the library — and in a domain's own widgets — derives from these.
const (
	GlyphW = 6
	GlyphH = 16
)

// drawText renders s in colour c, glyph by glyph out of the atlas in
// i18n.font.go. Every cell carries a black drop shadow under a white glyph,
// so scaling the blit by c tints the glyph and leaves the shadow black.
func drawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
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
	// The shadow pass is baked one pixel right, so the run starts one pixel
	// right of x — which is where ebitenutil put it, and what every layout
	// in the engine was measured against.
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
// A newline in s breaks the line where it stands.
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

// ClipText cuts s to maxPx pixels, ending in an ellipsis when it had to
// cut. A width with no room at all returns the empty string.
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
