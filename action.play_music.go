package inkwell

type playMusicAction struct{ music *Music }

// PlayMusic starts a track, looping if it says so; the same track again is
// a no-op.
func PlayMusic(music *Music) Action      { return &playMusicAction{music: music} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlayMusic(a.music)
	return StatusDone
}
