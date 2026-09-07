package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// SetPose holds a character in a pose — leaning on the wall through a
// cutscene — until the character next walks or speaks, which releases it.
func SetPose(character *Character, pose Pose) Action {
	return act.SetPose(character.Name, pose)
}
