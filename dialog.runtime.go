package inkwell

// runtimeDialog holds the in-flight Dialog state. Mutated by the
// startDialog/gotoDialogNode/endDialog plumbing and read by the
// DialogBox widget.
type runtimeDialog struct {
	Dialog     *Dialog
	Node       *DialogNode
	LineIdx    int
	ChoiceHits []Rectangle
}

func (g *Game) startDialog(name string) {
	d, ok := g.DialogManager.Get(name)
	if !ok {
		logf("startDialog: unknown %q", name)
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
	g.dialog = &runtimeDialog{Dialog: &d, Node: &node, LineIdx: 0}
}

func (g *Game) endDialog() {
	g.activeDialog = ""
	g.dialog = nil
}

func (g *Game) gotoDialogNode(name string) {
	if g.dialog == nil || g.activeDialog == "" {
		return
	}
	d := g.DialogManager.MustGet(g.activeDialog)
	node, ok := d.Node(name)
	if !ok {
		logf("gotoNode: unknown %q in %q", name, g.activeDialog)
		return
	}
	g.dialog.Node = &node
	g.dialog.LineIdx = 0
	g.dialog.ChoiceHits = nil
}

func (g *Game) dialogActive() bool { return g.dialog != nil }
