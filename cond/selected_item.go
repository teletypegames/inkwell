package cond

import "git.teletypegames.org/engines/inkwell/inkbase"

type selectedItemCond struct{ name string }

// SelectedItem is true while the named item is the one in hand.
func SelectedItem(name string) inkbase.Condition { return &selectedItemCond{name: name} }
func (s *selectedItemCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(InventoryChecker); ok {
		return host.InventorySelected() == s.name
	}
	return false
}
