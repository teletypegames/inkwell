package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// Quitter can be told to quit.
type Quitter interface {
	SetQuitting()
}

type quitAction struct{}

// Quit closes the window at the end of the frame.
func Quit() inkbase.Action { return &quitAction{} }

func (a *quitAction) Start() inkbase.Runner { return a }

func (a *quitAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(Quitter); ok {
		h.SetQuitting()
	}
	return inkbase.StatusDone
}
