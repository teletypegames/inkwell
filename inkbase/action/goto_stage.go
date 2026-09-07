package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// StageHost provides stage navigation.
type StageHost interface {
	GoToStage(stage *inkbase.Stage)
	PreviousStage() *inkbase.Stage
}

type gotoStageAction struct{ stage *inkbase.Stage }

// GoToStage shows another Stage.
func GoToStage(stage *inkbase.Stage) inkbase.Action { return &gotoStageAction{stage: stage} }

func (a *gotoStageAction) Start() inkbase.Runner { return a }

func (a *gotoStageAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(StageHost); ok {
		h.GoToStage(a.stage)
	}
	return inkbase.StatusDone
}

type goBackAction struct{}

// GoBack shows the Stage that was showing before the current one.
func GoBack() inkbase.Action { return &goBackAction{} }

func (a *goBackAction) Start() inkbase.Runner { return a }

func (a *goBackAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(StageHost); ok {
		if previous := h.PreviousStage(); previous != nil {
			h.GoToStage(previous)
		}
	}
	return inkbase.StatusDone
}
