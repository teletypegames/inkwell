package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// VerbBar is the SCUMM-style permanent verb panel. It reads the
// VerbManager every frame so newly-registered verbs show up automatically.
type VerbBar struct {
	Name   string
	When   Condition
	Origin Point
	Cols   int
	Button Size
	Gap    Point

	// PanelBG paints a Theme.PanelBG backdrop behind the buttons.
	PanelBG bool
}

// GetName returns the name the bar is registered under.
func (v *VerbBar) GetName() string { return v.Name }

// VisibleWhen returns the condition that switches the bar on; nil means
// always.
func (v *VerbBar) VisibleWhen() Condition { return v.When }

// Layer puts the bar on [LayerHUD].
func (v *VerbBar) Layer() Layer { return LayerHUD }

func (v *VerbBar) buttons(ctx *UICtx) []verbButton {
	g := gameOfUI(ctx)
	cols := v.Cols
	if cols <= 0 {
		cols = 2
	}
	names := g.VerbManager.Names()
	out := make([]verbButton, 0, len(names))
	for i, n := range names {
		col := i % cols
		row := i / cols
		x := int(v.Origin.X) + col*(v.Button.W+int(v.Gap.X))
		y := int(v.Origin.Y) + row*(v.Button.H+int(v.Gap.Y))
		verb := g.VerbManager.MustGet(n)
		out = append(out, verbButton{
			Name:   n,
			Label:  verb.Label.String(),
			Bounds: Rect(float64(x), float64(y), float64(v.Button.W), float64(v.Button.H)),
		})
	}
	return out
}

// Tick selects the verb whose button is left-clicked and consumes the
// click.
func (v *VerbBar) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	if !g.Input.LeftClicked() {
		return
	}
	mp := g.Input.Point()
	for _, b := range v.buttons(ctx) {
		if b.Bounds.Contains(mp) {
			g.SetSelectedVerb(b.Name)
			g.Input.ConsumeLeft()
			return
		}
	}
}

// Draw paints the PanelBG backdrop when asked, then one button per
// registered verb in a grid of Cols columns (2 when zero), the selected
// verb's button highlighted.
func (v *VerbBar) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	th := g.Theme()
	if v.PanelBG {
		// rough bounding box of the buttons
		buttons := v.buttons(ctx)
		if len(buttons) > 0 {
			x, y := buttons[0].Bounds.X, buttons[0].Bounds.Y
			rx, ry := x, y
			for _, b := range buttons {
				if b.Bounds.X+b.Bounds.W > rx {
					rx = b.Bounds.X + b.Bounds.W
				}
				if b.Bounds.Y+b.Bounds.H > ry {
					ry = b.Bounds.Y + b.Bounds.H
				}
			}
			vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(rx-x+4), float32(ry-y+4), th.PanelBG, false)
		}
	}
	for _, b := range v.buttons(ctx) {
		bg := th.VerbButtonBG
		if b.Name == g.SelectedVerb() {
			bg = th.VerbButtonSelectedBG
		}
		vector.DrawFilledRect(dst, float32(b.Bounds.X), float32(b.Bounds.Y), float32(b.Bounds.W), float32(b.Bounds.H), bg, false)
		g.DrawText(dst, b.Label, int(b.Bounds.X)+2, int(b.Bounds.Y)-1, th.VerbButtonText)
	}
}

type verbButton struct {
	Name   string
	Label  string
	Bounds Rectangle
}

// BlocksClickAt reports whether p falls inside the bar, so the verb coin
// does not open over it. The bar is taken to be two rows high whatever the
// number of verbs.
func (v *VerbBar) BlocksClickAt(p Point) bool {
	cols := v.Cols
	if cols <= 0 {
		cols = 2
	}
	// rough bounding box: rows = ceil(verbCount/cols), but without a Game
	// reference we approximate by assuming the bar is small (at most 2 rows
	// in the default SCUMM layout).
	x := v.Origin.X
	y := v.Origin.Y
	w := float64(cols*(v.Button.W+int(v.Gap.X)) - int(v.Gap.X))
	h := float64(2*(v.Button.H+int(v.Gap.Y)) - int(v.Gap.Y))
	return Rect(x, y, w, h).Contains(p)
}
