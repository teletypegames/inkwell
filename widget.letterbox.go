package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Letterbox is the pair of bars that crop the picture top and bottom while
// a cutscene plays. It draws over the scene and under the HUD, so a game
// that also hides its HUD for a cutscene gets the whole convention out of
// two registrations with the same condition.
type Letterbox struct {
	Name string
	When Condition
	// Bar is each bar's height in pixels; 0 means 18.
	Bar float64
	// Color, if nil, is black.
	Color color.Color
}

// GetName returns the name the letterbox is registered under.
func (l *Letterbox) GetName() string { return l.Name }

// VisibleWhen returns the condition that brings the bars in, typically
// "a cutscene is playing"; nil means always.
func (l *Letterbox) VisibleWhen() Condition { return l.When }

// Layer puts the bars on [LayerScene], over the picture and under the HUD.
func (l *Letterbox) Layer() Layer { return LayerScene }

// Tick does nothing; the bars only display.
func (l *Letterbox) Tick(ctx *UICtx) {}

// Draw paints one bar of height Bar across the top of the screen and one
// across the bottom.
func (l *Letterbox) Draw(dst *ebiten.Image, ctx *UICtx) {
	bar := l.Bar
	if bar <= 0 {
		bar = 18
	}
	fill := l.Color
	if fill == nil {
		fill = color.Black
	}
	g := gameOfUI(ctx)
	w, h := float32(g.Width), float32(g.Height)
	vector.DrawFilledRect(dst, 0, 0, w, float32(bar), fill, false)
	vector.DrawFilledRect(dst, 0, h-float32(bar), w, float32(bar), fill, false)
}
