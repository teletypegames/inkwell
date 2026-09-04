package inkwell

type setVoiceAction struct{ on bool }

// SetVoice switches recorded lines on or off and saves the settings — the
// row on a settings screen.
func SetVoice(on bool) Action { return &setVoiceAction{on: on} }

func (a *setVoiceAction) Start() Runner { return a }

func (a *setVoiceAction) Tick(ctx *Ctx) Status {
	ctx.Game.settings.Voice = a.on
	if !a.on {
		ctx.Game.Audio.StopVoice()
	}
	_ = ctx.Game.SaveSettings()
	return StatusDone
}
