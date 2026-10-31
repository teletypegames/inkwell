package widget

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

// Widget is the minimal surface a HUD element must implement.
type Widget interface {
	inkbase.Named
	Tick(ctx *UICtx)
	Draw(dst *ebiten.Image, ctx *UICtx)
}

// Layer is a widget's place in the paint order.
type Layer int

// The standard layers, from the back to the front. A higher Layer is drawn
// later, on top; the gaps leave room for layers in between.
const (
	LayerScene   Layer = 0   // the scene itself and what stands in it
	LayerPanel   Layer = 100 // panels docked around the scene
	LayerHUD     Layer = 200 // heads-up elements: inventory, hover labels
	LayerSpeech  Layer = 300 // speech bubbles
	LayerDialog  Layer = 400 // dialog choices
	LayerMenu    Layer = 500 // menus over the game
	LayerCurtain Layer = 600 // full-screen fades and flashes
	LayerCursor  Layer = 700 // the mouse cursor, above everything
)

// Layered is implemented by a widget that declares its own layer.
type Layered interface {
	Layer() Layer
}

// Gated is implemented by a widget that can be switched off by game state.
type Gated interface {
	VisibleWhen() inkbase.Condition
}

// Awake is implemented by a widget that ticks even during a running script.
type Awake interface {
	TicksDuringScript() bool
}

// ClickBlocker is the optional widget contract for "the screen under me is
// mine": a visible widget that reports p as blocked keeps a click there from
// reaching whatever the host puts behind its widgets — in inkwell the scene,
// with its hotspots and its walk — and keeps a widget that pops up where the
// pointer is (a verb coin) from opening over it. The widget's own Tick still
// gets the click first.
type ClickBlocker interface {
	BlocksClickAt(ctx *UICtx, p inkbase.Point) bool
}

// Screen is the drawing surface most widgets need from their host: the
// logical screen size and text drawing in the host's font. The widget
// packages embed it in their own Host interfaces.
type Screen interface {
	// ScreenSize is the logical screen's width and height in pixels.
	ScreenSize() (width, height int)
	// DrawText draws s in colour c with its top-left corner at x, y.
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
}
