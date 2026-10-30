package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type notCond struct{ c inkbase.Condition }

// Not holds when the condition it wraps does not.
func Not(c inkbase.Condition) inkbase.Condition { return &notCond{c: c} }
func (n *notCond) Eval(ctx *inkbase.Ctx) bool   { return !n.c.Eval(ctx) }
