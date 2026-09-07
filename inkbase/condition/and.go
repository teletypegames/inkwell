package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type andCond struct{ cs []inkbase.Condition }

// And holds when every one of its conditions does.
func And(cs ...inkbase.Condition) inkbase.Condition { return &andCond{cs: cs} }
func (a *andCond) Eval(ctx *inkbase.Ctx) bool {
	for _, c := range a.cs {
		if !c.Eval(ctx) {
			return false
		}
	}
	return true
}
