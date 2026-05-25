package pncdsl

type Verb struct {
	Name    string
	Label   string
	Default Action
}

func (v Verb) GetName() string   { return v.Name }
func (v Verb) TypeLabel() string { return "verb" }

// defaultVerbs returns the built-in SCUMM-style verb set. A game can replace
// or extend these by registering its own Verbs after NewGame.
func defaultVerbs() []Verb {
	return []Verb{
		{Name: "look", Label: "Nézd"},
		{Name: "use", Label: "Használd"},
		{Name: "talk", Label: "Beszélj"},
		{Name: "take", Label: "Vedd fel"},
	}
}
