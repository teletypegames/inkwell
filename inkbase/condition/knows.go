package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type knowsCond struct{ entity inkbase.Knowable }

// Knows is true while the player has met the entity.
func Knows(entity inkbase.Knowable) inkbase.Condition { return &knowsCond{entity: entity} }
func (k *knowsCond) Eval(ctx *inkbase.Ctx) bool       { return ctx.State.Known(k.entity) > 0 }
