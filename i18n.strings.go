package inkwell

// The engine's own voice: the handful of sentences inkwell puts on screen
// without a game having authored them. A game's own wording lives in its own
// content files — these are the ones that belong to the machinery, and they
// are gathered here rather than left at their call sites because a translator
// adding a language wants one file for the chrome, not a grep.
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

	// MsgNeedItem is what RequireItem says. %s is the item.
	MsgNeedItem = Text{
		{Lang: "en", Text: "I'd need a %s for that."},
		{Lang: "hu", Text: "Ehhez kellene egy %s."},
	}

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

	// ScoreLabel prefixes the TopBar score.
	ScoreLabel = Text{
		{Lang: "en", Text: "Score"},
		{Lang: "hu", Text: "Pont"},
	}

	// MenuEmpty is what a Menu page with no rows says.
	MenuEmpty = Text{
		{Lang: "en", Text: "(nothing here)"},
		{Lang: "hu", Text: "(nincs itt semmi)"},
	}

	// MenuKeys is the key legend along the foot of a Menu window.
	MenuKeys = Text{
		{Lang: "en", Text: "↑↓ move  ENTER open  BACKSPACE back  ESC close"},
		{Lang: "hu", Text: "↑↓ mozgás  ENTER nyit  BACKSPACE vissza  ESC zár"},
	}
)
