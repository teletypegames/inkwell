package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type flashAction struct{ text inkbase.Text }

// Flash puts a line in the status strip for about two seconds.
func Flash(text inkbase.Text) inkbase.Action { return &flashAction{text: text} }

func (a *flashAction) Start() inkbase.Runner { return a }

func (a *flashAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(FlashHost)
	host.FlashLine(a.text)
	return inkbase.StatusDone
}
