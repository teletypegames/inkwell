package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// adaptCondition wraps an inkbase.Condition so it can be evaluated with
// inkwell.Ctx.
func adaptCondition(inner inkbase.Condition) Condition {
	if inner == nil {
		return nil
	}
	return &adaptedCondition{inner: inner}
}

type adaptedCondition struct{ inner inkbase.Condition }

func (c *adaptedCondition) Eval(ctx *Ctx) bool {
	ibctx := &inkbase.Ctx{
		DT:    ctx.DT,
		State: ctx.State,
		Game:  ctx.Game,
	}
	return c.inner.Eval(ibctx)
}
