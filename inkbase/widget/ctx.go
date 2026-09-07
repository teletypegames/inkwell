package widget

// UICtx is the per-tick context handed to widgets. Game is opaque; widgets
// type-assert it to their own host interface.
type UICtx struct {
	DT   float64
	Game any
}
