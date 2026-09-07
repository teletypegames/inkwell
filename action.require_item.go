package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// RequireItem fails the sequence, with a flash line, unless the item is held.
func RequireItem(item *Item) Action { return act.RequireItem(item) }
