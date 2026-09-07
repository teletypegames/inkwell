package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type seqAction struct{ children []inkbase.Action }
type seqRunner struct {
	spec    *seqAction
	idx     int
	current inkbase.Runner
}

// Seq runs its actions one after another. A Seq handed to another Seq is
// flattened.
func Seq(actions ...inkbase.Action) inkbase.Action {
	return &seqAction{children: flattenSeq(actions)}
}

func flattenSeq(in []inkbase.Action) []inkbase.Action {
	out := make([]inkbase.Action, 0, len(in))
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

func (a *seqAction) Start() inkbase.Runner { return &seqRunner{spec: a} }

func (r *seqRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	for {
		if r.idx >= len(r.spec.children) {
			return inkbase.StatusDone
		}
		if r.current == nil {
			r.current = r.spec.children[r.idx].Start()
		}
		s := r.current.Tick(ctx)
		switch s {
		case inkbase.StatusDone:
			r.idx++
			r.current = nil
			continue
		case inkbase.StatusFailed:
			return inkbase.StatusFailed
		default:
			return inkbase.StatusRunning
		}
	}
}
