package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// SetVarIfEmpty writes a world-state variable only while nothing has
// written it yet. It is what a scene's OnEnter wants for a value the
// player is meant to see the first time they arrive and never again —
// re-entering the scene must not undo whatever has happened since.
func SetVarIfEmpty(name string, v any) Action { return adaptAction(ibaction.SetVarIfEmpty(name, v)) }
