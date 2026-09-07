package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// TakeAway removes an item from the inventory.
func TakeAway(item *Item) Action { return act.TakeAway(item) }
