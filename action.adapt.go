package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// adaptAction is now identity: inkwell.Action = inkbase.Action.
func adaptAction(inner inkbase.Action) Action { return inner }
