package inkwell

type playSoundAction struct{ sound *Sound }

// PlaySound plays a noise once, over whatever else is playing.
func PlaySound(sound *Sound) Action      { return &playSoundAction{sound: sound} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.Audio.PlaySound(a.sound)
	return StatusDone
}
