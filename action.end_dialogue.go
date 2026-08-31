package inkwell

type endDialogueAction struct{}

// EndDialogue closes the open conversation, and so lets the RunDialogue
// that opened it finish.
func EndDialogue() Action                  { return &endDialogueAction{} }
func (a *endDialogueAction) Start() Runner { return a }
func (a *endDialogueAction) Tick(ctx *Ctx) Status {
	ctx.Game.endDialogue()
	return StatusDone
}
