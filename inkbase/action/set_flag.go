package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type setFlagAction struct{ name string }

// SetFlag raises a world-state flag.
func SetFlag(name string) inkbase.Action       { return &setFlagAction{name: name} }
func (a *setFlagAction) Start() inkbase.Runner { return a }
func (a *setFlagAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	ctx.State.SetFlag(a.name)
	return inkbase.StatusDone
}
