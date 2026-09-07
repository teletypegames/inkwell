package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// GoTo takes the player to another scene, through the scene fade.
func GoTo(scene *Scene) Action {
	name := ""
	if scene != nil {
		name = scene.Name
	}
	return act.GoTo(name)
}
