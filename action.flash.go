package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// Flash puts a line in the status strip for about two seconds — the
// engine's own way of telling the player something happened that no
// character would say out loud. Game.FlashLine is the same thing for code
// that is not an action.
func Flash(text Text) Action { return act.Flash(text) }
