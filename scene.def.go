package inkwell

// Scene is a place in the fiction: a picture, the things in it a verb can
// reach, the ways out of it and the people standing in it. The engine
// always has one showing on the world Stage.
type Scene struct {
	Name string
	// Title is the human-readable name a TopBar shows; empty shows Name.
	Title Text
	// Background is the picture, stretched over the scene area. Nil leaves
	// the theme's backdrop showing — a scene that is all widgets.
	Background *Asset
	// Music starts when the scene is entered and keeps playing until a
	// scene with different music is. Nil changes nothing.
	Music     *Asset
	Hotspots  []Hotspot
	Exits     []Exit
	Walkboxes []Polygon
	Triggers  []Trigger
	Actors    []SceneActor
	OnEnter   Action
	OnLeave   Action

	// Derived is asked once, by Run, for exits the scene computes rather
	// than writes — a map screen's pins, one per scene that leads back to
	// it. What comes back is appended to Exits before Validate.
	Derived func(g *Game) []Exit
}

func (s Scene) GetName() string   { return s.Name }
func (s Scene) TypeLabel() string { return "scene" }

// Scene makes *Scene its own SceneRef, so a scene that is safe to name
// directly can be.
func (s *Scene) Scene() *Scene { return s }

// SceneRef is a scene named indirectly. Two scenes that lead to each other
// cannot each hold the other's pointer in their own literal — Go rejects
// the initialisation cycle — so an Exit and a GoTo take a SceneRef and ask
// it at run time. *Scene is one; a game's own empty type whose Scene()
// method returns the variable is the other, and breaks the cycle because a
// type is not a dependency.
type SceneRef interface {
	Scene() *Scene
}

// SceneActor is a character standing in a scene when it is entered.
type SceneActor struct {
	Character *Character
	At        Point
}
