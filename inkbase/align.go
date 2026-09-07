package inkbase

// Align controls horizontal text alignment for widgets that draw a single line.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Size is a {Width, Height} pair in screen-space pixels.
type Size struct {
	W, H int
}

// MouseButton identifies which mouse button a widget binds to.
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
)
