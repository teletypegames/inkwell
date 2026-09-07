package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type gotoAction struct{ scene string }

// GoTo takes the player to another scene, through the scene fade.
// The scene is identified by name; the root wrapper converts *Scene to name.
func GoTo(scene string) inkbase.Action          { return &gotoAction{scene: scene} }
func (a *gotoAction) Start() inkbase.Runner     { return a }
func (a *gotoAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(SceneHost)
	if a.scene != "" {
		host.ChangeScene(a.scene)
	}
	return inkbase.StatusDone
}
