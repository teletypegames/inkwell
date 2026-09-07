package inkwell

import ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"

// Wait holds the frame for a number of seconds — the pause between two
// lines of a cutscene.
func Wait(seconds float64) Action { return adaptAction(ibaction.Wait(seconds)) }
