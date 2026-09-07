package panel

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

type Colors struct {
	BG     color.Color
	Border color.Color
}

type host interface {
	PanelColors() Colors
}

type Panel struct {
	Name        string
	When        inkbase.Condition
	Bounds      inkbase.Rectangle
	BG          color.Color
	Border      int
	BorderColor color.Color
	Blocks      bool
}

func (p *Panel) GetName() string                { return p.Name }
func (p *Panel) VisibleWhen() inkbase.Condition { return p.When }
func (p *Panel) Layer() widget.Layer            { return widget.LayerPanel }
func (p *Panel) Tick(ctx *widget.UICtx)         {}

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

func (p *Panel) BlocksClickAt(pt inkbase.Point) bool {
	return p.Blocks && p.Bounds.Contains(pt)
}
