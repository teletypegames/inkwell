package inkwell

type inSceneCond struct{ scene *Scene }

// InScene is true while the given scene is the current one.
func InScene(scene *Scene) Condition      { return &inSceneCond{scene: scene} }
func (i *inSceneCond) Eval(ctx *Ctx) bool { return ctx.Scene != nil && ctx.Scene == i.scene }
