// Package statusline is a stand-alone inkbase widget: [StatusLine], a single
// line of text that shows the host's flash message while one is up, and the
// hover hint ("verb + target") otherwise.
//
// The host provides ScreenSize, DrawText, FlashText and SetFlash (the
// message and its remaining time), HoverHint and StatusLineColors() [Colors].
// inkwell has its own StatusLine, wired to the Game; this one is for a game
// built on inkbase alone.
package statusline
