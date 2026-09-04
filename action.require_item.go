package inkwell

type requireItemAction struct{ item *Item }

// RequireItem fails the sequence, with a flash line, unless the item is held.
func RequireItem(item *Item) Action        { return &requireItemAction{item: item} }
func (a *requireItemAction) Start() Runner { return a }
func (a *requireItemAction) Tick(ctx *Ctx) Status {
	if ctx.Game.Inventory.Has(a.item.Name) {
		return StatusDone
	}
	ctx.Game.FlashLine(Format(MsgNeedItem, ctx.Game.ItemLabel(a.item.Name)))
	return StatusFailed
}
