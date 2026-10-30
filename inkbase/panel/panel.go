package panel

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Colors is the fallback palette a [Panel] reads from its host when its own
// BG or BorderColor is nil.
type Colors struct {
	BG     color.Color
	Border color.Color
}

type host interface {
	PanelColors() Colors
}

// Panel is a colored rectangle: a backdrop for grouping other widgets or a
// frame for a log.
type Panel struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Bounds is the rectangle filled.
	Bounds inkbase.Rectangle
	// BG fills it; nil takes the host's [Colors].BG, and black without one.
	BG color.Color
	// Border is the outline's thickness in pixels; 0 draws none.
	Border int
	// BorderColor is the outline's color; nil takes the host's
	// [Colors].Border, and no outline without one.
	BorderColor color.Color
	// Blocks makes the panel swallow clicks inside Bounds.
	Blocks bool
}

// GetName returns the panel's Name, the key it is registered under.
func (p *Panel) GetName() string { return p.Name }

// VisibleWhen returns When, the condition that shows the panel; nil shows
// it always.
func (p *Panel) VisibleWhen() inkbase.Condition { return p.When }

// Layer puts the panel on [widget.LayerPanel], under the HUD.
func (p *Panel) Layer() widget.Layer { return widget.LayerPanel }

// Tick does nothing: the panel has no state of its own to advance.
func (p *Panel) Tick(ctx *widget.UICtx) {}

// Draw fills Bounds with the background and strokes the border when one
// is set.
func (p *Panel) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	bg := p.BG
	bc := p.BorderColor
	if h, ok := ctx.Game.(host); ok {
		colors := h.PanelColors()
		if bg == nil {
			bg = colors.BG
		}
		if bc == nil {
			bc = colors.Border
		}
	}
	if bg == nil {
		bg = color.Black
	}
	vector.DrawFilledRect(dst,
		float32(p.Bounds.X), float32(p.Bounds.Y),
		float32(p.Bounds.W), float32(p.Bounds.H), bg, false)
	if p.Border > 0 && bc != nil {
		vector.StrokeRect(dst,
			float32(p.Bounds.X), float32(p.Bounds.Y),
			float32(p.Bounds.W), float32(p.Bounds.H),
			float32(p.Border), bc, false)
	}
}

// BlocksClickAt reports whether the panel swallows a click at pt: only when
// Blocks is set and pt is inside Bounds.
func (p *Panel) BlocksClickAt(pt inkbase.Point) bool {
	return p.Blocks && p.Bounds.Contains(pt)
}
