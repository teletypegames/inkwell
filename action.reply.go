package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// Reply is Say with the line decided when it is spoken: answer is called
// with the State at that moment, and what it returns is what is said.
func Reply(speaker Talker, answer func(state *State) Text) Action {
	return act.Reply(speaker, answer)
}
