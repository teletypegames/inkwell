package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// ClearFlag lowers a world-state flag, which does not have to have been
// raised for it to be lowered.
func ClearFlag(name string) Action { return adaptAction(ibaction.ClearFlag(name)) }
