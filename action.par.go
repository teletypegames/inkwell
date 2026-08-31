package inkwell

type parAction struct{ children []Action }
type parRunner struct {
	runners []Runner
	done    []bool
}

// Par starts every action on the same frame and is done when all of them
// are. The first one to fail fails the whole.
func Par(actions ...Action) Action { return &parAction{children: actions} }

func (a *parAction) Start() Runner {
	rs := make([]Runner, len(a.children))
	for i, c := range a.children {
		rs[i] = c.Start()
	}
	return &parRunner{runners: rs, done: make([]bool, len(rs))}
}

func (r *parRunner) Tick(ctx *Ctx) Status {
	allDone := true
	for i, rn := range r.runners {
		if r.done[i] {
			continue
		}
		s := rn.Tick(ctx)
		if s == StatusFailed {
			return StatusFailed
		}
		if s == StatusDone {
			r.done[i] = true
			continue
		}
		allDone = false
	}
	if allDone {
		return StatusDone
	}
	return StatusRunning
}
