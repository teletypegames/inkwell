package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// holdAction is a Stage's Hold on the queue: it holds the frame for the
// Stage's seconds, then moves on to Next. It is one runner rather than a
// Wait followed by a GoToStage so that a Skippable Stage can end the hold
// early and leave nothing behind — a GoToStage still queued after a skip
// would fire on whatever Stage was showing by then.
type holdAction struct{ stage *Stage }
type holdRunner struct {
	spec    *holdAction
	elapsed float64
}

func (a *holdAction) Start() Runner { return &holdRunner{spec: a} }

func (r *holdRunner) Tick(ctx *Ctx) Status {
	r.elapsed += ctx.DT
	if r.elapsed >= r.spec.stage.Hold || r.skipped(ctx.Game.Input) {
		ctx.Game.GoToStage(r.spec.stage.Next)
		return StatusDone
	}
	return StatusRunning
}

// skipped is the player asking to move on: Space, Enter or a click, each
// claimed so that nothing underneath reads it as its own.
func (r *holdRunner) skipped(input *Input) bool {
	if !r.spec.stage.Skippable {
		return false
	}
	if input.LeftClicked() {
		input.ConsumeLeft()
		return true
	}
	return input.ClaimKey(ebiten.KeySpace) || input.ClaimKey(ebiten.KeyEnter)
}
