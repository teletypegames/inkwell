package inkwell

type takeAwayAction struct{ item *Item }

// TakeAway removes an item from the inventory.
func TakeAway(item *Item) Action        { return &takeAwayAction{item: item} }
func (a *takeAwayAction) Start() Runner { return a }
func (a *takeAwayAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Remove(a.item.Name)
	return StatusDone
}
