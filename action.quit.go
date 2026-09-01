package inkwell

type quitAction struct{}

// Quit closes the window at the end of the frame. Run returns nil, so a game
// that wants to write anything on the way out can do it after Run.
func Quit() Action { return &quitAction{} }

func (a *quitAction) Start() Runner { return a }

func (a *quitAction) Tick(ctx *Ctx) Status {
	ctx.Game.quitting = true
	return StatusDone
}
