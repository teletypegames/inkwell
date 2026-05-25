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

	// HotspotOnly (trigger mode only) makes the wheel open only when the
	// trigger lands on a registered hotspot in the current scene. Right-
	// clicks on empty space fall through to handleSceneInput (so the
	// SCUMM "reset verb to look" behaviour still works there).
	HotspotOnly bool

	// Labels overrides the rendered verb label per verb Name. Useful when
	// the regular Verb.Label is too long to fit inside the wheel — e.g.
	// {"use": "Tedd", "take": "Vedd"}. Falls back to Verb.Label.
	Labels map[string]string

	// runtime
	visible        bool
	center         Point
	pendingHotspot *Hotspot // remembered hotspot when HotspotOnly opens the coin
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
		// HotspotOnly: only open when the trigger lands on a hotspot.
		// Empty-space clicks fall through (so the scene's right-click
		// "reset verb" behaviour still fires from handleSceneInput).
		var pending *Hotspot
		if r.HotspotOnly {
			pending = g.HotspotAt(mp)
			if pending == nil {
				return
			}
		}
		// Don't open over a widget that would also claim this click —
		// e.g. an InventoryBar. Each widget can declare its dead-zone
		// via the optional clickBlocker interface.
		if r.blockedAt(ctx, mp) {
			return
		}
		r.consumeTrigger(g)
		r.center = mp
		r.visible = true
		r.pendingHotspot = pending
		return
	}
	// modal: swallow secondary clicks (close without picking)
	if r.triggerPressed(g) {
		r.consumeTrigger(g)
		r.visible = false
		r.pendingHotspot = nil
		return
	}
	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	idx := r.hitTest(g, g.Input.Point())
	if idx >= 0 {
		names := g.VerbManager.Names()
		verb := names[idx]
		g.SetSelectedVerb(verb)
		// Verb-coin behaviour: when the coin was opened over a hotspot
		// (HotspotOnly mode), picking a slice immediately fires that
		// verb's action on that hotspot — one round-trip per command,
		// instead of forcing a follow-up left-click on the same target.
		if r.pendingHotspot != nil {
			g.invokeHotspotVerb(r.pendingHotspot, verb)
		}
	}
	r.visible = false
	r.pendingHotspot = nil
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
		// DialogBox advertises its bounds via clickBlocker even when no
		// dialog is open, so the verb-coin would refuse to pop up over
		// scenery that happens to sit under the dialog rect. Skip the
		// dialog when it's actually inactive.
		if _, isDialog := w.(*DialogBox); isDialog && !ctx.Game.dialogueActive() {
			continue
		}
		// CharacterPanel auto-hides when its character isn't in the
		// current scene; treat the hidden panel as transparent so the
		// coin can open over the empty space it would occupy.
		if cp, isCP := w.(*CharacterPanel); isCP {
			if cp.Character == "" || !ctx.Game.CharacterInScene(cp.Character) {
				continue
			}
		}
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

	cx, cy := float32(r.center.X), float32(r.center.Y)

	// translucent backdrop disk that defines the wheel area
	vector.DrawFilledCircle(dst, cx, cy, float32(radius), color.RGBA{0, 0, 0, 170}, true)

	hover := r.hitTest(g, g.Input.Point())
	selected := g.SelectedVerb()
	slice := 2 * math.Pi / float64(n)

	// Labels live INSIDE the disk, at 0.62×radius from the center, so a
	// short (3–5 char) word stays comfortably within the rim.
	labelDist := radius * 0.62
	for i, name := range verbs {
		a := slice*float64(i) + slice/2
		lx := r.center.X + math.Cos(a)*labelDist
		ly := r.center.Y + math.Sin(a)*labelDist

		label := r.labelFor(g, name)
		tw := textWidth(label)
		col := th.VerbButtonText
		isSelected := name == selected
		isHover := i == hover
		if isSelected || isHover {
			// Filled halo behind the label. Selected = stronger, hover =
			// lighter overlay on top so a hover over a non-selected verb
			// still shows feedback.
			haloR := float32(radius * 0.28)
			if isSelected {
				vector.DrawFilledCircle(dst, float32(lx), float32(ly+4), haloR, th.VerbButtonSelectedBG, true)
			}
			if isHover && !isSelected {
				vector.DrawFilledCircle(dst, float32(lx), float32(ly+4), haloR, th.VerbButtonBG, true)
			}
		}
		g.DrawText(dst, label, int(lx)-tw/2, int(ly)-7, col)
	}

	// center dot
	vector.DrawFilledCircle(dst, cx, cy, 2.5, th.CursorColor, true)
}

func (r *RadialVerbs) labelFor(g *Game, verbName string) string {
	if r.Labels != nil {
		if alt, ok := r.Labels[verbName]; ok {
			return alt
		}
	}
	if v, ok := g.VerbManager.Get(verbName); ok {
		return v.Label
	}
	return verbName
}
