package inkwell

// UIManager registers Widget instances. Same shape as every other manager.
type UIManager = Manager[Widget]

// reversedWidgets iterates a manager's contents in reverse registration
// order — used by the engine for top-down input dispatch (the widget
// drawn on top gets the click first).
func reversedWidgets(m *UIManager) []Widget {
	names := m.Names()
	out := make([]Widget, len(names))
	for i, n := range names {
		out[len(names)-1-i] = m.MustGet(n)
	}
	return out
}

// orderedWidgets iterates in registration order — bottom-up draw.
func orderedWidgets(m *UIManager) []Widget {
	names := m.Names()
	out := make([]Widget, len(names))
	for i, n := range names {
		out[i] = m.MustGet(n)
	}
	return out
}
