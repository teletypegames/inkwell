package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type clearFlagAction struct{ name string }

// ClearFlag lowers a world-state flag.
func ClearFlag(name string) inkbase.Action       { return &clearFlagAction{name: name} }
func (a *clearFlagAction) Start() inkbase.Runner { return a }
func (a *clearFlagAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	ctx.State.ClearFlag(a.name)
	return inkbase.StatusDone
}
