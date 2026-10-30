// Package dialog holds dialog trees: the conversations a game opens with
// RunDialog and the DialogBox widget shows.
//
// A [Dialog] is a set of named [DialogNode] values. Each node is a few
// [DialogLine] values, spoken one at a time, followed by the [DialogChoice]
// values the player picks from. A line's Speaker is any [inkbase.Talker] — a
// character or a bodiless speaker — and its Text is an [inkbase.Text], so a
// dialog is written in every language it is played in. The [Manager] holds
// the registered dialogs.
//
//	var Dex = DialogManager.Register(&inkwell.Dialog{
//		Name: "dex",
//		Nodes: []inkwell.DialogNode{{
//			Name:  "hello",
//			Lines: []inkwell.DialogLine{{Speaker: dex, Text: hello}},
//			Choices: []inkwell.DialogChoice{
//				{Text: askTape, Show: inkwell.Given{Item: tape}, Once: true,
//					Actions: []inkwell.Action{inkwell.GotoNode("tape")}},
//				{Text: bye, Actions: []inkwell.Action{inkwell.EndDialog()}},
//			},
//		}, …},
//	})
//
// # The flow
//
//  1. RunDialog starts the dialog at its Start node (the first node when
//     Start is empty); the runner keeps the script running until the dialog
//     closes.
//  2. The DialogBox renders the node's lines one at a time; a click
//     advances.
//  3. After the last line, the choices whose Show holds are offered.
//  4. Picking a choice runs its Actions in sequence. A Once choice is then
//     hidden for good; it is remembered by its node and its position in
//     Choices rather than its wording, so it stays used up across a save
//     and a change of language.
//  5. EndDialog closes the conversation; GotoNode moves it to another node
//     of the same dialog.
package dialog
