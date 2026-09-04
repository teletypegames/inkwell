package inkwell

type gotoAction struct{ scene *Scene }

// GoTo takes the player to another scene, through the scene fade.
func GoTo(scene *Scene) Action      { return &gotoAction{scene: scene} }
func (a *gotoAction) Start() Runner { return a }
func (a *gotoAction) Tick(ctx *Ctx) Status {
	if a.scene != nil {
		ctx.Game.changeScene(a.scene.Name)
	}
	return StatusDone
}
