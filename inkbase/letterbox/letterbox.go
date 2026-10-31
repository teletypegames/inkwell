package letterbox

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Host is what a [Letterbox] needs from the game it is drawn for, the
// UICtx's Game. On a Game that does not implement it the bars are not drawn.
type Host interface {
	// ScreenSize is the logical screen's width and height in pixels.
	ScreenSize() (width, height int)
}

// Letterbox is the pair of bars that crop the picture top and bottom while
// a cutscene plays. It draws over the scene and under the HUD, so a game
// that also hides its HUD for a cutscene gets the whole convention out of
// two registrations with the same condition.
type Letterbox struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Bar is each bar's height in pixels; 0 means 18.
	Bar float64
	// Color, if nil, is black.
	Color color.Color
}

// GetName returns the name the letterbox is registered under.
func (l *Letterbox) GetName() string { return l.Name }

// VisibleWhen returns the condition that brings the bars in, typically
// "a cutscene is playing"; nil means always.
func (l *Letterbox) VisibleWhen() inkbase.Condition { return l.When }

// Layer puts the bars on [widget.LayerScene], over the picture and under
// the HUD.
func (l *Letterbox) Layer() widget.Layer { return widget.LayerScene }

// Tick does nothing; the bars only display.
func (l *Letterbox) Tick(ctx *widget.UICtx) {}

// Draw paints one bar of height Bar across the top of the screen and one
// across the bottom.
func (l *Letterbox) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	bar := l.Bar
	if bar <= 0 {
		bar = 18
	}
	fill := l.Color
	if fill == nil {
		fill = color.Black
	}
	h, ok := ctx.Game.(Host)
	if !ok {
		return
	}
	w, sh := h.ScreenSize()
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(bar), fill, false)
	vector.DrawFilledRect(dst, 0, float32(sh)-float32(bar), float32(w), float32(bar), fill, false)
}
