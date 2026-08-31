package inkwell

type giveAction struct{ item string }

// Give puts an item in the player's inventory. An item already held is
// not added twice.
func Give(item string) Action       { return &giveAction{item: item} }
func (a *giveAction) Start() Runner { return a }
func (a *giveAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Add(a.item)
	return StatusDone
}
