package verb

import "git.teletypegames.org/engines/inkwell/inkbase"

// Defaults returns the built-in SCUMM-style verb set. New registers them
// so the demo works out of the box.
func Defaults() []*Verb {
	return []*Verb{
		{
			Name: "look",
			Label: inkbase.Text{
				{Lang: "en", Text: "Look at"},
				{Lang: "hu", Text: "Nézd"},
			},
		},
		{
			Name: "use",
			Label: inkbase.Text{
				{Lang: "en", Text: "Use"},
				{Lang: "hu", Text: "Használd"},
			},
		},
		{
			Name: "talk",
			Label: inkbase.Text{
				{Lang: "en", Text: "Talk to"},
				{Lang: "hu", Text: "Beszélj"},
			},
		},
		{
			Name: "take",
			Label: inkbase.Text{
				{Lang: "en", Text: "Take"},
				{Lang: "hu", Text: "Vedd fel"},
			},
		},
	}
}
