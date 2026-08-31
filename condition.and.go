package inkwell

type andCond struct{ cs []Condition }

// And holds when every one of its conditions does. No conditions at all
// holds.
func And(cs ...Condition) Condition { return &andCond{cs: cs} }
func (a *andCond) Eval(ctx *Ctx) bool {
	for _, c := range a.cs {
		if !c.Eval(ctx) {
			return false
		}
	}
	return true
}
