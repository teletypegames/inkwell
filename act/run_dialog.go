package act

import (
	"git.teletypegames.org/engines/inkwell/dialog"
	"git.teletypegames.org/engines/inkwell/inkbase"
)

type runDialogAction struct{ dialog *dialog.Dialog }

// RunDialog opens a dialog tree and holds the script until it closes.
func RunDialog(dialog *dialog.Dialog) inkbase.Action { return &runDialogAction{dialog: dialog} }
func (a *runDialogAction) Start() inkbase.Runner     { return &runDialogRunner{spec: a} }

type runDialogRunner struct {
	spec    *runDialogAction
	started bool
}

func (r *runDialogRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(DialogHost)
	if !r.started {
		r.started = true
		host.StartDialog(r.spec.dialog)
	}
	if host.DialogActive() {
		return inkbase.StatusRunning
	}
	return inkbase.StatusDone
}
