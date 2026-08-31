package inkwell

type seqAction struct{ children []Action }
type seqRunner struct {
	spec    *seqAction
	idx     int
	current Runner
}

// Seq runs its actions one after another and is done when the last one is.
// A Seq handed to another Seq is flattened into it, so composing sequences
// out of sequences costs nothing at runtime.
func Seq(actions ...Action) Action {
	return &seqAction{children: flattenSeq(actions)}
}

func flattenSeq(in []Action) []Action {
	out := make([]Action, 0, len(in))
	for _, a := range in {
		if a == nil {
			continue
		}
		if s, ok := a.(*seqAction); ok {
			out = append(out, s.children...)
		} else {
			out = append(out, a)
		}
	}
	return out
}

func (a *seqAction) Start() Runner { return &seqRunner{spec: a} }

func (r *seqRunner) Tick(ctx *Ctx) Status {
	for {
		if r.idx >= len(r.spec.children) {
			return StatusDone
		}
		if r.current == nil {
			r.current = r.spec.children[r.idx].Start()
		}
		s := r.current.Tick(ctx)
		switch s {
		case StatusDone:
			r.idx++
			r.current = nil
			// keep going only if the just-finished action consumed no time
			// (zero-dt would loop forever otherwise — that's fine here since
			// immediate actions complete in one Tick call).
			continue
		case StatusFailed:
			return StatusFailed
		default:
			return StatusRunning
		}
	}
}
