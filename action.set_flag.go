package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// SetFlag raises a world-state flag — what a Condition asks about and what
// a locked Exit waits for.
func SetFlag(name string) Action { return adaptAction(ibaction.SetFlag(name)) }
