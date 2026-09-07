package inkwell

import ibcond "git.teletypegames.org/engines/inkwell/inkbase/condition"

// VarEq holds while a world-state variable compares equal to a value.
func VarEq(name string, v any) Condition { return adaptCondition(ibcond.VarEq(name, v)) }
