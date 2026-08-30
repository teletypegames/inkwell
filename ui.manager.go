package inkwell

import "sort"

// WidgetManager registers Widget instances. Same shape as every other manager.
type WidgetManager = Manager[Widget]

// reversedWidgets iterates from the top layer down — used by the engine
// for top-down input dispatch (the widget drawn on top gets the click
// first).
func reversedWidgets(m *WidgetManager) []Widget {
	ordered := orderedWidgets(m)
	out := make([]Widget, len(ordered))
	for i, w := range ordered {
		out[len(ordered)-1-i] = w
	}
	return out
}

// orderedWidgets iterates bottom-up draw order: by layer, and by
// registration order inside a layer.
func orderedWidgets(m *WidgetManager) []Widget {
	all := m.All()
	sort.SliceStable(all, func(i, j int) bool {
		return LayerOf(all[i]) < LayerOf(all[j])
	})
	return all
}
