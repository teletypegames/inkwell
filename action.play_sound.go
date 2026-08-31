package inkwell

type playSoundAction struct{ name string }

// PlaySound plays the named audio asset once, over whatever music is
// running.
func PlaySound(name string) Action       { return &playSoundAction{name: name} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlaySound(a.name)
	return StatusDone
}
