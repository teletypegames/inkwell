package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// Say puts one line in a character's mouth. The line is a Text, so the same
// call carries every language it has been written in; which one reaches the
// player is decided when the runner starts, not when the content is
// authored.
func Say(speaker Talker, text Text) Action {
	return act.Say(speaker, text)
}
