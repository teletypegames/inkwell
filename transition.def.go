package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// Transition is how one picture gives way to another: a scene to the next
// scene, a Stage to the next Stage. It runs in two halves — the old picture
// is covered, the swap happens at the moment it is fully hidden, and the new
// picture is uncovered — so whatever a Transition draws, the player never
// sees the frame the swap lands on.
//
// A Transition is a value a Stage carries (Stage.Transition), and nil is a
// cut. Fade is the one the engine ships; a wipe or an iris is another
// implementation of the same two methods.
type Transition interface {
	// Seconds is how long one half takes: covering the old picture, and
	// again uncovering the new one.
	Seconds() float64
	// Draw paints the cover over the frame. cover runs 0 → 1 while the old
	// picture goes and 1 → 0 while the new one comes, so 0 shows the frame
	// as it is and 1 hides it completely.
	Draw(dst *ebiten.Image, width, height int, cover float64)
}
