package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// SetVar writes a world-state variable, which is how a game keeps a value
// a flag cannot hold: a score, a time of day, the mode the HUD is in.
func SetVar(name string, v any) Action { return adaptAction(ibaction.SetVar(name, v)) }
