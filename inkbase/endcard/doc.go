// Package endcard is an inkbase widget: [EndCard], a full-screen overlay with
// the host's end text centred on it, shown whenever that text is not empty.
// While it is up it swallows every click, so nothing underneath reacts.
//
// The widget reaches its game through [Host]: the screen, the input, the end
// text and the [Colors]. inkwell's *Game is a Host — its end text is what
// ShowEnd put up, its colours the theme's end card tokens — and the root
// package re-exports EndCard; a game built on inkbase alone implements Host
// itself.
package endcard
