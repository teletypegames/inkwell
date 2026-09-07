package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// LifecycleHost provides restart and quit.
type LifecycleHost interface {
	Restart()
}

type restartAction struct{}

// Restart begins a new session.
func Restart() inkbase.Action { return &restartAction{} }

func (a *restartAction) Start() inkbase.Runner { return a }

func (a *restartAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(LifecycleHost); ok {
		h.Restart()
	}
	return inkbase.StatusDone
}
