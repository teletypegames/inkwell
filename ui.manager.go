package inkwell

// UIManager registers Widget instances. Same shape as every other manager.
type UIManager = Manager[Widget]

// reversedWidgets iterates a manager's contents in reverse registration
// order — used by the engine for top-down input dispatch (the widget
// drawn on top gets the click first).
func reversedWidgets(m *UIManager) []Widget {
	all := m.All()
	out := make([]Widget, len(all))
	for i, w := range all {
		out[len(all)-1-i] = w
	}
	return out
}

// orderedWidgets iterates in registration order — bottom-up draw.
func orderedWidgets(m *UIManager) []Widget { return m.All() }
