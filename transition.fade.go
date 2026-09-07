package inkwell

import (
	"image/color"

	iwtransition "git.teletypegames.org/engines/inkwell/transition"
)

// Fade is a Transition through black: the old picture darkens over the
// given seconds, the swap lands on the black frame, and the new picture
// brightens over the same span. FadeTo is the same through another colour.
func Fade(seconds float64) Transition                   { return iwtransition.Fade(seconds) }
func FadeTo(seconds float64, to color.Color) Transition { return iwtransition.FadeTo(seconds, to) }
