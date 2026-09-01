package inkwell

type useThemeAction struct{ name string }

// UseTheme switches the active theme from inside a script, which is how a
// game changes its own look at a story beat — the moment the lights go out,
// the moment the machine takes over — rather than only at boot.
func UseTheme(name string) Action { return &useThemeAction{name: name} }

func (a *useThemeAction) Start() Runner { return a }

func (a *useThemeAction) Tick(ctx *Ctx) Status {
	ctx.Game.UseTheme(a.name)
	return StatusDone
}
