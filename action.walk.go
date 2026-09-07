package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// Walk sends a character to a point through the walkboxes and holds the
// script until it arrives.
func Walk(character *Character, to Point) Action {
	return act.Walk(character.Name, to)
}
