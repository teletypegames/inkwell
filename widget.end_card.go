package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// EndCard is a fullscreen overlay shown when the ShowEnd action fires.
// While active it consumes every input so nothing underneath reacts.
type EndCard struct {
	Name string
	When Condition
}

// GetName returns the name the card is registered under.
func (e *EndCard) GetName() string { return e.Name }

// VisibleWhen returns the condition that switches the card on; nil means
// always.
func (e *EndCard) VisibleWhen() Condition { return e.When }

// Layer puts the card on [LayerCurtain], over the whole game.
func (e *EndCard) Layer() Layer { return LayerCurtain }

// Tick consumes every left and right click while the end card is up.
func (e *EndCard) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	if g.endCard == "" {
		return
	}
	// swallow every click while the end card is up
	if g.Input.LeftClicked() {
		g.Input.ConsumeLeft()
	}
	if g.Input.RightClicked() {
		g.Input.ConsumeRight()
	}
}

// Draw fills the screen and prints the ShowEnd text centred on it, while the
// end card is up.
func (e *EndCard) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	if g.endCard == "" {
		return
	}
	th := g.Theme()
	vector.DrawFilledRect(dst, 0, 0, float32(g.Width), float32(g.Height), th.EndCardBG, false)
	w := TextWidth(g.endCard)
	g.DrawText(dst, g.endCard, (g.Width-w)/2, g.Height/2-8, th.EndCardText)
}
