package theme

// Presets returns the four built-in themes: classic-scumm, sierra-coin,
// paper-notebook and terminal-green. New calls RegisterPresets so a game
// that never mentions a theme still has one.
func Presets() []*Theme {
	return []*Theme{
		classicScummTheme(),
		sierraCoinTheme(),
		paperNotebookTheme(),
		terminalGreenTheme(),
	}
}

// RegisterPresets installs the built-in themes into a manager, skipping
// any name the caller has already taken.
func RegisterPresets(m *Manager) {
	for _, t := range Presets() {
		if !m.Has(t.Name) {
			m.Register(t)
		}
	}
}
