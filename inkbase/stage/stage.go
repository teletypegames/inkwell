package stage

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/transition"
)

// Stage is one full-screen thing the game can be showing, and exactly one of
// them is showing at a time. A main menu is a Stage, a settings screen is a
// Stage, a credits card is a Stage — and so is the game itself, the one with
// World set: the world is not the thing every other screen floats over, it
// is one of the screens.
//
// Off the world Stage the host paints no world and answers no click in it;
// widgets follow by being gated on a condition that asks which Stage is
// showing (OnStage in the condition package), the same way they hide for a
// cutscene.
type Stage struct {
	// Name is the stage's key in its registry.
	Name string
	// Title is what a menu or a bar calls this screen.
	Title inkbase.Text

	// World says this Stage is the game: the current scene is painted, its
	// actors walk, its triggers fire, and a click reaches its hotspots. A
	// game has at least one; every other Stage is a front end.
	World bool

	// Background is a full-screen image for a Stage that is not the world.
	// Nil leaves the theme's backdrop showing.
	Background *inkbase.Asset

	// Music plays while this Stage shows — a title theme. Nil changes
	// nothing, so a menu over the world keeps the world's track.
	Music *inkbase.Music

	// Transition is how this Stage gives way to the next one: a change out
	// of it — GoToStage, or Hold running out — goes through it, and the
	// swap lands on the frame the picture is fully covered. Nil cuts, which
	// is what a menu and its sub-screens want; a splash card says Fade.
	Transition transition.Transition

	// Hold shows this Stage for a number of seconds and then moves on to
	// Next, through Transition — a splash card, or one card of a title
	// sequence. The clock starts once OnEnter is done. 0 stays until
	// something says GoToStage.
	Hold float64
	// Next is the Stage a Hold moves on to.
	Next *Stage

	// Skippable lets the player cut a Hold short: Space, Enter or a click
	// moves on to Next at once, through the same Transition. A splash card
	// says it; a title card timed to its music does not.
	Skippable bool

	// OnEnter runs when the Stage starts showing.
	OnEnter inkbase.Action
	// OnLeave runs when the Stage gives way to another.
	OnLeave inkbase.Action
}

// GetName returns the stage's Name, its key in the registry.
func (s Stage) GetName() string { return s.Name }

// TypeLabel returns "stage", the kind a Stage is registered as and the word
// used in log and panic messages.
func (s Stage) TypeLabel() string { return "stage" }
