package inkwell

import ibcond "git.teletypegames.org/engines/inkwell/inkbase/condition"

// Flag holds while a named world-state flag is raised.
func Flag(name string) Condition { return adaptCondition(ibcond.Flag(name)) }
