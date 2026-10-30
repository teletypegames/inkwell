package inkbase

// Align controls horizontal text alignment for widgets that draw a single line.
type Align int

// The horizontal alignments.
const (
	AlignLeft   Align = iota // flush with the left edge
	AlignCenter              // centered
	AlignRight               // flush with the right edge
)

// Size is a {Width, Height} pair in screen-space pixels.
type Size struct {
	W, H int
}

// MouseButton identifies which mouse button a widget binds to.
type MouseButton int

// The mouse buttons a widget can bind to.
const (
	MouseButtonLeft  MouseButton = iota // the left (primary) button
	MouseButtonRight                    // the right (secondary) button
)
