// Package letterbox is an inkbase widget: [Letterbox], the pair of bars that
// crop the picture top and bottom while a cutscene plays. It sits on the
// scene layer, over the picture and under the HUD, so a game that also hides
// its HUD for a cutscene gets the whole convention out of two registrations
// sharing one condition. Bar is each bar's height (18 when zero) and Color
// its colour (black when nil).
//
// The widget reaches its game through [Host], which only reports the screen
// size. inkwell's *Game is a Host and the root package re-exports Letterbox.
package letterbox
