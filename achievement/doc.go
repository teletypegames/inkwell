// Package achievement holds achievements: what the player can be recognised
// for, declared as content.
//
// An [Achievement] is a Name a platform knows it by — the Steam API name — a
// Title and a Text in every language the game speaks, and the
// [inkbase.Condition] that earns it: a flag raised, an item in hand, a scene
// entered. It is registered in a [Manager] like any other entity and handed
// to the game as Game.AchievementManager; Game.Validate rejects one with no
// When.
//
//	AchievementManager.Register(&inkwell.Achievement{
//		Name:  "ACH_LADDER",
//		Title: inkwell.Text{{Lang: "en", Text: "Up and Out"}},
//		Text:  inkwell.Text{{Lang: "en", Text: "Climbed the fire escape."}},
//		When:  inkwell.Flag("on_roof"),
//	})
//
// The engine evaluates every achievement on each frame the world is showing,
// a running script included, and hands the ones that have just come true to
// every Unlocker among the game's extensions. It keeps no record of its own
// beyond a per-process set that stops a condition that stays true from
// firing every frame: a platform that has a record, such as Steam, answers a
// repeated unlock idempotently, and a game without one has nothing to
// persist. A loaded game whose condition is already met therefore reports it
// again in the next process, which is what lets an achievement earned
// offline reach the platform.
//
// Title and Text are what a platform's own listing repeats — the words typed
// into the Steamworks partner site are these — and what a screen of the
// game's own would show; nothing in the engine draws them.
package achievement
