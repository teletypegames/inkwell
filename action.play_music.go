package inkwell

type playMusicAction struct{ asset *Asset }

// PlayMusic starts an audio asset looping; the same asset again is a no-op.
func PlayMusic(asset *Asset) Action      { return &playMusicAction{asset: asset} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlayMusic(a.asset)
	return StatusDone
}
