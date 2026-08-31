package inkwell

type showEndAction struct{ text Text }

// ShowEnd puts the end card up and never finishes: the game is over, and
// the player closes the window.
func ShowEnd(text Text) Action         { return &showEndAction{text: text} }
func (a *showEndAction) Start() Runner { return &showEndRunner{spec: a} }

type showEndRunner struct {
	spec    *showEndAction
	started bool
}

func (r *showEndRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.showEndCard(r.spec.text.String())
	}
	return StatusRunning // never finishes; player closes the window
}
