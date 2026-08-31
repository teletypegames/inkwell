package inkwell

type waitAction struct{ seconds float64 }
type waitRunner struct {
	spec    *waitAction
	elapsed float64
}

// Wait holds the frame for a number of seconds — the pause between two
// lines of a cutscene.
func Wait(seconds float64) Action   { return &waitAction{seconds: seconds} }
func (a *waitAction) Start() Runner { return &waitRunner{spec: a} }
func (r *waitRunner) Tick(ctx *Ctx) Status {
	r.elapsed += ctx.DT
	if r.elapsed >= r.spec.seconds {
		return StatusDone
	}
	return StatusRunning
}
