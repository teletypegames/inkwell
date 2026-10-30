package widget

import (
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
