package inkwell

type stopMusicAction struct{}

// StopMusic silences the music track.
func StopMusic() Action                         { return &stopMusicAction{} }
func (a *stopMusicAction) Start() Runner        { return a }
func (a *stopMusicAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.Audio.StopMusic()
	return StatusDone
}
