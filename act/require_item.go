package act

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/item"
)

// MsgNeedItem is what RequireItem says when the player hasn't got the item.
// %s is the item label.
var MsgNeedItem = inkbase.Text{
	{Lang: "en", Text: "I'd need a %s for that."},
	{Lang: "hu", Text: "Ehhez kellene egy %s."},
}

type requireItemAction struct{ item *item.Item }

// RequireItem fails the sequence, with a flash line, unless the item is held.
func RequireItem(item *item.Item) inkbase.Action   { return &requireItemAction{item: item} }
func (a *requireItemAction) Start() inkbase.Runner { return a }
func (a *requireItemAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(InventoryHost)
	if host.InventoryHas(a.item.Name) {
		return inkbase.StatusDone
	}
	host.FlashLine(inkbase.Format(MsgNeedItem, host.ItemLabel(a.item.Name)))
	return inkbase.StatusFailed
}
