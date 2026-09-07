package inkwell

type notingAction struct {
	name  string
	value any
	inner Action
}

// Noting writes a world-state variable for exactly as long as the action it
// wraps is running, and clears it afterwards. It is what puts "— paused" in
// a top bar while a conversation holds the game, without anybody having to
// remember to take it back down again: the note is tied to the action's
// lifetime rather than to a pair of calls somebody has to keep in step.
//
// The variable is State, so a widget's condition can read it — which is the
// point. Nothing has to be told that the game is busy; it can ask.
func Noting(name string, value any, inner Action) Action {
	return &notingAction{name: name, value: value, inner: inner}
}

func (a *notingAction) Start() Runner { return &notingRunner{spec: a} }

type notingRunner struct {
	spec  *notingAction
	inner Runner
}

func (r *notingRunner) Tick(ctx *Ctx) Status {
	if r.spec.inner == nil {
		return StatusDone
	}
	if r.inner == nil {
		r.inner = r.spec.inner.Start()
		ctx.State.SetVar(r.spec.name, r.spec.value)
	}
	status := r.inner.Tick(ctx)
	if status != StatusRunning {
		ctx.State.SetVar(r.spec.name, nil)
	}
	return status
}
