package widget

import "slices"

// OrderedWidgets sorts widgets bottom-to-top by layer for drawing.
func OrderedWidgets(ws []Widget) []Widget {
	out := slices.Clone(ws)
	slices.SortStableFunc(out, func(a, b Widget) int {
		return int(LayerOf(a)) - int(LayerOf(b))
	})
	return out
}

// ReversedWidgets returns top-to-bottom order for ticking (input).
func ReversedWidgets(ws []Widget) []Widget {
	out := OrderedWidgets(ws)
	slices.Reverse(out)
	return out
}
