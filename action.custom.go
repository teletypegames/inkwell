package inkwell

type immediateAction struct {
	fn func(*Ctx) Status
}

func (a *immediateAction) Start() Runner { return a }
func (a *immediateAction) Tick(ctx *Ctx) Status {
	if a.fn == nil {
		return StatusDone
	}
	return a.fn(ctx)
}

// Custom wraps a user function as an Action. Returns Done on first tick
// unless the function itself returns StatusRunning.
func Custom(fn func(*Ctx) Status) Action { return &immediateAction{fn: fn} }
