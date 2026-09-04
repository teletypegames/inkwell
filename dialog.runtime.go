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

func (g *Game) startDialog(d *Dialog) {
	if d == nil {
		logf("startDialog: nil dialog")
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
	g.dialog = &runtimeDialog{Dialog: d, Node: &node, LineIdx: 0}
}

func (g *Game) endDialog() {
	g.dialog = nil
}

func (g *Game) gotoDialogNode(name string) {
	if g.dialog == nil {
		return
	}
	node, ok := g.dialog.Dialog.Node(name)
	if !ok {
		logf("gotoNode: unknown %q in %q", name, g.dialog.Dialog.Name)
		return
	}
	g.dialog.Node = &node
	g.dialog.LineIdx = 0
	g.dialog.ChoiceHits = nil
}

func (g *Game) dialogActive() bool { return g.dialog != nil }
