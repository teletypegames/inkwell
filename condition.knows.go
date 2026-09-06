package inkwell

type knowsCond struct{ entity Knowable }

// Knows is true once the player has come to know an entity — see
// State.Known — which is what a dialog choice that asks about somebody by
// name wants to be gated on.
func Knows(entity Knowable) Condition { return &knowsCond{entity: entity} }

func (c *knowsCond) Eval(ctx *Ctx) bool { return ctx.Game.State.Known(c.entity) > 0 }
