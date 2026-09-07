package inkbase

// Condition is a question about the game as it stands: a flag, a variable,
// the stage showing. Actions branch on one, widgets are shown by one, and a
// dialog choice is offered by one.
type Condition interface {
	Eval(ctx *Ctx) bool
}
