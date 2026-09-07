package cond

import "git.teletypegames.org/engines/inkwell/inkbase"

type inSceneCond struct{ name string }

// InScene is true while the named scene is the current one.
func InScene(name string) inkbase.Condition { return &inSceneCond{name: name} }
func (i *inSceneCond) Eval(ctx *inkbase.Ctx) bool {
	if host, ok := ctx.Game.(SceneChecker); ok {
		return host.CurrentScene() == i.name
	}
	return false
}
