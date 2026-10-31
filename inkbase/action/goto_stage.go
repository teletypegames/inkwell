package action

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/stage"
)

// StageHost provides stage navigation.
type StageHost interface {
	GoToStage(s *stage.Stage)
	PreviousStage() *stage.Stage
}

type gotoStageAction struct{ stage *stage.Stage }

// GoToStage shows another Stage.
func GoToStage(s *stage.Stage) inkbase.Action { return &gotoStageAction{stage: s} }

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
