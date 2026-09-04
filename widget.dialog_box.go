package inkwell

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DialogBox renders the active conversation: one line at a time (advanced
// on click) followed by the player's choices. The widget is dormant
// unless ctx.Game.dialogActive() — and while active it consumes input
// so nothing underneath reacts.
type DialogBox struct {
	Name       string
	When       Condition
	Bounds     Rectangle
	LineHeight int
	Padding    int
}

func (d *DialogBox) GetName() string        { return d.Name }
func (d *DialogBox) VisibleWhen() Condition { return d.When }
func (d *DialogBox) Layer() Layer           { return LayerDialog }

// TicksDuringScript: the conversation *is* the running script. RunDialog
// holds the frame for as long as the dialog is up, so the box has to tick
// inside it — otherwise its lines never advance, no choice is ever clicked,
// and the dialog can never reach the EndDialog that would let go.
func (d *DialogBox) TicksDuringScript() bool { return true }

func (d *DialogBox) Tick(ctx *UICtx) {
	g := ctx.Game
	if !g.dialogActive() {
		return
	}
	// while a dialog is open, ALL clicks belong to it
	if g.Input.RightClicked() {
		g.Input.ConsumeRight()
	}
	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	mp := g.Input.Point()

	dlg := g.dialog
	if dlg.LineIdx >= 0 && dlg.LineIdx < len(dlg.Node.Lines) {
		dlg.LineIdx++
		if dlg.LineIdx >= len(dlg.Node.Lines) {
			dlg.LineIdx = -1
			dlg.ChoiceHits = nil
		}
		return
	}
	// pick the clicked choice
	for i, r := range dlg.ChoiceHits {
		if r.Contains(mp) {
			idx := d.resolveChoices(ctx)[i]
			choice := dlg.Node.Choices[idx]
			if choice.Once {
				g.State.NoteTalked(choiceKey(dlg.Node, idx))
			}
			if len(choice.Actions) > 0 {
				g.queueAction(Seq(choice.Actions...), "choice")
			}
			return
		}
	}
}

// resolveChoices is the choices the player can see right now, as positions
// in the node's own list.
func (d *DialogBox) resolveChoices(ctx *UICtx) []int {
	g := ctx.Game
	dlg := g.dialog
	out := make([]int, 0, len(dlg.Node.Choices))
	for i, c := range dlg.Node.Choices {
		gctx := g.makeCtx()
		if c.Show != nil && !c.Show.Eval(gctx) {
			continue
		}
		if c.Once && g.State.Talked(choiceKey(dlg.Node, i)) > 0 {
			continue
		}
		out = append(out, i)
	}
	return out
}

// choiceKey identifies a Once choice in saved state. It is the node plus the
// choice's position in it, never its wording: a choice the player has used
// up has to stay used up when the game is played in another language, and a
// translator has to be free to reword a line without reopening somebody's
// save.
func choiceKey(n *DialogNode, i int) string {
	return n.Name + ":" + strconv.Itoa(i)
}

// BlocksClickAt — when a dialog is open, every click is ours.
func (d *DialogBox) BlocksClickAt(p Point) bool {
	// without ctx we can't know if a dialog is active; conservatively block
	// when the cursor is in the dialog bounds (used only by RadialVerbs).
	b := d.Bounds
	if b.W == 0 {
		return false
	}
	return b.Contains(p)
}

func (d *DialogBox) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	if !g.dialogActive() {
		return
	}
	th := g.Theme()
	b := d.Bounds
	if b.W == 0 {
		b = Rect(0, 140, float64(g.Width), 60)
	}
	pad := d.Padding
	if pad <= 0 {
		pad = 6
	}
	lh := d.LineHeight
	if lh <= 0 {
		lh = 14
	}

	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), th.DialogBG, false)
	vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, th.DialogBorder, false)

	dlg := g.dialog
	if dlg.LineIdx >= 0 && dlg.LineIdx < len(dlg.Node.Lines) {
		ln := dlg.Node.Lines[dlg.LineIdx]
		g.DrawText(dst, g.SpeakerLabel(talkerName(ln.Speaker)).String()+":", int(b.X)+pad, int(b.Y)+2, th.DialogSpeaker)
		lines := WrapText(ln.Text.String(), int(b.W)-pad*2)
		for i, l := range lines {
			g.DrawText(dst, l, int(b.X)+pad, int(b.Y)+18+i*lh, th.DialogText)
		}
		return
	}

	choices := d.resolveChoices(ctx)
	dlg.ChoiceHits = dlg.ChoiceHits[:0]
	mp := g.Input.Point()
	for i, idx := range choices {
		c := dlg.Node.Choices[idx]
		r := Rect(b.X+float64(pad), b.Y+4+float64(i*lh), b.W-float64(pad)*2, float64(lh-1))
		dlg.ChoiceHits = append(dlg.ChoiceHits, r)
		bg := th.DialogChoiceBG
		if r.Contains(mp) {
			bg = th.DialogChoiceHover
		}
		vector.DrawFilledRect(dst, float32(r.X), float32(r.Y), float32(r.W), float32(r.H), bg, false)
		g.DrawText(dst, c.Text.String(), int(r.X)+2, int(r.Y)-1, th.DialogText)
	}
}
