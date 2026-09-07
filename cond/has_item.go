package cond

import "git.teletypegames.org/engines/inkwell/inkbase"

type hasItemCond struct{ name string }

// HasItem is true while the named item is in the inventory.
func HasItem(name string) inkbase.Condition { return &hasItemCond{name: name} }
func (h *hasItemCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(InventoryChecker); ok {
		return host.InventoryHas(h.name)
	}
	return false
}
