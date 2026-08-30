package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// EndCard is a fullscreen overlay shown when the ShowEnd action fires.
// While active it consumes every input so nothing underneath reacts.
type EndCard struct {
	Name string
}

func (e *EndCard) GetName() string { return e.Name }
func (e *EndCard) Layer() Layer    { return LayerCurtain }

func (e *EndCard) Tick(ctx *UICtx) {
	if ctx.Game.endCard == "" {
		return
	}
	// swallow every click while the end card is up
	if ctx.Game.Input.LeftClicked() {
		ctx.Game.Input.ConsumeLeft()
	}
	if ctx.Game.Input.RightClicked() {
		ctx.Game.Input.ConsumeRight()
	}
}

func (e *EndCard) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	if g.endCard == "" {
		return
	}
	th := g.Theme()
	vector.DrawFilledRect(dst, 0, 0, float32(g.Width), float32(g.Height), th.EndCardBG, false)
	w := textWidth(g.endCard)
	g.DrawText(dst, g.endCard, (g.Width-w)/2, g.Height/2-8, th.EndCardText)
}
