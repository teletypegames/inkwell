package inkwell

// RegisterPresetThemes installs the four built-in themes, each of which is a
// file of its own. New calls it, so a game that never mentions a theme still
// has one — and it adds rather than replaces: a domain that hands over its
// own ThemeManager keeps everything already in it, and a preset whose name
// the domain has taken is left alone.
func RegisterPresetThemes(g *Game) {
	for _, t := range []Theme{
		classicScummTheme(),
		sierraCoinTheme(),
		paperNotebookTheme(),
		terminalGreenTheme(),
	} {
		if !g.ThemeManager.Has(t.Name) {
			g.ThemeManager.Register(t)
		}
	}
}
