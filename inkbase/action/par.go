package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type parAction struct{ children []inkbase.Action }
type parRunner struct {
	runners []inkbase.Runner
	done    []bool
}

// Par starts every action on the same frame and is done when all are.
func Par(actions ...inkbase.Action) inkbase.Action { return &parAction{children: actions} }

func (a *parAction) Start() inkbase.Runner {
	rs := make([]inkbase.Runner, len(a.children))
	for i, c := range a.children {
		rs[i] = c.Start()
	}
	return &parRunner{runners: rs, done: make([]bool, len(rs))}
}

func (r *parRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	allDone := true
	for i, rn := range r.runners {
		if r.done[i] {
			continue
		}
		s := rn.Tick(ctx)
		if s == inkbase.StatusFailed {
			return inkbase.StatusFailed
		}
		if s == inkbase.StatusDone {
			r.done[i] = true
			continue
		}
		allDone = false
	}
	if allDone {
		return inkbase.StatusDone
	}
	return inkbase.StatusRunning
}
