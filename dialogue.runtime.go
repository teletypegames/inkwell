package inkwell

// runtimeDialogue holds the in-flight Dialogue state. Mutated by the
// startDialogue/gotoDialogueNode/endDialogue plumbing and read by the
// DialogueBox widget.
type runtimeDialogue struct {
	Dialogue   *Dialogue
	Node       *DialogueNode
	LineIdx    int
	ChoiceHits []Rectangle
}

func (g *Game) startDialogue(name string) {
	d, ok := g.DialogueManager.Get(name)
	if !ok {
		logf("startDialogue: unknown %q", name)
		return
	}
	startNode := d.Start
	if startNode == "" && len(d.Nodes) > 0 {
		startNode = d.Nodes[0].Name
	}
	node, ok := d.Node(startNode)
	if !ok {
		return
	}
	g.activeDialogue = name
	g.dialogue = &runtimeDialogue{Dialogue: &d, Node: &node, LineIdx: 0}
}

func (g *Game) endDialogue() {
	g.activeDialogue = ""
	g.dialogue = nil
}

func (g *Game) gotoDialogueNode(name string) {
	if g.dialogue == nil || g.activeDialogue == "" {
		return
	}
	d := g.DialogueManager.MustGet(g.activeDialogue)
	node, ok := d.Node(name)
	if !ok {
		logf("gotoNode: unknown %q in %q", name, g.activeDialogue)
		return
	}
	g.dialogue.Node = &node
	g.dialogue.LineIdx = 0
	g.dialogue.ChoiceHits = nil
}

func (g *Game) dialogueActive() bool { return g.dialogue != nil }
