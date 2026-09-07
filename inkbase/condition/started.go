package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

// StartedReader reports whether a session has begun.
type StartedReader interface {
	Started() bool
}

type startedCond struct{}

// Started is true once a session exists.
func Started() inkbase.Condition { return &startedCond{} }
func (s *startedCond) Eval(ctx *inkbase.Ctx) bool {
	if h, ok := ctx.Game.(StartedReader); ok {
		return h.Started()
	}
	return false
}
