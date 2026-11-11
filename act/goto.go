package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type gotoAction struct{ scene string }

// GoTo takes the player to another scene, through the scene fade, and holds
// the script until the new scene is showing: the old scene's OnLeave and
// the new scene's OnEnter run inside it, so the actions after it happen in
// the new scene, after its OnEnter. The scene is identified by name; the
// root wrapper converts *Scene to name. An empty or unknown name does
// nothing.
func GoTo(scene string) inkbase.Action { return gotoAction{scene: scene} }

// Start returns the Runner of one scene change.
func (a gotoAction) Start() inkbase.Runner {
	return &enterRunner{scene: func(SceneHost) string { return a.scene }}
}

// enterRunner asks the host for the scene change on its first tick and
// then ticks the Runner the host hands back.
type enterRunner struct {
	scene func(SceneHost) string
	run   inkbase.Runner
	asked bool
}

func (r *enterRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if !r.asked {
		r.asked = true
		host := ctx.Game.(SceneHost)
		if name := r.scene(host); name != "" {
			r.run = host.EnterScene(name)
		}
	}
	if r.run == nil {
		return inkbase.StatusDone
	}
	return r.run.Tick(ctx)
}
