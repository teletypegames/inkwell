package inkwell

type backAction struct{}

// Back returns to the scene the player came from. Most exits name their
// destination, because a door leads where it leads; a scene reached from
// several different rooms cannot, so its way out is a direction, not a place.
// A no-op when there is nowhere to go back to.
func Back() Action { return backAction{} }

func (a backAction) Start() Runner { return a }
func (a backAction) Tick(ctx *Ctx) Status {
	if prev := ctx.Game.PreviousScene(); prev != "" {
		ctx.Game.changeScene(prev)
	}
	return StatusDone
}
