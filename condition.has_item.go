package inkwell

import "git.teletypegames.org/engines/inkwell/cond"

// HasItem is true while the item is in the inventory.
func HasItem(item *Item) Condition { return cond.HasItem(item.Name) }
