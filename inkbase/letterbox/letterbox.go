package letterbox

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

type host interface {
	ScreenSize() (int, int)
}

type Letterbox struct {
	Name  string
	When  inkbase.Condition
	Bar   float64
	Color color.Color
}

func (l *Letterbox) GetName() string                { return l.Name }
func (l *Letterbox) VisibleWhen() inkbase.Condition { return l.When }
func (l *Letterbox) Layer() widget.Layer            { return widget.LayerScene }
func (l *Letterbox) Tick(ctx *widget.UICtx)         {}

func (l *Letterbox) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	bar := l.Bar
	if bar <= 0 {
		bar = 18
	}
	fill := l.Color
	if fill == nil {
		fill = color.Black
	}
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	w, sh := h.ScreenSize()
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(bar), fill, false)
	vector.DrawFilledRect(dst, 0, float32(sh)-float32(bar), float32(w), float32(bar), fill, false)
}
