package inkwell

import (
	"image/color"

	iwtransition "git.teletypegames.org/engines/inkwell/inkbase/transition"
)

// Fade is a Transition through black: the old picture darkens over the
// given seconds, the swap lands on the black frame, and the new picture
// brightens over the same span. See [iwtransition.Fade].
func Fade(seconds float64) Transition { return iwtransition.Fade(seconds) }

// FadeTo is [Fade] through the colour to instead of black; see
// [iwtransition.FadeTo].
func FadeTo(seconds float64, to color.Color) Transition { return iwtransition.FadeTo(seconds, to) }
