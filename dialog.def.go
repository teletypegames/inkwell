package inkwell

type Dialogue struct {
	Name  string
	Start string
	Nodes []DialogueNode
}

func (d Dialogue) GetName() string   { return d.Name }
func (d Dialogue) TypeLabel() string { return "dialogue" }

func (d Dialogue) Node(name string) (DialogueNode, bool) {
	for _, n := range d.Nodes {
		if n.Name == name {
			return n, true
		}
	}
	return DialogueNode{}, false
}

type DialogueNode struct {
	Name    string
	Lines   []DialogueLine
	Choices []DialogueChoice
}

type DialogueLine struct {
	Speaker string
	Text    Text
}

type DialogueChoice struct {
	Text    Text
	Show    Condition // nil = always
	Once    bool
	Actions []Action
}
