package inkwell

type clearFlagAction struct{ name string }

// ClearFlag lowers a world-state flag, which does not have to have been
// raised for it to be lowered.
func ClearFlag(name string) Action       { return &clearFlagAction{name: name} }
func (a *clearFlagAction) Start() Runner { return a }
func (a *clearFlagAction) Tick(ctx *Ctx) Status {
	ctx.Game.State.ClearFlag(a.name)
	return StatusDone
}
