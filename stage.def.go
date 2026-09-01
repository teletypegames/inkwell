package inkwell

// Stage is one full-screen thing the game can be showing, and exactly one of
// them is showing at a time. A main menu is a Stage, a settings screen is a
// Stage, a credits card is a Stage — and so is the game itself, which is the
// point: the world is not the thing every other screen floats over, it is one
// of the screens.
//
// That is the difference between a Stage and a Scene. A Scene is a place in
// the fiction, and the engine always has one; walking from a room to an alley
// is a Scene change inside the world. A Stage is what owns the frame. Off the
// world Stage the scene is not drawn, characters do not walk, triggers do not
// fire and hotspots answer nothing — which is what a game showing its main
// menu wants, and what it could not have while the only way to hold the world
// still was a script that never finished.
//
// Widgets need no new field to follow this: OnStage(name) is a Condition like
// any other, so a HUD hides itself off the world the same way it hides itself
// for a cutscene.
type Stage struct {
	Name string
	// Title is what a menu or a bar calls this screen.
	Title Text

	// World says this Stage is the game: the current scene is painted, its
	// actors walk, its triggers fire, and a click reaches its hotspots. A
	// game has at least one; every other Stage is a front end.
	World bool

	// Background is a full-screen image for a Stage that is not the world.
	// Empty leaves the theme's SceneBackdrop showing.
	Background string

	OnEnter Action
	OnLeave Action
}

func (s Stage) GetName() string   { return s.Name }
func (s Stage) TypeLabel() string { return "stage" }
