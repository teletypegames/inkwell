package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type immediateAction struct {
	fn func(*inkbase.Ctx) inkbase.Status
}

func (a *immediateAction) Start() inkbase.Runner { return a }
func (a *immediateAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if a.fn == nil {
		return inkbase.StatusDone
	}
	return a.fn(ctx)
}

// Custom wraps a user function as an Action.
func Custom(fn func(*inkbase.Ctx) inkbase.Status) inkbase.Action { return &immediateAction{fn: fn} }
