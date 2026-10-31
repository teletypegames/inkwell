// Package panel is an inkbase widget: [Panel], a coloured rectangle with an
// optional border, used as a backdrop for grouping other widgets. With Blocks
// set it is opaque to clicks (see widget.ClickBlocker) — what a plate the HUD
// is mounted on wants; a rule drawn on top of one does not need it.
//
// The widget reads its default colours from a game that implements [Host]
// (PanelColors() [Colors]) on the UICtx's Game, and draws with its own
// colours, or black, on one that does not. inkwell's *Game is a Host, with
// the theme's PanelBG and DialogBorder, and the root package re-exports
// Panel.
package panel
