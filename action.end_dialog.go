package inkwell

type endDialogAction struct{}

// EndDialog closes the open conversation, and so lets the RunDialog
// that opened it finish.
func EndDialog() Action                  { return &endDialogAction{} }
func (a *endDialogAction) Start() Runner { return a }
func (a *endDialogAction) Tick(ctx *Ctx) Status {
	ctx.Game.endDialog()
	return StatusDone
}
