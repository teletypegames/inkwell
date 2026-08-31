package inkwell

type varEqCond struct {
	name string
	v    any
}

// VarEq holds while a world-state variable compares equal to a value. The
// comparison is Go's own ==, so the value has to be a comparable type.
func VarEq(name string, v any) Condition { return &varEqCond{name: name, v: v} }
func (e *varEqCond) Eval(ctx *Ctx) bool  { return ctx.Game.State.Var(e.name) == e.v }
