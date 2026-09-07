package inkwell

import ibcond "git.teletypegames.org/engines/inkwell/inkbase/condition"

// Knows is true once the player has come to know an entity.
func Knows(entity Knowable) Condition { return adaptCondition(ibcond.Knows(entity)) }
