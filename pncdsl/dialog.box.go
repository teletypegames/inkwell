package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// dialogBox is the active-conversation UI: it shows the current node's
// lines (one at a time, advanced on click) followed by the visible choices.
type dialogBox struct {
	dialogue   *Dialogue
	node       *DialogueNode
	lineIdx    int   // -1 = lines done, showing choices
	choiceHits []Rectangle
}

func newDialogBox(d *Dialogue, n *DialogueNode) *dialogBox {
	return &dialogBox{dialogue: d, node: n, lineIdx: 0}
}

// resolveChoices returns the list of choices currently visible to the player.
func (db *dialogBox) resolveChoices(ctx *Ctx) []DialogueChoice {
	out := make([]DialogueChoice, 0, len(db.node.Choices))
	for _, c := range db.node.Choices {
		if c.Show != nil && !c.Show.Eval(ctx) {
			continue
		}
		if c.Once && ctx.Game.State.Talked(db.node.Name+":"+c.Text) > 0 {
			continue
		}
		out = append(out, c)
	}
	return out
}

// handleClick advances the dialog state. Returns the picked choice (if any).
func (db *dialogBox) handleClick(ctx *Ctx, p Point) (picked *DialogueChoice) {
	if db.lineIdx >= 0 && db.lineIdx < len(db.node.Lines) {
		db.lineIdx++
		if db.lineIdx >= len(db.node.Lines) {
			db.lineIdx = -1
			db.choiceHits = nil
		}
		return nil
	}
	for i, r := range db.choiceHits {
		if r.Contains(p) {
			ch := db.resolveChoices(ctx)[i]
			return &ch
		}
	}
	return nil
}

func (db *dialogBox) currentLine() (DialogueLine, bool) {
	if db.lineIdx >= 0 && db.lineIdx < len(db.node.Lines) {
		return db.node.Lines[db.lineIdx], true
	}
	return DialogueLine{}, false
}

func (db *dialogBox) draw(dst *ebiten.Image, g *Game) {
	// box across the bottom 60 px
	vector.DrawFilledRect(dst, 0, 140, 320, 60, color.RGBA{15, 15, 25, 240}, false)
	vector.StrokeRect(dst, 0, 140, 320, 60, 1, color.RGBA{80, 80, 110, 255}, false)

	if ln, ok := db.currentLine(); ok {
		speakerCol := color.RGBA{220, 220, 120, 255}
		drawText(dst, ln.Speaker+":", 6, 142, speakerCol)
		lines := wrapText(ln.Text, 300)
		for i, l := range lines {
			drawText(dst, l, 6, 158+i*14, color.White)
		}
		return
	}

	// choices
	ctx := g.makeCtx()
	choices := db.resolveChoices(ctx)
	db.choiceHits = db.choiceHits[:0]
	for i, c := range choices {
		y := 144 + i*14
		r := Rect(6, float64(y), 308, 13)
		db.choiceHits = append(db.choiceHits, r)
		mx, my := g.input.Pos()
		hover := r.Contains(Point{X: float64(mx), Y: float64(my)})
		bg := color.RGBA{30, 30, 45, 255}
		if hover {
			bg = color.RGBA{70, 70, 100, 255}
		}
		vector.DrawFilledRect(dst, float32(r.X), float32(r.Y), float32(r.W), float32(r.H), bg, false)
		drawText(dst, c.Text, int(r.X)+2, int(r.Y)-1, color.White)
		_ = i
	}
}
