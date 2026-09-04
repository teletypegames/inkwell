package inkwell

type useThemeAction struct{ theme *Theme }

// UseTheme switches the active theme mid-game.
func UseTheme(theme *Theme) Action { return &useThemeAction{theme: theme} }

func (a *useThemeAction) Start() Runner { return a }

func (a *useThemeAction) Tick(ctx *Ctx) Status {
	ctx.Game.UseTheme(a.theme)
	_ = ctx.Game.SaveSettings()
	return StatusDone
}
