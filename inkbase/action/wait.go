package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type waitAction struct{ seconds float64 }
type waitRunner struct {
	spec    *waitAction
	elapsed float64
}

// Wait holds the frame for a number of seconds.
func Wait(seconds float64) inkbase.Action   { return &waitAction{seconds: seconds} }
func (a *waitAction) Start() inkbase.Runner { return &waitRunner{spec: a} }
func (r *waitRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	r.elapsed += ctx.DT
	if r.elapsed >= r.spec.seconds {
		return inkbase.StatusDone
	}
	return inkbase.StatusRunning
}
