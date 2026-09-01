package inkwell

type setVarIfEmptyAction struct {
	name string
	v    any
}

// SetVarIfEmpty writes a world-state variable only while nothing has
// written it yet. It is what a scene's OnEnter wants for a value the
// player is meant to see the first time they arrive and never again —
// re-entering the scene must not undo whatever has happened since.
func SetVarIfEmpty(name string, v any) Action { return &setVarIfEmptyAction{name: name, v: v} }

func (a *setVarIfEmptyAction) Start() Runner { return a }

func (a *setVarIfEmptyAction) Tick(ctx *Ctx) Status {
	if ctx.Game.State.Var(a.name) == nil {
		ctx.Game.State.SetVar(a.name, a.v)
	}
	return StatusDone
}
