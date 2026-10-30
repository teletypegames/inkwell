package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type backAction struct{}

// Back returns to the scene the player came from.
func Back() inkbase.Action                 { return backAction{} }
func (a backAction) Start() inkbase.Runner { return a }
func (a backAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(SceneHost)
	if prev := host.PreviousScene(); prev != "" {
		host.ChangeScene(prev)
	}
	return inkbase.StatusDone
}
