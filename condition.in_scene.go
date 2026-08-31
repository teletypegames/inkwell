package inkwell

type inSceneCond struct{ name string }

// InScene holds while the player is in the named scene.
func InScene(name string) Condition       { return &inSceneCond{name: name} }
func (i *inSceneCond) Eval(ctx *Ctx) bool { return ctx.Scene != nil && ctx.Scene.Name == i.name }
