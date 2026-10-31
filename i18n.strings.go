package inkwell

import (
	"git.teletypegames.org/engines/inkwell/act"
	"git.teletypegames.org/engines/inkwell/inkbase/menu"
	"git.teletypegames.org/engines/inkwell/inkbase/topbar"
)

// The engine's own voice: the handful of sentences inkwell puts on screen
// without a game having authored them. A game's own wording lives in its own
// content files — these are the ones that belong to the machinery, and they
// are gathered here rather than left at their call sites because a translator
// adding a language wants one file for the chrome, not a grep.
//
// The ones a sub-package prints itself — RequireItem's line, the TopBar's
// score label and the Menu's chrome — are written once, in that package;
// the variables here start as those and the engine copies them back every
// frame, so a game changes them here, like the rest.
//
// They are variables, not constants, so a game that speaks a language the
// engine has not been taught can teach it one:
//
//	inkwell.MsgNothing = append(inkwell.MsgNothing, inkwell.Translation{
//		Lang: "de",
//		Text: "Nichts von Bedeutung.",
//	})
var (
	// MsgNothing is the answer to a verb nothing is bound to.
	MsgNothing = Text{
		{Lang: "en", Text: "Nothing worth mentioning."},
		{Lang: "hu", Text: "Semmi említésre méltó."},
	}

	// MsgNotThat is the fallback when an item is used on something nobody
	// paired it with and Game.UseWithFail is unset.
	MsgNotThat = Text{
		{Lang: "en", Text: "Not that."},
		{Lang: "hu", Text: "Nem ehhez."},
	}

	// MsgNeedItem is what RequireItem says. %s is the item. It starts as
	// [act.MsgNeedItem], and the engine hands it back to act every frame.
	MsgNeedItem = act.MsgNeedItem

	// HintVerb is the status line over a plain target: verb, then target.
	HintVerb = Text{
		{Lang: "en", Text: "%s %s"},
		{Lang: "hu", Text: "%s %s"},
	}

	// HintVerbOn is the status line while an item is held: verb, item,
	// target. The preposition is the reason this is a sentence and not a
	// join — English puts it between the two, Hungarian after them.
	HintVerbOn = Text{
		{Lang: "en", Text: "%s %s on %s"},
		{Lang: "hu", Text: "%s %s ezen: %s"},
	}

	// ScoreLabel prefixes the TopBar score. It starts as
	// [topbar.ScoreLabel], and the engine hands it back to topbar every
	// frame.
	ScoreLabel = topbar.ScoreLabel

	// MenuEmpty is what a Menu page with no rows says; it starts as
	// [menu.MenuEmpty], and the engine hands it back to menu every frame.
	MenuEmpty = menu.MenuEmpty

	// MenuKeys is the key legend along the foot of a Menu window; it starts
	// as [menu.MenuKeys], and the engine hands it back to menu every frame.
	MenuKeys = menu.MenuKeys

	// MenuKeysSticky is the same legend for a Menu that is the screen: it
	// cannot be closed, so it does not offer to be. It starts as
	// [menu.MenuKeysSticky], and the engine hands it back to menu every
	// frame.
	MenuKeysSticky = menu.MenuKeysSticky
)

// syncStrings hands the root package's copies of the strings the
// sub-packages print back to them, so that a game that rewrote or extended
// inkwell.MenuKeys, say, is read by the Menu that draws it. Like syncLang,
// it runs at the top of every frame and once at boot.
func syncStrings() {
	act.MsgNeedItem = MsgNeedItem
	topbar.ScoreLabel = ScoreLabel
	menu.MenuEmpty = MenuEmpty
	menu.MenuKeys = MenuKeys
	menu.MenuKeysSticky = MenuKeysSticky
}
