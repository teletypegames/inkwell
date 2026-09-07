package transition

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type fade struct {
	seconds float64
	to      color.Color
}

// Fade is a Transition through black: the old picture darkens over the
// given seconds, the swap lands on the black frame, and the new picture
// brightens over the same span. FadeTo is the same through another colour.
func Fade(seconds float64) Transition { return FadeTo(seconds, color.Black) }

// FadeTo is Fade through the given colour — white for a flashbulb, the
// theme's paper for a page turn.
func FadeTo(seconds float64, to color.Color) Transition {
	return &fade{seconds: seconds, to: to}
}

func (f *fade) Seconds() float64 { return f.seconds }

func (f *fade) Draw(dst *ebiten.Image, width, height int, cover float64) {
	r, g, b, _ := f.to.RGBA()
	c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(cover * 0xffff)}
	vector.DrawFilledRect(dst, 0, 0, float32(width), float32(height), c, false)
}
