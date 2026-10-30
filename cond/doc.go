// Package cond holds the adventure-specific conditions: the ones that need
// more of the game than the [inkbase.State] in the Ctx. Each one
// type-asserts ctx.Game to the narrow host interface it requires —
// [InventoryChecker] for [HasItem] and [SelectedItem], [SceneChecker] for
// [InScene], [StageChecker] for [OnStage], [OnWorld] and [Started] — and an
// inkwell *Game satisfies all of them. Games normally use the re-exports in
// the root inkwell package, which take the entity itself (a *Scene, an
// *Item) rather than its name.
//
// The generic conditions and the combinators (Not, And, Or, Flag, VarEq,
// Knows) live in the inkbase/condition package.
//
// # Conditions
//
// A condition is stateless and evaluated lazily, every time it is asked —
// in an If action, in a dialog choice's Show, in a widget's When, in a
// trigger, in an achievement. Nothing is cached, so a condition always
// answers for the game as it stands.
//
// [Given] is the one that is a struct rather than a constructor, and it is
// there for the place a condition is read most: a dialog choice, where "only
// if the player holds the tape and has not been told yet" is two stacked
// fields rather than a call tree:
//
//	Show: inkwell.Given{
//		Item:    mysteryTape,
//		NotFlag: "dex_told",
//	},
//
// Every field set on a Given has to hold — Flag raised, NotFlag not, Item in
// the inventory, NotItem not, Var equal to Is — so a Given is an And of its
// fields, an empty Given holds always, and Or(Given{…}, Given{…}) is still
// how alternatives read.
package cond
