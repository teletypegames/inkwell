package item

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/sprite"
)

type Item struct {
	Name        string
	Sprite      *sprite.Sprite
	Description inkbase.Text
	// Lore is what the player's record says about the item once it is
	// known — a paragraph for a dossier or a codex screen, not the one-line
	// Description the inventory shows. Nil keeps the item out of such a
	// screen.
	Lore inkbase.Text

	OnUseSelf inkbase.Action
	// OnUseWith: targetName -> action. Target may be a hotspot Name or another item Name.
	OnUseWith map[*Item]inkbase.Action
}

func (i Item) GetName() string   { return i.Name }
func (i Item) TypeLabel() string { return "item" }
