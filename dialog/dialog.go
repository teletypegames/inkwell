package dialog

import "git.teletypegames.org/engines/inkwell/inkbase"

// Dialog is a conversation: a set of named nodes, each a few lines spoken
// in turn followed by the choices the player picks from. A choice's Actions
// move it on — to another node, out of the dialog, or anywhere else an
// [inkbase.Action] can go.
type Dialog struct {
	// Name is the key the dialog is registered and started by.
	Name string
	// Start is the node the dialog opens on. Empty opens on the first of
	// Nodes.
	Start string
	// Nodes are the dialog's nodes, looked up by their Name.
	Nodes []DialogNode
}

// GetName returns the dialog's Name, the key it is registered under.
func (d Dialog) GetName() string { return d.Name }

// TypeLabel returns "dialog", the kind that tells a dialog apart from
// other entities of the same Name in saved state and in panic messages.
func (d Dialog) TypeLabel() string { return "dialog" }

// Node returns the node called name and true, or the zero DialogNode and
// false when the dialog has none by that name. When two nodes share a name
// the first wins.
func (d Dialog) Node(name string) (DialogNode, bool) {
	for _, n := range d.Nodes {
		if n.Name == name {
			return n, true
		}
	}
	return DialogNode{}, false
}

// DialogNode is one step of a [Dialog]: Lines are shown one per click, and
// once they run out the Choices the player may currently see are offered.
type DialogNode struct {
	// Name is the node's key within its dialog, used by Start and by a
	// choice that goes to it.
	Name    string
	Lines   []DialogLine
	Choices []DialogChoice
}

// DialogLine is one line of a [DialogNode]: who says it and what, in every
// language the game speaks.
type DialogLine struct {
	Speaker inkbase.Talker
	Text    inkbase.Text
}

// DialogChoice is an answer the player can pick at the end of a
// [DialogNode]. Picking it runs its Actions in sequence.
type DialogChoice struct {
	// Text is what the choice reads.
	Text inkbase.Text
	// Show, when set, hides the choice while it evaluates false. Nil shows
	// it always.
	Show inkbase.Condition
	// Once hides the choice for good after the player has picked it. It is
	// remembered by node name and position in Choices, so it survives a
	// save and a change of language.
	Once bool
	// Actions run, one after the other, when the choice is picked.
	Actions []inkbase.Action
}
