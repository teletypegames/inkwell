package inkwell

// CursorKind is the pointer a Hotspot asks for while the mouse is over it.
type CursorKind int

// The pointers a Hotspot can ask for.
const (
	CursorDefault CursorKind = iota // the ordinary arrow
	CursorLook                      // something to look at
	CursorUse                       // something to use
	CursorTalk                      // someone to talk to
	CursorTake                      // something to pick up
	CursorExit                      // a way out; every Exit gets it
)

// Hotspot is a thing in a scene a verb can reach: an Area to click, a Label
// to show on hover, and the Action each verb runs. A verb with no Action
// does nothing here.
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

	// OnUseWith answers for an item the player has selected when clicking
	// this hotspot, keyed by the item itself; it is looked up first.
	OnUseWith map[*Item]Action

	// OnUseWithAny answers for an item OnUseWith does not name — a thing
	// with one reply for whatever you put on it, which is what a slot, a
	// lock or a person usually is. It stands between the authored pairs
	// and Game.UseWithFail, so a hotspot that has something of its own to
	// say says it, and everything else still falls through to the game's
	// general refusal.
	OnUseWithAny Action

	// OnVerb holds the actions of custom verbs registered in the
	// VerbManager, keyed by the verb's Name.
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
