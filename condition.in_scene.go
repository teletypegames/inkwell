package inkwell

import "git.teletypegames.org/engines/inkwell/cond"

// InScene is true while the given scene is the current one.
func InScene(scene *Scene) Condition { return cond.InScene(scene.Name) }
