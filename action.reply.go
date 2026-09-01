package inkwell

// Reply is Say with the line decided when it is spoken rather than when it
// is written: answer is handed the game's State and returns the Text for
// this moment. It is what lets a speaker answer differently as the game
// moves on — a cassette that knows what the player has done, a device that
// reports a reading — without a node per case in the dialog tree.
func Reply(speaker string, answer func(state *State) Text) Action {
	return &sayAction{
		speaker: speaker,
		answer:  answer,
	}
}
