package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// EndDialog closes the open conversation, and so lets the RunDialog
// that opened it finish.
func EndDialog() Action { return act.EndDialog() }
