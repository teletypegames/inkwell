package inkwell

// Given is a Condition written as a literal: every field that is set has
// to hold, and an empty Given holds always. It is the shape a dialog
// choice or a widget reads best in — one struct, the requirements
// stacked — where the combinators would spell the same thing as a call
// tree. It is a Condition like any other, so it stands in a Show, an If
// or a When, and Or(Given{…}, Given{…}) is still how alternatives read.
type Given struct {
	// Flag must be raised; NotFlag must not be.
	Flag    string
	NotFlag string
	// Item must be in the inventory; NotItem must not be.
	Item    *Item
	NotItem *Item
	// Var must compare equal to Is, by Go's own ==. Is alone means
	// nothing; Var alone compares the variable to nil.
	Var string
	Is  any
}

func (given Given) Eval(ctx *Ctx) bool {
	state, inventory := ctx.Game.State, ctx.Game.Inventory
	if given.Flag != "" && !state.Flag(given.Flag) {
		return false
	}
	if given.NotFlag != "" && state.Flag(given.NotFlag) {
		return false
	}
	if given.Item != nil && !inventory.Has(given.Item.Name) {
		return false
	}
	if given.NotItem != nil && inventory.Has(given.NotItem.Name) {
		return false
	}
	if given.Var != "" && state.Var(given.Var) != given.Is {
		return false
	}
	return true
}
