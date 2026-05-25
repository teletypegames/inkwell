package pncdsl

type CursorKind int

const (
	CursorDefault CursorKind = iota
	CursorLook
	CursorUse
	CursorTalk
	CursorTake
	CursorExit
)

type Hotspot struct {
	Name   string
	Area   Shape
	Label  string
	Cursor CursorKind

	OnLook Action
	OnUse  Action
	OnTalk Action
	OnTake Action
	OnGive Action

	// OnUseWith: when the player has an item selected and clicks this hotspot,
	// the engine looks up the item's Name here first.
	OnUseWith map[string]Action

	// OnVerb: custom verbs registered via g.VerbManager.
	OnVerb map[string]Action
}

// handler returns the action bound to verb v for this hotspot, or nil.
func (h Hotspot) handler(v string) Action {
	switch v {
	case "look":
		return h.OnLook
	case "use":
		return h.OnUse
	case "talk":
		return h.OnTalk
	case "take":
		return h.OnTake
	case "give":
		return h.OnGive
	}
	if h.OnVerb != nil {
		return h.OnVerb[v]
	}
	return nil
}
