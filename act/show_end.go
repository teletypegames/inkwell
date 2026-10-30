package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type showEndAction struct{ text inkbase.Text }

// ShowEnd puts the end card up and never finishes: the game is over.
func ShowEnd(text inkbase.Text) inkbase.Action { return &showEndAction{text: text} }
func (a *showEndAction) Start() inkbase.Runner { return &showEndRunner{spec: a} }

type showEndRunner struct {
	spec    *showEndAction
	started bool
}

func (r *showEndRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(EndCardHost)
	if !r.started {
		r.started = true
		host.ShowEndCard(r.spec.text.String())
	}
	return inkbase.StatusRunning
}
