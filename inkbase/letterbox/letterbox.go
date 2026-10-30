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

// Letterbox is the pair of bars that crop the picture top and bottom while
// a cutscene plays. It draws over the scene and under the HUD.
type Letterbox struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Bar is each bar's height in pixels; 0 means 18.
	Bar float64
	// Color fills the bars; nil means black.
	Color color.Color
}

// GetName returns the letterbox's Name, the key it is registered under.
func (l *Letterbox) GetName() string { return l.Name }

// VisibleWhen returns When, the condition that shows the letterbox; nil shows
// it always.
func (l *Letterbox) VisibleWhen() inkbase.Condition { return l.When }

// Layer puts the letterbox on [widget.LayerScene], with the scene.
func (l *Letterbox) Layer() widget.Layer { return widget.LayerScene }

// Tick does nothing: the letterbox has no state of its own to advance.
func (l *Letterbox) Tick(ctx *widget.UICtx) {}

// Draw fills a bar across the top and one across the bottom of the screen.
// It draws nothing on a host that does not report the screen size.
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
