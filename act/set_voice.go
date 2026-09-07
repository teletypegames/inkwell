package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type setVoiceAction struct{ on bool }

// SetVoice switches recorded lines on or off and saves the settings.
func SetVoice(on bool) inkbase.Action { return &setVoiceAction{on: on} }

func (a *setVoiceAction) Start() inkbase.Runner { return a }

func (a *setVoiceAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(VoiceHost)
	host.SetVoiceEnabled(a.on)
	if !a.on {
		host.StopVoice()
	}
	_ = host.SaveSettings()
	return inkbase.StatusDone
}
