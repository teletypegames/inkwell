package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// drawText is a minimal text draw using ebiten's built-in debug font.
// The glyphs are 6×16-ish; good enough for a SCUMM-style 320×200 demo.
func drawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	if s == "" {
		return
	}
	// ebitenutil.DebugPrintAt only draws white; for color we render to a
	// tiny offscreen, then ColorScale-tint when blitting. Simpler: just
	// use the white draw and skip color. (TODO: text/v2 once needed.)
	_ = c
	ebitenutil.DebugPrintAt(dst, s, x, y)
}

// textWidth is a coarse pixel-width estimate, used for centering.
func textWidth(s string) int { return len(s) * 6 }

// wrapText splits s at word boundaries into lines no wider than maxPx.
func wrapText(s string, maxPx int) []string {
	if maxPx <= 0 || textWidth(s) <= maxPx {
		return []string{s}
	}
	var (
		out  []string
		line string
	)
	flush := func() {
		if line != "" {
			out = append(out, line)
			line = ""
		}
	}
	word := ""
	for _, r := range s {
		if r == ' ' || r == '\n' {
			if line == "" {
				line = word
			} else if textWidth(line+" "+word) <= maxPx {
				line += " " + word
			} else {
				flush()
				line = word
			}
			word = ""
			if r == '\n' {
				flush()
			}
			continue
		}
		word += string(r)
	}
	if word != "" {
		if line == "" {
			line = word
		} else if textWidth(line+" "+word) <= maxPx {
			line += " " + word
		} else {
			flush()
			line = word
		}
	}
	flush()
	return out
}
