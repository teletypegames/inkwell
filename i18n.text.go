package inkwell

import "fmt"

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
var Lang string

// Fallback is the language used when a Text says nothing in Lang. Empty
// falls through to the first translation the Text lists, which is the
// language it was written in — so a line nobody has translated yet still
// reaches the player rather than vanishing.
var Fallback string

// Translation is one language's wording of one string.
type Translation struct {
	Lang string
	Text string
	// Voice is the recording of this line in this language, if there is one.
	Voice *Voice
}

// Text is a string in as many languages as somebody has written it in:
//
//	Label: inkwell.Text{
//		{Lang: "en", Text: "locked back door"},
//		{Lang: "hu", Text: "zárt hátsó ajtó"},
//	},
//
// The wording stays where the thing it names is: a hotspot's label is in the
// hotspot, a line of dialog is in the scene that speaks it. Nothing is
// lifted out into a catalogue keyed by identifier, because an identifier is
// one more name to invent and one more thing to keep in step — and because a
// translator reading a scene file can see what the line is answering.
//
// Adding a language is adding a line, and it touches nothing around it. The
// zero value says nothing in any language, so a field nobody has translated
// behaves exactly like the "" it used to be.
type Text []Translation

// plain wraps a string that is not really a translation of anything — a
// registry key standing in for a label nobody wrote. It answers in every
// language, because there is nothing about it to translate.
func plain(s string) Text {
	return Text{{
		Text: s,
	}}
}

// String resolves the Text against Lang: that language if it is there, then
// Fallback, then whichever translation was written first.
func (t Text) String() string { return t.In(Lang) }

// resolved is the Translation String would read from: the one in Lang,
// else Fallback, else the first — with its Voice.
func (t Text) resolved() Translation {
	for _, lang := range []string{Lang, Fallback} {
		if lang == "" {
			continue
		}
		for _, tr := range t {
			if tr.Lang == lang {
				return tr
			}
		}
	}
	if len(t) == 0 {
		return Translation{}
	}
	return t[0]
}

// In is the Text's wording in one named language, falling back the same way
// String does. Composing one Text around another needs it — the frame and
// the part it frames have to come out in the same language — though Format
// is usually the shorter way to say that.
func (t Text) In(lang string) string {
	if len(t) == 0 {
		return ""
	}
	if s, ok := t.lookup(lang); ok {
		return s
	}
	if s, ok := t.lookup(Fallback); ok {
		return s
	}
	return t[0].Text
}

func (t Text) lookup(lang string) (string, bool) {
	if lang == "" {
		return "", false
	}
	for _, tr := range t {
		if tr.Lang == lang {
			return tr.Text, true
		}
	}
	return "", false
}

// Empty reports whether the Text carries no translation at all — the test a
// field that used to be compared against "" now makes. A Text that lists a
// language but words it as "" is not empty; nothing authors one.
func (t Text) Empty() bool { return len(t) == 0 }

// Map rewrites every translation, keeping the languages. A derived label —
// a scene title cut down to the part that fits on a map pin — is the same
// cut in every language.
func (t Text) Map(fn func(string) string) Text {
	out := make(Text, len(t))
	for i, tr := range t {
		out[i] = Translation{
			Lang:  tr.Lang,
			Text:  fn(tr.Text),
			Voice: tr.Voice,
		}
	}
	return out
}

// Format composes a sentence around other Texts, resolving every part in the
// same language as the frame it lands in:
//
//	var exitLook = inkwell.Text{
//		{Lang: "en", Text: "That way: %s."},
//		{Lang: "hu", Text: "Arra van: %s."},
//	}
//
//	inkwell.Format(exitLook, e.Label)
//
// An argument that is a Text resolves per language; anything else is printed
// as it stands, which is what a number or a bare identifier wants. The
// result is a Text covering exactly the languages the frame covers.
func Format(format Text, args ...any) Text {
	out := make(Text, len(format))
	for i, tr := range format {
		in := make([]any, len(args))
		for j, a := range args {
			if t, ok := a.(Text); ok {
				in[j] = t.In(tr.Lang)
			} else {
				in[j] = a
			}
		}
		out[i] = Translation{
			Lang: tr.Lang,
			Text: fmt.Sprintf(tr.Text, in...),
		}
	}
	return out
}
