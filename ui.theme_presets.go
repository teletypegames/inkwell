package inkwell

// RegisterPresetThemes installs the four built-in themes, each of which is a
// file of its own. NewGame calls this and selects "classic-scumm", so a game
// that never mentions a theme still has one.
func RegisterPresetThemes(g *Game) {
	g.ThemeManager.Register(classicScummTheme())
	g.ThemeManager.Register(sierraCoinTheme())
	g.ThemeManager.Register(paperNotebookTheme())
	g.ThemeManager.Register(terminalGreenTheme())
}
