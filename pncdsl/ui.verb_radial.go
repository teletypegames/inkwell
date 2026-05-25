package pncdsl

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// RadialVerbs is the verb-coin alternative to VerbBar: secondary-click on
// the scene opens a radial menu of every registered verb, click on a
// slice (or close to its center) picks that verb.
//
// To use it instead of a VerbBar, register RadialVerbs and skip VerbBar.
// Both can coexist if you really want.
type RadialVerbs struct {
	Name    string
	Trigger MouseButton
	Radius  float64

	// AlwaysVisible turns the widget into a permanent verb-wheel fixed at
	// Center. The trigger button is ignored in this mode; left-click picks
	// the verb slice under the cursor.
	AlwaysVisible bool
	Center        Point

	// runtime
	visible bool
	center  Point
}

func (r *RadialVerbs) GetName() string { return r.Name }

func (r *RadialVerbs) Tick(ctx *UICtx) {
	g := ctx.Game
	if r.Radius <= 0 {
		r.Radius = 40
	}

	if r.AlwaysVisible {
		r.visible = true
		r.center = r.Center
		if !g.Input.LeftClicked() {
			return
		}
		idx := r.hitTest(g, g.Input.Point())
		if idx < 0 {
			return
		}
		g.Input.ConsumeLeft()
		names := g.VerbManager.Names()
		g.SetSelectedVerb(names[idx])
		return
	}

	if !r.visible {
		if !r.triggerPressed(g) {
			return
		}
		mp := g.Input.Point()
		// Don't open over a widget that would also claim this click —
		// e.g. an InventoryBar. Each widget can declare its dead-zone
		// via the optional clickBlocker interface.
		if r.blockedAt(ctx, mp) {
			return
		}
		r.consumeTrigger(g)
		r.center = mp
		r.visible = true
		return
	}
	// modal: swallow secondary clicks (close without picking)
	if r.triggerPressed(g) {
		r.consumeTrigger(g)
		r.visible = false
		return
	}
	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	idx := r.hitTest(g, g.Input.Point())
	if idx >= 0 {
		names := g.VerbManager.Names()
		g.SetSelectedVerb(names[idx])
	}
	r.visible = false
}

func (r *RadialVerbs) triggerPressed(g *Game) bool {
	if r.Trigger == MouseButtonLeft {
		return g.Input.LeftClicked()
	}
	return g.Input.RightClicked()
}

func (r *RadialVerbs) consumeTrigger(g *Game) {
	if r.Trigger == MouseButtonLeft {
		g.Input.ConsumeLeft()
		return
	}
	g.Input.ConsumeRight()
}

// blockedAt walks the registered widgets and asks any clickBlocker
// whether its area covers the cursor — so the radial menu doesn't open
// over an InventoryBar that would also want this click.
func (r *RadialVerbs) blockedAt(ctx *UICtx, p Point) bool {
	for _, name := range ctx.Game.UIManager.Names() {
		if name == r.Name {
			continue
		}
		w := ctx.Game.UIManager.MustGet(name)
		if b, ok := w.(clickBlocker); ok && b.BlocksClickAt(p) {
			return true
		}
	}
	return false
}

// clickBlocker is the optional widget contract for "if the cursor is on
// me, claim the click — don't let modal pop-ups open over me." Implemented
// by VerbBar, InventoryBar, DialogBox by default.
type clickBlocker interface {
	BlocksClickAt(p Point) bool
}

func (r *RadialVerbs) hitTest(g *Game, p Point) int {
	verbs := g.VerbManager.Names()
	n := len(verbs)
	if n == 0 {
		return -1
	}
	dx := p.X - r.center.X
	dy := p.Y - r.center.Y
	dist := math.Hypot(dx, dy)
	if dist < r.Radius*0.25 || dist > r.Radius*1.4 {
		return -1
	}
	// angle 0 = right, going counter-clockwise; map to slice index
	a := math.Atan2(dy, dx)
	if a < 0 {
		a += 2 * math.Pi
	}
	slice := 2 * math.Pi / float64(n)
	idx := int(a / slice)
	if idx >= n {
		idx = n - 1
	}
	return idx
}

// BlocksClickAt — when the wheel is open (or always visible), clicks on
// it should not pass through to the scene underneath.
func (r *RadialVerbs) BlocksClickAt(p Point) bool {
	if !r.visible && !r.AlwaysVisible {
		return false
	}
	c := r.center
	if r.AlwaysVisible {
		c = r.Center
	}
	dx := p.X - c.X
	dy := p.Y - c.Y
	d := dx*dx + dy*dy
	radius := r.Radius * 1.4
	return d <= radius*radius
}

func (r *RadialVerbs) Draw(dst *ebiten.Image, ctx *UICtx) {
	if !r.visible {
		return
	}
	g := ctx.Game
	th := g.Theme()
	verbs := g.VerbManager.Names()
	n := len(verbs)
	if n == 0 {
		return
	}
	radius := r.Radius
	// translucent backdrop disk
	vector.DrawFilledCircle(dst, float32(r.center.X), float32(r.center.Y), float32(radius*1.15), color.RGBA{0, 0, 0, 140}, true)

	hover := r.hitTest(g, g.Input.Point())
	slice := 2 * math.Pi / float64(n)
	for i, name := range verbs {
		a := slice*float64(i) + slice/2
		x := r.center.X + math.Cos(a)*radius
		y := r.center.Y + math.Sin(a)*radius
		bg := th.VerbButtonBG
		if i == hover {
			bg = th.VerbButtonSelectedBG
		}
		bw, bh := 34, 14
		bx := int(x) - bw/2
		by := int(y) - bh/2
		vector.DrawFilledRect(dst, float32(bx), float32(by), float32(bw), float32(bh), bg, false)
		verb := g.VerbManager.MustGet(name)
		// center label as best we can
		tw := textWidth(verb.Label)
		g.DrawText(dst, verb.Label, bx+(bw-tw)/2, by-1, th.VerbButtonText)
	}
	// center dot
	vector.DrawFilledCircle(dst, float32(r.center.X), float32(r.center.Y), 2, th.CursorColor, true)
}
