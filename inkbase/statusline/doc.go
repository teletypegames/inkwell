// Package statusline is an inkbase widget: [StatusLine], a single line of
// text that shows the host's flash message while one is up, and the hover
// hint ("verb + target") otherwise.
//
// The widget only displays. The flash message and its expiry belong to the
// host: it is up for as long as the host says, whether or not a StatusLine
// is registered, visible or gated off by its When. The widget reaches its
// game through [Host]: the screen, FlashText, HoverHint and the [Colors].
// inkwell's *Game is a Host — its flash is what FlashLine put up, and the
// game loop expires it after about two seconds — and the root package
// re-exports StatusLine.
package statusline
