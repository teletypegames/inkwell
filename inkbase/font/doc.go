// Package font provides the built-in glyph atlas and the text drawing and
// measuring functions for inkbase widgets.
//
// The atlas is a bitmap font in a fixed [GlyphW] × [GlyphH] (6×16) cell,
// covering ASCII, the Latin-1 supplement, Latin Extended-A, dashes, curly
// quotes, the bullet, the ellipsis, the euro sign and the arrows — enough for
// English, Hungarian and most European languages. A rune outside the
// coverage draws a hollow box, so a gap is visible rather than silent.
//
// Each cell is a white glyph over a baked black drop shadow, so [DrawText]
// tints the glyph with the colour it is given and leaves the shadow black.
// [TextWidth] measures a string, [WrapText] word-wraps it to a pixel width (a
// newline breaks where it stands) and [ClipText] cuts it to a width with an
// ellipsis. All of them take a resolved string, not an inkbase.Text, and
// every vertical measurement in a HUD should derive from GlyphH.
package font
