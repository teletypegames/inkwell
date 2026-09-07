package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

// StageReader provides the current stage for conditions.
type StageReader interface {
	CurrentStage() *inkbase.Stage
	OnWorldStage() bool
}

type onStageCond struct{ stage *inkbase.Stage }

// OnStage is true while the given Stage is the one showing.
func OnStage(stage *inkbase.Stage) inkbase.Condition { return &onStageCond{stage: stage} }
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
