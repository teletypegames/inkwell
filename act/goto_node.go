package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type gotoNodeAction struct{ node string }

// GotoNode moves the open conversation to another of its nodes.
func GotoNode(node string) inkbase.Action          { return &gotoNodeAction{node: node} }
func (a *gotoNodeAction) Start() inkbase.Runner    { return a }
func (a *gotoNodeAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(DialogHost)
	host.GotoDialogNode(a.node)
	return inkbase.StatusDone
}
