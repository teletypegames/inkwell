package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type setVarAction struct {
	name string
	v    any
}

// SetVar writes a world-state variable.
func SetVar(name string, v any) inkbase.Action { return &setVarAction{name: name, v: v} }
func (a *setVarAction) Start() inkbase.Runner  { return a }
func (a *setVarAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	ctx.State.SetVar(a.name, a.v)
	return inkbase.StatusDone
}
