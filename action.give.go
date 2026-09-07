package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// Give puts an item in the inventory.
func Give(item *Item) Action { return act.Give(item) }
