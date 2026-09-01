package inkwell

type restartAction struct{}

// Restart begins a new session: fresh State (carrying Game.Flags forward),
// empty inventory, every character back where the scene puts them, the start
// scene entered and the start script queued. It is what a "New Game" row does.
//
// Like GoToStage the reset lands at the top of the next frame rather than
// inside the action that asked for it, which is what lets a menu row be
// Seq(Restart(), GoToStage(…)) and mean it.
func Restart() Action { return &restartAction{} }

func (a *restartAction) Start() Runner { return a }

func (a *restartAction) Tick(ctx *Ctx) Status {
	ctx.Game.Restart()
	return StatusDone
}
