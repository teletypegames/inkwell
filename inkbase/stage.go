package inkbase

// Stage is one full-screen thing the game can be showing, and exactly one of
// them is showing at a time. A main menu is a Stage, a settings screen is a
// Stage, a credits card is a Stage — and so is the game itself, the one with
// World set.
type Stage struct {
	Name  string
	Title Text

	// World says this Stage is the game: the current scene is painted, its
	// actors walk, its triggers fire, and a click reaches its hotspots.
	World bool

	// Background is a full-screen image for a Stage that is not the world.
	Background *Asset

	// Music plays while this Stage shows.
	Music *Music

	// Hold shows this Stage for a number of seconds and then moves on to
	// Next. 0 stays until something says GoToStage.
	Hold float64
	Next *Stage

	// Skippable lets the player cut a Hold short.
	Skippable bool

	OnEnter Action
	OnLeave Action
}

// GetName returns the stage's Name, its key in the registry.
func (s Stage) GetName() string { return s.Name }

// TypeLabel returns "stage", the kind used in log and panic messages.
func (s Stage) TypeLabel() string { return "stage" }
