package inkwell

// Condition is a question about the game as it stands: a flag, an item in
// hand, the scene the player is in. Actions branch on one, widgets are shown
// by one, and a dialog choice is offered by one, so the same eight little
// types answer for all three.
type Condition interface {
	Eval(ctx *Ctx) bool
}
