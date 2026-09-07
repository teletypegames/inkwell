package cond

import "git.teletypegames.org/engines/inkwell/inkbase"

type onStageCond struct{ stage any }

// OnStage is true while the given Stage is the one showing. The argument is
// a *Stage pointer; identity comparison is used.
func OnStage(stage any) inkbase.Condition { return &onStageCond{stage: stage} }
func (o *onStageCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(StageChecker); ok {
		return host.CurrentStagePtr() == o.stage
	}
	return false
}

type onWorldCond struct{}

// OnWorld is true while the Stage showing is the world.
func OnWorld() inkbase.Condition { return &onWorldCond{} }
func (o *onWorldCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(StageChecker); ok {
		return host.OnWorldStage()
	}
	return false
}

type startedCond struct{}

// Started is true once a session exists: the world has been begun or loaded.
func Started() inkbase.Condition { return &startedCond{} }
func (s *startedCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(StageChecker); ok {
		return host.Started()
	}
	return false
}
