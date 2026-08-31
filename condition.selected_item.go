package inkwell

type selectedItemCond struct{ name string }

// SelectedItem holds while the named item is the one in the player's hand.
func SelectedItem(name string) Condition       { return &selectedItemCond{name: name} }
func (s *selectedItemCond) Eval(ctx *Ctx) bool { return ctx.Game.Inventory.Selected() == s.name }
