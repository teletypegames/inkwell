package condition

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/stage"
)

// StageReader provides the current stage for conditions.
type StageReader interface {
	CurrentStage() *stage.Stage
	OnWorldStage() bool
}

type onStageCond struct{ stage *stage.Stage }

// OnStage is true while the given Stage is the one showing.
func OnStage(s *stage.Stage) inkbase.Condition { return &onStageCond{stage: s} }
func (o *onStageCond) Eval(ctx *inkbase.Ctx) bool {
	if h, ok := ctx.Game.(StageReader); ok {
		return h.CurrentStage() == o.stage
	}
	return false
}

type onWorldCond struct{}

// OnWorld is true while the Stage showing is the world.
func OnWorld() inkbase.Condition { return &onWorldCond{} }
func (o *onWorldCond) Eval(ctx *inkbase.Ctx) bool {
	if h, ok := ctx.Game.(StageReader); ok {
		return h.OnWorldStage()
	}
	return false
}
