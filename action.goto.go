package inkwell

type gotoAction struct{ scene string }

// GoTo takes the player to another scene, through the engine's fade. An
// unregistered name logs and leaves the player where they stand.
func GoTo(scene string) Action      { return &gotoAction{scene: scene} }
func (a *gotoAction) Start() Runner { return a }
func (a *gotoAction) Tick(ctx *Ctx) Status {
	ctx.Game.changeScene(a.scene)
	return StatusDone
}
