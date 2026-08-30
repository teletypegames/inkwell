package inkwell

import (
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// The built-in debug font is a fixed 6×16 cell. Every layout measurement
// in the library — and in a domain's own widgets — derives from these.
const (
	GlyphW = 6
	GlyphH = 16
)

// textScratch is the offscreen the coloured draw goes through. It grows
// to fit the widest line ever drawn and is reused after that.
var textScratch *ebiten.Image

// drawText renders s in colour c. ebitenutil.DebugPrintAt only draws
// white, so the glyphs go onto an offscreen first and are blitted back
// tinted with a ColorScale.
func drawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	if s == "" {
		return
	}
	w, h := TextWidth(s)+GlyphW, GlyphH
	if textScratch == nil || textScratch.Bounds().Dx() < w || textScratch.Bounds().Dy() < h {
		nw := w
		if nw < 320 {
			nw = 320
		}
		textScratch = ebiten.NewImage(nw, h)
	}
	textScratch.Clear()
	ebitenutil.DebugPrintAt(textScratch, s, 0, 0)

	if c == nil {
		c = color.White
	}
	r, g, b, a := c.RGBA()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.Scale(
		float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff,
	)
	dst.DrawImage(textScratch.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image), op)
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
