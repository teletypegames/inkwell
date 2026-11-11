// Package act holds the adventure actions: the ones that need more of the
// game than the [inkbase.State] in the Ctx — speech, the inventory, dialogs,
// scene changes, walking, poses, the voice setting and the end card. Each
// action type-asserts ctx.Game to the narrow host interface it needs
// ([SpeechHost], [InventoryHost], [DialogHost], [SceneHost], [MovementHost],
// [PoseHost], [VoiceHost], [EndCardHost], [FlashHost]), and an inkwell *Game
// satisfies all of them. Games normally use the re-exports in the root
// inkwell package, which take the entities themselves (a *Scene, a
// *Character) where these take names.
//
// The generic actions and the combinators (Seq, Par, If, Wait, SetFlag, …)
// live in the inkbase/action package, which also describes the Action and
// Runner model and how to write an action of your own.
//
// # The actions
//
//   - [Say] shows a line over the speaker (a speech bubble, or the log alone
//     for a DeliveryLog speaker) and appends it to the chat log. A click
//     skips it. It lasts as long as the line's recorded [inkbase.Voice] when
//     voices are on, otherwise a time that scales with the resolved text's
//     length, 1.2 seconds at least. The speaker becomes known.
//   - [Reply] is Say with the line decided when it is spoken: answer is
//     handed the State and returns the Text for this moment, so a speaker
//     can answer from what the player has done without a dialog node per
//     case.
//   - [Give] adds an item to the inventory and marks it known; [TakeAway]
//     removes it.
//   - [RequireItem] flashes MsgNeedItem and fails — stopping the sequence it
//     is in — unless the item is held.
//   - [RunDialog] opens a dialog and holds the script until it closes;
//     [EndDialog] closes the open one, and [GotoNode] moves it to another of
//     its nodes.
//   - [GoTo] changes scene through the scene fade and holds the script until
//     the new scene is showing, running the old scene's OnLeave and the new
//     scene's OnEnter on the way, so one script can lead the player through
//     several scenes; [Back] returns to the previous scene the same way, and
//     does nothing when there is none.
//   - [Walk] sends a character to a point through the walkboxes and holds
//     the script until it arrives.
//   - [SetPose] holds a character in a pose until it next walks or speaks.
//   - [SetVoice] switches recorded lines on or off and saves the settings.
//   - [Flash] puts a line in the status strip for about two seconds.
//   - [ShowEnd] puts the end card up and never finishes: the game is over.
//
// For example, a hotspot that only opens with the key in hand:
//
//	OnUse: inkwell.Seq(
//		inkwell.RequireItem(key),
//		inkwell.Say(paul, inkwell.Text{{Lang: "en", Text: "Click."}}),
//		inkwell.SetFlag("door_open"),
//	),
package act
