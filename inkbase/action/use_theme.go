package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// ThemeHost can switch themes.
type ThemeHost interface {
	UseThemeByName(name string)
	SaveSettings() error
}

type useThemeAction struct{ name string }

// UseTheme switches the active theme by name.
func UseTheme(name string) inkbase.Action { return &useThemeAction{name: name} }

func (a *useThemeAction) Start() inkbase.Runner { return a }

func (a *useThemeAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(ThemeHost); ok {
		h.UseThemeByName(a.name)
		_ = h.SaveSettings()
	}
	return inkbase.StatusDone
}
