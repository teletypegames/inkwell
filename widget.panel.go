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
}

func (p *Panel) GetName() string        { return p.Name }
func (p *Panel) VisibleWhen() Condition { return p.When }
func (p *Panel) Layer() Layer           { return LayerPanel }
func (p *Panel) Tick(ctx *UICtx)        {}

func (p *Panel) Draw(dst *ebiten.Image, ctx *UICtx) {
	bg := p.BG
	if bg == nil {
		bg = ctx.Game.Theme().PanelBG
	}
	vector.DrawFilledRect(dst,
		float32(p.Bounds.X), float32(p.Bounds.Y),
		float32(p.Bounds.W), float32(p.Bounds.H), bg, false)
	if p.Border > 0 {
		bc := p.BorderColor
		if bc == nil {
			bc = ctx.Game.Theme().DialogBorder
		}
		vector.StrokeRect(dst,
			float32(p.Bounds.X), float32(p.Bounds.Y),
			float32(p.Bounds.W), float32(p.Bounds.H),
			float32(p.Border), bc, false)
	}
}
