package inkwell

type runDialogAction struct{ name string }

// RunDialog opens a registered dialog and holds the frame until it
// ends, which is what lets a conversation be one step of a script.
func RunDialog(name string) Action       { return &runDialogAction{name: name} }
func (a *runDialogAction) Start() Runner { return &runDialogRunner{spec: a} }

type runDialogRunner struct {
	spec    *runDialogAction
	started bool
}

func (r *runDialogRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.startDialog(r.spec.name)
	}
	if ctx.Game.dialogActive() {
		return StatusRunning
	}
	return StatusDone
}
