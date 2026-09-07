package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type orCond struct{ cs []inkbase.Condition }

// Or holds when any one of its conditions does.
func Or(cs ...inkbase.Condition) inkbase.Condition { return &orCond{cs: cs} }
func (o *orCond) Eval(ctx *inkbase.Ctx) bool {
	for _, c := range o.cs {
		if c.Eval(ctx) {
			return true
		}
	}
	return false
}
