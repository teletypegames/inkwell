package inkwell

type playMusicAction struct{ name string }

// PlayMusic starts the named audio asset as the music track. Naming the
// track that is already playing changes nothing.
func PlayMusic(name string) Action       { return &playMusicAction{name: name} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlayMusic(a.name)
	return StatusDone
}
