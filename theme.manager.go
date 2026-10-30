package inkwell

import "git.teletypegames.org/engines/inkwell/theme"

// ThemeManager is the registry of a game's Themes, keyed by name; the
// active one is chosen with [Game.UseTheme]. See [theme.Manager].
type ThemeManager = theme.Manager
