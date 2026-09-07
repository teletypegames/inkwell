package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type varEqCond struct {
	name string
	v    any
}

// VarEq holds while a world-state variable compares equal to a value.
func VarEq(name string, v any) inkbase.Condition { return &varEqCond{name: name, v: v} }
func (e *varEqCond) Eval(ctx *inkbase.Ctx) bool  { return ctx.State.Var(e.name) == e.v }
