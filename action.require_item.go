package inkwell

// RequireItem fails silently with a generic line if the player doesn't have
// the named item selected or in inventory. Used at the top of Use handlers.
type requireItemAction struct{ item string }

func RequireItem(item string) Action       { return &requireItemAction{item: item} }
func (a *requireItemAction) Start() Runner { return a }
func (a *requireItemAction) Tick(ctx *Ctx) Status {
	if ctx.Game.Inventory.Has(a.item) {
		return StatusDone
	}
	ctx.Game.FlashLine(Format(MsgNeedItem, ctx.Game.ItemLabel(a.item)))
	return StatusFailed
}
