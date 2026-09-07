package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// Flasher is the host interface a Flash action needs.
type Flasher interface {
	FlashLine(text inkbase.Text)
}

type flashAction struct{ text inkbase.Text }

// Flash puts a line in the status strip.
func Flash(text inkbase.Text) inkbase.Action { return &flashAction{text: text} }

func (a *flashAction) Start() inkbase.Runner { return a }

func (a *flashAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(Flasher); ok {
		h.FlashLine(a.text)
	}
	return inkbase.StatusDone
}
