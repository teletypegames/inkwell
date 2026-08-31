package inkwell

// runtimeDialog holds the in-flight Dialogue state. Mutated by the
// startDialogue/gotoDialogueNode/endDialogue plumbing and read by the
// DialogBox widget.
type runtimeDialog struct {
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
	g.activeDialog = name
	g.dialog = &runtimeDialog{Dialogue: &d, Node: &node, LineIdx: 0}
}

func (g *Game) endDialogue() {
	g.activeDialog = ""
	g.dialog = nil
}

func (g *Game) gotoDialogueNode(name string) {
	if g.dialog == nil || g.activeDialog == "" {
		return
	}
	d := g.DialogueManager.MustGet(g.activeDialog)
	node, ok := d.Node(name)
	if !ok {
		logf("gotoNode: unknown %q in %q", name, g.activeDialog)
		return
	}
	g.dialog.Node = &node
	g.dialog.LineIdx = 0
	g.dialog.ChoiceHits = nil
}

func (g *Game) dialogueActive() bool { return g.dialog != nil }
