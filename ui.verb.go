package pncdsl

// Verb is a registered "action word" (Look, Use, Talk, Take, …). The
// VerbBar / RadialVerbs widgets read the VerbManager to know which verbs
// to render; Hotspot.handler(name) looks up the bound action by Verb.Name.
type Verb struct {
	Name    string
	Label   string
	Default Action
}

func (v Verb) GetName() string   { return v.Name }
func (v Verb) TypeLabel() string { return "verb" }

type VerbManager = Manager[Verb]

// defaultVerbs returns the built-in SCUMM-style verb set. NewGame
// registers them so the demo works out of the box.
func defaultVerbs() []Verb {
	return []Verb{
		{Name: "look", Label: "Nézd"},
		{Name: "use", Label: "Használd"},
		{Name: "talk", Label: "Beszélj"},
		{Name: "take", Label: "Vedd fel"},
	}
}
