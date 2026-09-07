package act

import (
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/sprite"
)

type setPoseAction struct {
	character string
	pose      sprite.Pose
}

// SetPose holds a character in a pose until the character next walks or
// speaks. The character is identified by name; the root wrapper converts
// *Character to name.
func SetPose(character string, pose sprite.Pose) inkbase.Action {
	return &setPoseAction{character: character, pose: pose}
}

func (a *setPoseAction) Start() inkbase.Runner { return a }

func (a *setPoseAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(PoseHost)
	host.SetCharacterPose(a.character, a.pose)
	return inkbase.StatusDone
}
