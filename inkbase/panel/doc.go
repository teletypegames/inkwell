// Package panel is a stand-alone inkbase widget: [Panel], a coloured
// rectangle with an optional border, used as a backdrop for grouping other
// widgets. With Blocks set it is opaque to clicks — what a plate the HUD is
// mounted on wants; a rule drawn on top of one does not need it.
//
// The widget reads its default colours from a host that implements
// PanelColors() [Colors] on the UICtx's Game. inkwell has its own Panel
// widget, drawn with the game's theme; this one is for a game built on
// inkbase alone.
package panel
