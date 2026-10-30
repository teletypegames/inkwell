package cond

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/item"
)

// Given is a Condition written as a literal: every field that is set has to
// hold, and an empty Given holds always. It is the shape a dialog choice or a
// widget reads best in -- one struct, the requirements stacked -- where the
// combinators would spell the same thing as a call tree.
type Given struct {
	// Flag must be raised; NotFlag must not be.
	Flag    string
	NotFlag string
	// Item must be in the inventory; NotItem must not be.
	Item    *item.Item
	NotItem *item.Item
	// Var must compare equal to Is, by Go's own ==. Is alone means nothing;
	// Var alone compares the variable to nil.
	Var string
	Is  any
}

// Eval reports whether every set field of given holds in ctx: the flags
// raised or not on ctx.State, the items in the inventory or not, and the
// variable equal to Is. An Item or NotItem condition is false when
// ctx.Game cannot answer [InventoryChecker].
func (given Given) Eval(ctx *inkbase.Ctx) bool {
	state := ctx.State
	if given.Flag != "" && !state.Flag(given.Flag) {
		return false
	}
	if given.NotFlag != "" && state.Flag(given.NotFlag) {
		return false
	}
	if given.Item != nil {
		if host, ok := ctx.Game.(InventoryChecker); ok {
			if !host.InventoryHas(given.Item.Name) {
				return false
			}
		} else {
			return false
		}
	}
	if given.NotItem != nil {
		if host, ok := ctx.Game.(InventoryChecker); ok {
			if host.InventoryHas(given.NotItem.Name) {
				return false
			}
		} else {
			return false
		}
	}
	if given.Var != "" && state.Var(given.Var) != given.Is {
		return false
	}
	return true
}
