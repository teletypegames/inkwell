package inkwell

type onStageCond struct{ name string }

// OnStage holds while the named Stage is the one showing. It is what a widget
// belonging to one screen says in its When, and what a menu row that only
// makes sense in the game asks in its Enabled.
func OnStage(name string) Condition       { return &onStageCond{name: name} }
func (o *onStageCond) Eval(ctx *Ctx) bool { return ctx.Game.currentStage == o.name }

type startedCond struct{}

// Started holds once a session exists — the world has been begun or loaded.
// A "Resume" or "Save" row is enabled by this, and at boot it is not.
func Started() Condition                  { return &startedCond{} }
func (s *startedCond) Eval(ctx *Ctx) bool { return ctx.Game.started }
