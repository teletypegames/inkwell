package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SceneNav is a development aid: the left and right arrow keys step
// through the registered scenes in registration order, wrapping at both
// ends. It draws nothing. Gate it with When — or leave it unregistered —
// when a build should not let the player walk the catalogue.
type SceneNav struct {
	Name string
	When Condition
}

func (n *SceneNav) GetName() string        { return n.Name }
func (n *SceneNav) VisibleWhen() Condition { return n.When }
func (n *SceneNav) Layer() Layer           { return LayerScene }

func (n *SceneNav) Draw(dst *ebiten.Image, ctx *UICtx) {}

func (n *SceneNav) Tick(ctx *UICtx) {
	step := 0
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
		step = 1
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		step = -1
	default:
		return
	}
	if next := n.neighbour(ctx.Game, step); next != "" {
		ctx.Game.Do(GoTo(next))
	}
}

func (n *SceneNav) neighbour(g *Game, step int) string {
	deck := g.SceneManager.Names()
	if len(deck) < 2 {
		return ""
	}
	for i, name := range deck {
		if name == g.CurrentScene() {
			return deck[((i+step)%len(deck)+len(deck))%len(deck)]
		}
	}
	return deck[0]
}
