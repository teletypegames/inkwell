package inkwell

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
	Label  Text
	Cursor CursorKind

	OnLook Action
	OnUse  Action
	OnTalk Action
	OnTake Action
	OnGive Action

	// OnUseWith: when the player has an item selected and clicks this hotspot,
	// the engine looks up the item's Name here first.
	OnUseWith map[*Item]Action

	// OnUseWithAny answers for an item OnUseWith does not name — a thing
	// with one reply for whatever you put on it, which is what a slot, a
	// lock or a person usually is. It stands between the authored pairs
	// and Game.UseWithFail, so a hotspot that has something of its own to
	// say says it, and everything else still falls through to the game's
	// general refusal.
	OnUseWithAny Action

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
