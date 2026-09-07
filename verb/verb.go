package verb

import "git.teletypegames.org/engines/inkwell/inkbase"

// Verb is a registered "action word" (Look, Use, Talk, Take, …). The
// VerbBar / RadialVerbs widgets read the VerbManager to know which verbs
// to render; Hotspot.handler(name) looks up the bound action by Verb.Name.
type Verb struct {
	Name    string
	Label   inkbase.Text
	Default inkbase.Action
}

func (v Verb) GetName() string   { return v.Name }
func (v Verb) TypeLabel() string { return "verb" }
