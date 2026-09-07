package inkwell

type useLangAction struct{ lang string }

// UseLang switches the language every Text resolves into. A Text resolves
// itself at the moment it is drawn, never when it was authored, so the change
// reaches the whole game on the next frame — the HUD, the log already on
// screen, a line half-spoken.
func UseLang(lang string) Action { return &useLangAction{lang: lang} }

func (a *useLangAction) Start() Runner { return a }

func (a *useLangAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	Lang = a.lang
	_ = game.SaveSettings()
	return StatusDone
}
