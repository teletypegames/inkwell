package inkwell

type hasItemCond struct{ name string }

// HasItem holds while the named item is in the player's inventory, whether
// or not it is the one in hand.
func HasItem(name string) Condition       { return &hasItemCond{name: name} }
func (h *hasItemCond) Eval(ctx *Ctx) bool { return ctx.Game.Inventory.Has(h.name) }
