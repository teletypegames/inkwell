package inkwell

type Item struct {
	Name        string
	Sprite      *Sprite
	Description Text
	// Lore is what the player's record says about the item once it is
	// known — a paragraph for a dossier or a codex screen, not the one-line
	// Description the inventory shows. Nil keeps the item out of such a
	// screen.
	Lore Text

	OnUseSelf Action
	// OnUseWith: targetName -> action. Target may be a hotspot Name or another item Name.
	OnUseWith map[*Item]Action
}

func (i Item) GetName() string   { return i.Name }
func (i Item) TypeLabel() string { return "item" }
