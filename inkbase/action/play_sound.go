package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type playSoundAction struct{ sound *inkbase.Sound }

// PlaySound plays a noise once.
func PlaySound(sound *inkbase.Sound) inkbase.Action { return &playSoundAction{sound: sound} }
func (a *playSoundAction) Start() inkbase.Runner    { return a }
func (a *playSoundAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(AudioHost); ok {
		h.GameAudio().PlaySound(a.sound)
	}
	return inkbase.StatusDone
}
