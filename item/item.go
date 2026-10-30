package item

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/sprite"
)

// Item is a thing the player can carry: a name to register and look it up
// by, a picture for the inventory, a line to show for it, and what happens
// when it is used on its own or on something else.
type Item struct {
	// Name is the key the item is registered under and the inventory holds.
	Name string
	// Sprite is the item's picture in the inventory.
	Sprite *sprite.Sprite
	// Description is the item's one-line text: what the inventory and the
	// action log call it, and what the player says on looking at it. An
	// empty Description falls back to Name.
	Description inkbase.Text
	// Lore is what the player's record says about the item once it is
	// known — a paragraph for a dossier or a codex screen, not the one-line
	// Description the inventory shows. Nil keeps the item out of such a
	// screen.
	Lore inkbase.Text

	// OnUseSelf is the action for using the item on its own.
	OnUseSelf inkbase.Action
	// OnUseWith maps another item to the action that runs when the two are
	// combined in the inventory; either item's map may hold the pairing. Using
	// an item on a hotspot reads the hotspot's own OnUseWith instead.
	OnUseWith map[*Item]inkbase.Action
}

// GetName returns the item's Name, the key it is registered under.
func (i Item) GetName() string { return i.Name }

// TypeLabel returns "item", the kind that tells an item apart from
// other entities of the same Name in saved state and in panic messages.
func (i Item) TypeLabel() string { return "item" }
