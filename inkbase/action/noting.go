package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type notingAction struct {
	name  string
	value any
	inner inkbase.Action
}

// Noting writes a world-state variable for exactly as long as the action it
// wraps is running, and clears it afterwards.
func Noting(name string, value any, inner inkbase.Action) inkbase.Action {
	return &notingAction{name: name, value: value, inner: inner}
}

func (a *notingAction) Start() inkbase.Runner { return &notingRunner{spec: a} }

type notingRunner struct {
	spec  *notingAction
	inner inkbase.Runner
}

func (r *notingRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if r.spec.inner == nil {
		return inkbase.StatusDone
	}
	if r.inner == nil {
		r.inner = r.spec.inner.Start()
		ctx.State.SetVar(r.spec.name, r.spec.value)
	}
	status := r.inner.Tick(ctx)
	if status != inkbase.StatusRunning {
		ctx.State.SetVar(r.spec.name, nil)
	}
	return status
}
