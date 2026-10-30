// Package condition holds the generic conditions: the ones that need nothing
// but the [inkbase.State] in the Ctx, or at most a narrow host interface any
// game can implement.
//
// A condition is an [inkbase.Condition]: Eval(ctx) bool. It is stateless and
// evaluated lazily, every time it is asked, so it always answers for the
// game as it stands. The constructors:
//
//   - [Not], [And] and [Or] combine conditions.
//   - [Flag] holds while a State flag is raised.
//   - [VarEq] holds while a State var equals a value, by Go's ==.
//   - [Knows] holds once the player has come to know an entity (its
//     State.Known is above zero).
//   - [OnStage] and [OnWorld] ask a [StageReader] which Stage is showing.
//   - [Started] asks a [StartedReader] whether a session exists.
//
// For example, a widget that shows only on the world while a flag is down:
//
//	When: condition.And(condition.OnWorld(), condition.Not(condition.Flag("cutscene")))
//
// The adventure conditions (HasItem, InScene, Given, …) are in the inkwell
// cond package, and the root inkwell package re-exports both sets.
package condition
