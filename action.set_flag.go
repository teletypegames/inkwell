package inkwell

type setFlagAction struct{ name string }

// SetFlag raises a world-state flag — what a Condition asks about and what
// a locked Exit waits for.
func SetFlag(name string) Action              { return &setFlagAction{name: name} }
func (a *setFlagAction) Start() Runner        { return a }
func (a *setFlagAction) Tick(ctx *Ctx) Status { ctx.Game.State.SetFlag(a.name); return StatusDone }
