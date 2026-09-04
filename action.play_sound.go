package inkwell

type playSoundAction struct{ asset *Asset }

// PlaySound plays an audio asset once, over whatever else is playing.
func PlaySound(asset *Asset) Action      { return &playSoundAction{asset: asset} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlaySound(a.asset)
	return StatusDone
}
