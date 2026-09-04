package inkwell

type runDialogAction struct{ dialog *Dialog }

// RunDialog opens a dialog tree and holds the script until it closes.
func RunDialog(dialog *Dialog) Action    { return &runDialogAction{dialog: dialog} }
func (a *runDialogAction) Start() Runner { return &runDialogRunner{spec: a} }

type runDialogRunner struct {
	spec    *runDialogAction
	started bool
}

func (r *runDialogRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.startDialog(r.spec.dialog)
	}
	if ctx.Game.dialogActive() {
		return StatusRunning
	}
	return StatusDone
}
