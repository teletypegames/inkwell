package dialog

import "git.teletypegames.org/engines/inkwell/inkbase"

type Dialog struct {
	Name  string
	Start string
	Nodes []DialogNode
}

func (d Dialog) GetName() string   { return d.Name }
func (d Dialog) TypeLabel() string { return "dialog" }

func (d Dialog) Node(name string) (DialogNode, bool) {
	for _, n := range d.Nodes {
		if n.Name == name {
			return n, true
		}
	}
	return DialogNode{}, false
}

type DialogNode struct {
	Name    string
	Lines   []DialogLine
	Choices []DialogChoice
}

type DialogLine struct {
	Speaker inkbase.Talker
	Text    inkbase.Text
}

type DialogChoice struct {
	Text    inkbase.Text
	Show    inkbase.Condition // nil = always
	Once    bool
	Actions []inkbase.Action
}
