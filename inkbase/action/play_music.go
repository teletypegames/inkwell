package action

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/asset"
)

// AudioHost provides access to the audio player.
type AudioHost interface {
	GameAudio() *asset.AudioPlayer
}

type playMusicAction struct{ music *inkbase.Music }

// PlayMusic starts a track.
func PlayMusic(music *inkbase.Music) inkbase.Action { return &playMusicAction{music: music} }
func (a *playMusicAction) Start() inkbase.Runner    { return a }
func (a *playMusicAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(AudioHost); ok {
		h.GameAudio().PlayMusic(a.music)
	}
	return inkbase.StatusDone
}
