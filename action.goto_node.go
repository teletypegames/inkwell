package inkwell

type gotoNodeAction struct{ node string }

// GotoNode moves the open conversation to another of its nodes.
func GotoNode(node string) Action       { return &gotoNodeAction{node: node} }
func (a *gotoNodeAction) Start() Runner { return a }
func (a *gotoNodeAction) Tick(ctx *Ctx) Status {
	ctx.Game.gotoDialogueNode(a.node)
	return StatusDone
}
