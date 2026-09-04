package inkwell

type Item struct {
	Name        string
	Sprite      *Asset
	Description Text

	OnUseSelf Action
	// OnUseWith: targetName -> action. Target may be a hotspot Name or another item Name.
	OnUseWith map[*Item]Action
}

func (i Item) GetName() string   { return i.Name }
func (i Item) TypeLabel() string { return "item" }
