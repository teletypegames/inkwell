package inkwell

type gotoStageAction struct{ stage *Stage }

// GoToStage shows another Stage. The change lands at the top of the next
// frame, along with any Restart or LoadSlot asked for in the same breath, so
// a menu row can say the whole of what it means in one Seq without having to
// know what order the engine would otherwise have applied it in.
func GoToStage(stage *Stage) Action { return &gotoStageAction{stage: stage} }

func (a *gotoStageAction) Start() Runner { return a }

func (a *gotoStageAction) Tick(ctx *Ctx) Status {
	ctx.Game.GoToStage(a.stage)
	return StatusDone
}
