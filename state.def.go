package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// State is what the player has done so far: flags, variables, scenes
// visited, dialog nodes talked through and entities known. It is what a
// save slot stores; see [inkbase.State].
type State = inkbase.State

// NewState returns an empty State; see [inkbase.NewState].
func NewState() *State {
	return inkbase.NewState()
}
