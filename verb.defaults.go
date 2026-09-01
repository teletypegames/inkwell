package inkwell

// defaultVerbs returns the built-in SCUMM-style verb set. NewGame
// registers them so the demo works out of the box.
func defaultVerbs() []Verb {
	return []Verb{
		{
			Name: "look",
			Label: Text{
				{Lang: "en", Text: "Look at"},
				{Lang: "hu", Text: "Nézd"},
			},
		},
		{
			Name: "use",
			Label: Text{
				{Lang: "en", Text: "Use"},
				{Lang: "hu", Text: "Használd"},
			},
		},
		{
			Name: "talk",
			Label: Text{
				{Lang: "en", Text: "Talk to"},
				{Lang: "hu", Text: "Beszélj"},
			},
		},
		{
			Name: "take",
			Label: Text{
				{Lang: "en", Text: "Take"},
				{Lang: "hu", Text: "Vedd fel"},
			},
		},
	}
}
