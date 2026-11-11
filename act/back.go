package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type backAction struct{}

// Back returns to the scene the player came from and, like [GoTo], holds
// the script until it is showing. It does nothing when there is none.
func Back() inkbase.Action { return backAction{} }

// Start returns the Runner of one return to the previous scene.
func (a backAction) Start() inkbase.Runner {
	return &enterRunner{scene: func(h SceneHost) string { return h.PreviousScene() }}
}
