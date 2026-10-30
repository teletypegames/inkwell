package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type stopMusicAction struct{}

// StopMusic silences the music track.
func StopMusic() inkbase.Action                  { return &stopMusicAction{} }
func (a *stopMusicAction) Start() inkbase.Runner { return a }
func (a *stopMusicAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(AudioHost); ok {
		h.GameAudio().StopMusic()
	}
	return inkbase.StatusDone
}
