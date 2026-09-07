package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type setVarIfEmptyAction struct {
	name string
	v    any
}

// SetVarIfEmpty writes a world-state variable only while nothing has written it yet.
func SetVarIfEmpty(name string, v any) inkbase.Action { return &setVarIfEmptyAction{name: name, v: v} }

func (a *setVarIfEmptyAction) Start() inkbase.Runner { return a }

func (a *setVarIfEmptyAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if ctx.State.Var(a.name) == nil {
		ctx.State.SetVar(a.name, a.v)
	}
	return inkbase.StatusDone
}
