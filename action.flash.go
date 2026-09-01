package inkwell

type flashAction struct{ text Text }

// Flash puts a line in the status strip for about two seconds — the
// engine's own way of telling the player something happened that no
// character would say out loud. Game.FlashLine is the same thing for code
// that is not an action.
func Flash(text Text) Action { return &flashAction{text: text} }

func (a *flashAction) Start() Runner { return a }

func (a *flashAction) Tick(ctx *Ctx) Status {
	ctx.Game.FlashLine(a.text)
	return StatusDone
}
