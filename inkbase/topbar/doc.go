// Package topbar is an inkbase widget: [TopBar], a strip across the top of
// the screen in three sections.
//
//   - Left: LeftText, or the scene title, followed after an em dash by the
//     value of the NoteVar State var when it holds anything ("ALLEY —
//     paused"). Notes maps each value to what the player reads, so the
//     var's value stays the same in every language.
//   - Centre: the ScoreVar var, as "Score: X/MAX" when ScoreMax is set and
//     "Score: X" otherwise; ScoreLabel overrides the word, which defaults to
//     the package's [ScoreLabel].
//   - Right: the TimeVar var.
//
// An empty section is skipped. [VarText] is how a var is printed. The strip
// blocks clicks above its bottom edge (widget.ClickBlocker).
//
// The widget reaches its game through [Host]: the screen, SceneTitle,
// StateVar and the [Colors]. inkwell's *Game is a Host, with the theme's top
// bar tokens, and the root package re-exports TopBar; there a game changes
// the default label through inkwell.ScoreLabel, which the engine copies here
// every frame.
package topbar
