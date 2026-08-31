package inkwell

type runDialogueAction struct{ name string }

// RunDialogue opens a registered dialogue and holds the frame until it
// ends, which is what lets a conversation be one step of a script.
func RunDialogue(name string) Action       { return &runDialogueAction{name: name} }
func (a *runDialogueAction) Start() Runner { return &runDialogueRunner{spec: a} }

type runDialogueRunner struct {
	spec    *runDialogueAction
	started bool
}

func (r *runDialogueRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.startDialogue(r.spec.name)
	}
	if ctx.Game.dialogueActive() {
		return StatusRunning
	}
	return StatusDone
}
