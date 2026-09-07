package inkwell

import (
	"git.teletypegames.org/engines/inkwell/act"
	ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"
)

// --- delegated to act/ (adventure actions) --------------------------------

// Say puts one line in a character's mouth. The line is a Text, so the same
// call carries every language it has been written in; which one reaches the
// player is decided when the runner starts, not when the content is
// authored.
func Say(speaker Talker, text Text) Action { return act.Say(speaker, text) }

// Reply is Say with the line decided when it is spoken: answer is called
// with the State at that moment, and what it returns is what is said.
func Reply(speaker Talker, answer func(state *State) Text) Action {
	return act.Reply(speaker, answer)
}

// Give puts an item in the inventory.
func Give(item *Item) Action { return act.Give(item) }

// TakeAway removes an item from the inventory.
func TakeAway(item *Item) Action { return act.TakeAway(item) }

// RequireItem fails the sequence, with a flash line, unless the item is held.
func RequireItem(item *Item) Action { return act.RequireItem(item) }

// RunDialog opens a dialog tree and holds the script until it closes.
func RunDialog(dialog *Dialog) Action { return act.RunDialog(dialog) }

// EndDialog closes the open conversation, and so lets the RunDialog
// that opened it finish.
func EndDialog() Action { return act.EndDialog() }

// GotoNode moves the open conversation to another of its nodes.
func GotoNode(node string) Action { return act.GotoNode(node) }

// GoTo takes the player to another scene, through the scene fade.
func GoTo(scene *Scene) Action {
	name := ""
	if scene != nil {
		name = scene.Name
	}
	return act.GoTo(name)
}

// Back returns to the scene the player came from.
func Back() Action { return act.Back() }

// Walk sends a character to a point through the walkboxes and holds the
// script until it arrives.
func Walk(character *Character, to Point) Action {
	return act.Walk(character.Name, to)
}

// SetPose holds a character in a pose until the character next walks or speaks.
func SetPose(character *Character, pose Pose) Action {
	return act.SetPose(character.Name, pose)
}

// SetVoice switches recorded lines on or off and saves the settings.
func SetVoice(on bool) Action { return act.SetVoice(on) }

// ShowEnd puts the end card up and never finishes: the game is over.
func ShowEnd(text Text) Action { return act.ShowEnd(text) }

// Flash puts a line in the status strip for about two seconds.
func Flash(text Text) Action { return act.Flash(text) }

// --- delegated to inkbase/action/ (generic actions) -----------------------

// SetFlag raises a world-state flag.
func SetFlag(name string) Action { return ibaction.SetFlag(name) }

// ClearFlag lowers a world-state flag.
func ClearFlag(name string) Action { return ibaction.ClearFlag(name) }

// SetVar writes a world-state variable.
func SetVar(name string, v any) Action { return ibaction.SetVar(name, v) }

// SetVarIfEmpty writes a world-state variable only while nothing has
// written it yet.
func SetVarIfEmpty(name string, v any) Action { return ibaction.SetVarIfEmpty(name, v) }

// Wait holds the frame for a number of seconds.
func Wait(seconds float64) Action { return ibaction.Wait(seconds) }

// Learn marks an entity Known.
func Learn(entity Knowable) Action { return ibaction.Learn(entity) }

// --- combinators (delegated to inkbase/action/) ---------------------------

// Seq runs its actions one after another.
func Seq(actions ...Action) Action { return ibaction.Seq(actions...) }

// Par runs its actions in parallel and is done when all are.
func Par(actions ...Action) Action { return ibaction.Par(actions...) }

// If runs then when cond holds and els otherwise.
func If(cond Condition, then Action, els ...Action) Action {
	return ibaction.If(cond, then, els...)
}

// Noting writes a var for the duration of the inner action.
func Noting(name string, value any, inner Action) Action {
	return ibaction.Noting(name, value, inner)
}
