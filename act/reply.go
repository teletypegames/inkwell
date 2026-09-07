package act

import "git.teletypegames.org/engines/inkwell/inkbase"

// Reply is Say with the line decided when it is spoken.
func Reply(speaker inkbase.Talker, answer func(state *inkbase.State) inkbase.Text) inkbase.Action {
	return &sayAction{
		speaker: speaker.TalkerName(),
		answer:  answer,
	}
}
