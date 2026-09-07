package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type ifAction struct {
	cond inkbase.Condition
	then inkbase.Action
	els  inkbase.Action
}

// If runs then when cond holds and els otherwise. The branch is chosen once,
// on the frame the action starts.
func If(cond inkbase.Condition, then inkbase.Action, els ...inkbase.Action) inkbase.Action {
	var e inkbase.Action
	if len(els) > 0 {
		e = Seq(els...)
	}
	return &ifAction{cond: cond, then: then, els: e}
}

func (a *ifAction) Start() inkbase.Runner { return &ifRunner{spec: a} }

type ifRunner struct {
	spec    *ifAction
	started bool
	inner   inkbase.Runner
}

func (r *ifRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if !r.started {
		r.started = true
		take := r.spec.then
		if r.spec.cond == nil || !r.spec.cond.Eval(ctx) {
			take = r.spec.els
		}
		if take == nil {
			return inkbase.StatusDone
		}
		r.inner = take.Start()
	}
	if r.inner == nil {
		return inkbase.StatusDone
	}
	return r.inner.Tick(ctx)
}
