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

// Host is the optional part of the game a [Panel] reads its fallback
// colours from, the UICtx's Game. A Panel on a Game that does not implement
// it still draws, with its own colours only.
type Host interface {
	// PanelColors is the fallback palette.
	PanelColors() Colors
}

// Panel is a generic background widget — a colored rectangle. Useful as a
// backdrop for grouping other widgets or as a chat/notebook frame.
type Panel struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Bounds is the rectangle filled.
	Bounds inkbase.Rectangle
	// BG fills it; nil takes the host's [Colors].BG (in inkwell the active
	// theme's PanelBG), and black without one.
	BG color.Color
	// Border is the outline's thickness in pixels; 0 draws none.
	Border int
	// BorderColor is the outline's colour; nil takes the host's
	// [Colors].Border (in inkwell the theme's DialogBorder), and no outline
	// without one.
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
func (p *Panel) VisibleWhen() inkbase.Condition { return p.When }

// Layer puts the panel on [widget.LayerPanel], under the HUD mounted on
// it.
func (p *Panel) Layer() widget.Layer { return widget.LayerPanel }

// Tick does nothing; the panel only displays.
func (p *Panel) Tick(ctx *widget.UICtx) {}

// Draw fills Bounds with BG and, when Border is positive, outlines it in
// BorderColor.
func (p *Panel) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	bg := p.BG
	bc := p.BorderColor
	if h, ok := ctx.Game.(Host); ok {
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

// BlocksClickAt reports whether pt falls inside Bounds of a panel that has
// Blocks set; such a panel keeps clicks off the scene behind it.
func (p *Panel) BlocksClickAt(ctx *widget.UICtx, pt inkbase.Point) bool {
	return p.Blocks && p.Bounds.Contains(pt)
}
