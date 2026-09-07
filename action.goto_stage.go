package inkwell

type gotoStageAction struct{ stage *Stage }

// GoToStage shows another Stage. The change lands at the top of the next
// frame, along with any Restart or LoadSlot asked for in the same breath, so
// a menu row can say the whole of what it means in one Seq without having to
// know what order the engine would otherwise have applied it in.
func GoToStage(stage *Stage) Action { return &gotoStageAction{stage: stage} }

func (a *gotoStageAction) Start() Runner { return a }

func (a *gotoStageAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.GoToStage(a.stage)
	return StatusDone
}

type goBackAction struct{}

// GoBack shows the Stage that was showing before the current one — what a
// screen reached from two places wants on its way out, so that a dossier
// opened over the game returns to the game and one opened from the main
// menu returns to the menu. Nothing happens when there is no previous Stage.
func GoBack() Action { return &goBackAction{} }

func (a *goBackAction) Start() Runner { return a }

func (a *goBackAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	if previous := game.PreviousStage(); previous != nil {
		game.GoToStage(previous)
	}
	return StatusDone
}
