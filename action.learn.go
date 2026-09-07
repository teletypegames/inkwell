package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// Learn marks an entity Known — the one a script has to say out loud, because
// the two rules the engine applies on its own do not reach it: a character
// is known once it has spoken, an item once it has been in the inventory. A
// name read off a list, a face seen through glass, is a Learn.
func Learn(entity Knowable) Action { return adaptAction(ibaction.Learn(entity)) }
