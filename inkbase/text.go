package inkbase

import "fmt"

// Lang is the language the game speaks. Empty means every Text falls
// through to Fallback, and then to the language it was authored in.
var Lang string

// Fallback is the language used when a Text says nothing in Lang.
var Fallback string

// Translation is one language's wording of one string.
type Translation struct {
	Lang string
	Text string
	// Voice is the recording of this line in this language, if there is one.
	Voice *Voice
}

// Text is a string in as many languages as somebody has written it in.
type Text []Translation

// Plain wraps a string that is not really a translation of anything — a
// registry key standing in for a label nobody wrote. It answers in every
// language, because there is nothing about it to translate.
func Plain(s string) Text {
	return Text{{
		Text: s,
	}}
}

// String resolves the Text against Lang: that language if it is there, then
// Fallback, then whichever translation was written first.
func (t Text) String() string { return t.In(Lang) }

// Resolved is the Translation String would read from: the one in Lang,
// else Fallback, else the first — with its Voice.
func (t Text) Resolved() Translation {
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
// String does.
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

// Empty reports whether the Text carries no translation at all.
func (t Text) Empty() bool { return len(t) == 0 }

// Map rewrites every translation, keeping the languages and voices.
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
// same language as the frame it lands in. An argument that is a Text resolves
// per language; anything else is printed as it stands.
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
