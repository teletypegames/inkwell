package inkwell

// Scene is a place in the fiction: a picture, the things in it a verb can
// reach, the ways out of it and the people standing in it. The engine
// always has one showing on the world Stage.
type Scene struct {
	Name string
	// Title is the human-readable name a TopBar shows; empty shows Name.
	Title Text
	// Layers is the picture, in planes: every one a full-screen image at a
	// depth Z, with the cast at 0. A negative Z is behind every character —
	// the background is -1 — and a positive Z is in front of every
	// character, a foreground the player walks behind. Nothing here leaves
	// the theme's backdrop showing — a scene that is all widgets.
	Layers []SceneLayer
	// Music starts when the scene is entered and keeps playing until a
	// scene or Stage with different music is. Nil changes nothing.
	Music     *Music
	Hotspots  []Hotspot
	Exits     []Exit
	Walkboxes []Polygon
	Triggers  []Trigger
	Actors    []SceneActor
	OnEnter   Action
	OnLeave   Action

	// Derived is asked once, by Run, for exits the scene computes rather
	// than writes — a map screen's pins, one per scene that leads back to
	// it. It is handed the scene itself, so it can ask "who leads to me"
	// without naming its own variable, which its own initialiser may not.
	// What comes back is appended to Exits before Validate.
	Derived func(g *Game, self *Scene) []Exit
}

func (s Scene) GetName() string   { return s.Name }
func (s Scene) TypeLabel() string { return "scene" }

// Leads appends exits to a scene after it is declared, and hands the scene
// back. Two scenes that lead to each other cannot each hold the other's
// pointer in their own literal — Go rejects the initialisation cycle — so a
// scene file declares the scene first and its exits second:
//
//	var Alley = SceneManager.Register(&Scene{Name: "alley", …})
//
//	var _ = Alley.Leads(
//		Exit{To: Street, Side: ExitLeft},
//	)
//
// The second declaration depends on both scenes and neither scene depends
// on it, so the graph may be as tangled as the city is.
func (s *Scene) Leads(exits ...Exit) *Scene {
	s.Exits = append(s.Exits, exits...)
	return s
}

// SceneActor is a character standing in a scene when it is entered.
type SceneActor struct {
	Character *Character
	At        Point
}
