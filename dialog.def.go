package inkwell

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
	Speaker Talker
	Text    Text
}

type DialogChoice struct {
	Text    Text
	Show    Condition // nil = always
	Once    bool
	Actions []Action
}
