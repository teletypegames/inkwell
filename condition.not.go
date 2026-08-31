package inkwell

type notCond struct{ c Condition }

// Not holds when the condition it wraps does not.
func Not(c Condition) Condition       { return &notCond{c: c} }
func (n *notCond) Eval(ctx *Ctx) bool { return !n.c.Eval(ctx) }
