// Package theme holds the [Theme] type, its [Manager] and the preset themes.
//
// A Theme is a named set of colour tokens every built-in widget draws with:
// the panel and status line, the verb buttons, the inventory slots, the
// speech bubble, the dialog box (whose tokens the menus use too), the end
// card, the cursor and hotspot outline, the scene backdrop, and the top bar,
// chat log and character panel. Widgets read the active theme every frame,
// so a switch shows on the next draw.
//
// # Presets
//
// [Presets] returns the four shipped themes and [RegisterPresets] adds them
// to a manager, leaving alone any name already taken, so a game that hands
// over its own ThemeManager keeps everything in it. inkwell.New does this
// and selects Game.StartTheme, or classic-scumm.
//
//   - classic-scumm: the SCUMM-era dark blue with warm amber accents (the
//     default).
//   - sierra-coin: darker purple and orange, designed for the verb-coin
//     layout.
//   - paper-notebook: cream paper and ink, with a touch of sepia.
//   - terminal-green: retro phosphor green on black.
//
// # Custom themes
//
// A theme is registered like any other entity and selected by pointer:
//
//	noir := g.ThemeManager.Register(&inkwell.Theme{
//		Name:       "midnight-noir",
//		PanelBG:    color.RGBA{6, 6, 12, 255},
//		StatusText: color.White,
//		// … the other colour tokens
//	})
//	g.UseTheme(noir)
//
// The UseTheme action does the same mid-script and saves the setting, so a
// choice on a settings screen is kept. The active theme is remembered by
// Name.
package theme
