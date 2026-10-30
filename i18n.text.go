package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Lang is the language the game speaks — whatever tag its content files
// spell, "en" and "hu" in the games this engine was written for. Empty means
// every Text falls through to Fallback, and then to the language it was
// authored in.
//
// It is a package variable rather than a field on Game because a Text
// resolves itself wherever it is read: inside a widget's Draw, inside the
// Say runner, inside a hover label, inside a hotspot's generated exit.
// Threading a *Game through all of those would put a parameter on every one
// of them for a value there is only ever one of. AssetFS is a package
// variable for the same reason.
//
// The engine syncs Lang and Fallback to inkbase at the top of every frame
// so sub-packages that call Text.String() resolve the same language.
var Lang string

// Fallback is the language used when a Text says nothing in Lang. Empty
// falls through to the first translation the Text lists, which is the
// language it was written in — so a line nobody has translated yet still
// reaches the player rather than vanishing.
var Fallback string

// Translation is one language's wording of one string, with its Voice
// recording if there is one; see [inkbase.Translation].
type Translation = inkbase.Translation

// Text is a string in as many languages as somebody has written it in; it
// resolves against [Lang] and [Fallback] when read. See [inkbase.Text].
type Text = inkbase.Text

// Format composes a Text around other Texts, each part resolving in the
// language of the frame it is read in; see [inkbase.Format].
var Format = inkbase.Format

// plain wraps a string that is not really a translation of anything — a
// registry key standing in for a label nobody wrote. It answers in every
// language, because there is nothing about it to translate.
func plain(s string) Text { return inkbase.Plain(s) }

// syncLang copies the package-level Lang and Fallback into inkbase so
// that sub-package Text.String() reads the same value. Called once per
// frame, and once at boot before the first frame.
func syncLang() {
	inkbase.Lang = Lang
	inkbase.Fallback = Fallback
}
