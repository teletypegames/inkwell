package inkwell

type giveAction struct{ item *Item }

// Give puts an item in the inventory.
func Give(item *Item) Action        { return &giveAction{item: item} }
func (a *giveAction) Start() Runner { return a }
func (a *giveAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Add(a.item.Name)
	return StatusDone
}
