package inkwell

type hasItemCond struct{ item *Item }

// HasItem is true while the item is in the inventory.
func HasItem(item *Item) Condition        { return &hasItemCond{item: item} }
func (h *hasItemCond) Eval(ctx *Ctx) bool { return ctx.Game.Inventory.Has(h.item.Name) }
