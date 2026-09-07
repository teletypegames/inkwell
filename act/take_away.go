package act

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/item"
)

type takeAwayAction struct{ item *item.Item }

// TakeAway removes an item from the inventory.
func TakeAway(item *item.Item) inkbase.Action   { return &takeAwayAction{item: item} }
func (a *takeAwayAction) Start() inkbase.Runner { return a }
func (a *takeAwayAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(InventoryHost)
	host.InventoryRemove(a.item.Name)
	return inkbase.StatusDone
}
