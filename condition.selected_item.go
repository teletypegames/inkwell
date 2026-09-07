package inkwell

import "git.teletypegames.org/engines/inkwell/cond"

// SelectedItem is true while the item is the one in hand.
func SelectedItem(item *Item) Condition { return cond.SelectedItem(item.Name) }
