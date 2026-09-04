package inkwell

type onStageCond struct{ stage *Stage }

// OnStage is true while the given Stage is the one showing.
func OnStage(stage *Stage) Condition      { return &onStageCond{stage: stage} }
func (o *onStageCond) Eval(ctx *Ctx) bool { return ctx.Game.currentStage == o.stage }

type onWorldCond struct{}

// OnWorld is true while the Stage showing is the world — whichever Stage
// that is. A HUD says this rather than naming the game's world Stage, which
// would make the HUD depend on the screen that carries it.
func OnWorld() Condition                  { return &onWorldCond{} }
func (o *onWorldCond) Eval(ctx *Ctx) bool { return ctx.Game.OnWorldStage() }

type startedCond struct{}

// Started is true once a session exists: the world has been begun or loaded.
func Started() Condition                  { return &startedCond{} }
func (s *startedCond) Eval(ctx *Ctx) bool { return ctx.Game.started }
