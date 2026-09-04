package inkwell

type gotoAction struct{ scene SceneRef }

// GoTo takes the player to another scene, through the scene fade. It takes
// a SceneRef so that a scene may name one that names it back; see SceneRef.
func GoTo(scene SceneRef) Action    { return &gotoAction{scene: scene} }
func (a *gotoAction) Start() Runner { return a }
func (a *gotoAction) Tick(ctx *Ctx) Status {
	if a.scene != nil {
		ctx.Game.changeScene(a.scene.Scene().Name)
	}
	return StatusDone
}
