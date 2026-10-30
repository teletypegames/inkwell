package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Panel is a generic background widget — a colored rectangle. Useful as a
// backdrop for grouping other widgets or as a chat/notebook frame.
type Panel struct {
	Name   string
	When   Condition
	Bounds Rectangle
	// BG, if nil, falls back to the active Theme.PanelBG.
	BG color.Color
	// Border thickness in pixels (0 = no border).
	Border int
	// BorderColor, if nil, falls back to Theme.DialogBorder.
	BorderColor color.Color
	// Blocks makes the panel opaque to clicks, so a widget that pops up
	// where the pointer is — the verb coin — does not open over it. A
	// plate the HUD is mounted on wants this; a rule drawn on top of one
	// does not, since the plate underneath has already said so.
	Blocks bool
}

// GetName returns the name the panel is registered under.
func (p *Panel) GetName() string { return p.Name }

// VisibleWhen returns the condition that switches the panel on; nil means
// always.
func (p *Panel) VisibleWhen() Condition { return p.When }

// Layer puts the panel on [LayerPanel], under the HUD mounted on it.
func (p *Panel) Layer() Layer { return LayerPanel }

// Tick does nothing; the panel only displays.
func (p *Panel) Tick(ctx *UICtx) {}

// Draw fills Bounds with BG and, when Border is positive, outlines it in
// BorderColor.
func (p *Panel) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	bg := p.BG
	if bg == nil {
		bg = g.Theme().PanelBG
	}
	vector.DrawFilledRect(dst,
		float32(p.Bounds.X), float32(p.Bounds.Y),
		float32(p.Bounds.W), float32(p.Bounds.H), bg, false)
	if p.Border > 0 {
		bc := p.BorderColor
		if bc == nil {
			bc = g.Theme().DialogBorder
		}
		vector.StrokeRect(dst,
			float32(p.Bounds.X), float32(p.Bounds.Y),
			float32(p.Bounds.W), float32(p.Bounds.H),
			float32(p.Border), bc, false)
	}
}

// BlocksClickAt reports whether p falls inside Bounds of a panel that has
// Blocks set, so the verb coin does not open over it.
func (p *Panel) BlocksClickAt(pt Point) bool {
	return p.Blocks && p.Bounds.Contains(pt)
}
