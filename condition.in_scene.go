package inkwell

type inSceneCond struct{ scene SceneRef }

// InScene is true while the named scene is the current one.
func InScene(scene SceneRef) Condition { return &inSceneCond{scene: scene} }
func (i *inSceneCond) Eval(ctx *Ctx) bool {
	return ctx.Scene != nil && i.scene != nil && ctx.Scene == i.scene.Scene()
}
