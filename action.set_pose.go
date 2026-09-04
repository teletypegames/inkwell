package inkwell

type setPoseAction struct {
	character *Character
	pose      Pose
}

// SetPose holds a character in a pose — leaning on the wall through a
// cutscene — until the character next walks or speaks, which releases it.
func SetPose(character *Character, pose Pose) Action {
	return &setPoseAction{character: character, pose: pose}
}

func (a *setPoseAction) Start() Runner { return a }

func (a *setPoseAction) Tick(ctx *Ctx) Status {
	if c, ok := ctx.Game.runtimeChar(a.character.Name); ok {
		c.posed = true
		c.pose = a.pose
		c.poseElapsed = 0
	}
	return StatusDone
}
