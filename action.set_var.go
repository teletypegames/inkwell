package inkwell

type setVarAction struct {
	name string
	v    any
}

// SetVar writes a world-state variable, which is how a game keeps a value
// a flag cannot hold: a score, a time of day, the mode the HUD is in.
func SetVar(name string, v any) Action { return &setVarAction{name: name, v: v} }
func (a *setVarAction) Start() Runner  { return a }
func (a *setVarAction) Tick(ctx *Ctx) Status {
	ctx.Game.State.SetVar(a.name, a.v)
	return StatusDone
}
