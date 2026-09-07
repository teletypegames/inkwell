package act

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/item"
)

type giveAction struct{ item *item.Item }

// Give puts an item in the inventory.
func Give(item *item.Item) inkbase.Action   { return &giveAction{item: item} }
func (a *giveAction) Start() inkbase.Runner { return a }
func (a *giveAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(InventoryHost)
	host.InventoryAdd(a.item.Name)
	ctx.State.NoteKnown(a.item)
	return inkbase.StatusDone
}
