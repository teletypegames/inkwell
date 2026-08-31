package inkwell

type orCond struct{ cs []Condition }

// Or holds when any one of its conditions does. No conditions at all does
// not hold.
func Or(cs ...Condition) Condition { return &orCond{cs: cs} }
func (o *orCond) Eval(ctx *Ctx) bool {
	for _, c := range o.cs {
		if c.Eval(ctx) {
			return true
		}
	}
	return false
}
