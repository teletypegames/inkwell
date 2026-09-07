package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// adaptCondition is now identity: inkwell.Condition = inkbase.Condition.
func adaptCondition(inner inkbase.Condition) Condition { return inner }
