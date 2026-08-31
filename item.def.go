package inkwell

type Item struct {
	Name        string
	Sprite      string // Asset.Name
	Description Text

	OnUseSelf Action
	// OnUseWith: targetName -> action. Target may be a hotspot Name or another item Name.
	OnUseWith map[string]Action
}

func (i Item) GetName() string   { return i.Name }
func (i Item) TypeLabel() string { return "item" }
