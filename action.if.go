package inkwell

type ifAction struct {
	cond Condition
	then Action
	els  Action
}

// If runs then when cond holds and els otherwise. The branch is chosen
// once, on the frame the action starts, and a nil condition takes els.
func If(cond Condition, then Action, els ...Action) Action {
	var e Action
	if len(els) > 0 {
		e = Seq(els...)
	}
	return &ifAction{cond: cond, then: then, els: e}
}

func (a *ifAction) Start() Runner { return &ifRunner{spec: a} }

type ifRunner struct {
	spec    *ifAction
	started bool
	inner   Runner
}

func (r *ifRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		take := r.spec.then
		if r.spec.cond == nil || !r.spec.cond.Eval(ctx) {
			take = r.spec.els
		}
		if take == nil {
			return StatusDone
		}
		r.inner = take.Start()
	}
	if r.inner == nil {
		return StatusDone
	}
	return r.inner.Tick(ctx)
}
