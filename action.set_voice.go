package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// SetVoice switches recorded lines on or off and saves the settings — the
// row on a settings screen.
func SetVoice(on bool) Action { return act.SetVoice(on) }
