package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// ShowEnd puts the end card up and never finishes: the game is over, and
// the player closes the window.
func ShowEnd(text Text) Action { return act.ShowEnd(text) }
