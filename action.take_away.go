package inkwell

type takeAwayAction struct{ item string }

// TakeAway drops an item from the player's inventory, and empties the hand
// when that item was the one in it.
func TakeAway(item string) Action       { return &takeAwayAction{item: item} }
func (a *takeAwayAction) Start() Runner { return a }
func (a *takeAwayAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Remove(a.item)
	return StatusDone
}
