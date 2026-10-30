package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type endDialogAction struct{}

// EndDialog closes the open conversation.
func EndDialog() inkbase.Action                  { return &endDialogAction{} }
func (a *endDialogAction) Start() inkbase.Runner { return a }
func (a *endDialogAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(DialogHost)
	host.EndDialog()
	return inkbase.StatusDone
}
