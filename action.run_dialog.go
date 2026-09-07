package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// RunDialog opens a dialog tree and holds the script until it closes.
func RunDialog(dialog *Dialog) Action { return act.RunDialog(dialog) }
