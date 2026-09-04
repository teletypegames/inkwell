package inkwell

type selectedItemCond struct{ item *Item }

// SelectedItem is true while the item is the one in hand.
func SelectedItem(item *Item) Condition { return &selectedItemCond{item: item} }
func (s *selectedItemCond) Eval(ctx *Ctx) bool {
	return ctx.Game.Inventory.Selected() == s.item.Name
}
