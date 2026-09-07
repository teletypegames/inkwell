// Package widget provides the Widget interface, layer constants, and the
// visibility/ordering helpers the engine uses to drive its per-frame loop.
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

const (
	LayerScene   Layer = 0
	LayerPanel   Layer = 100
	LayerHUD     Layer = 200
	LayerSpeech  Layer = 300
	LayerDialog  Layer = 400
	LayerMenu    Layer = 500
	LayerCurtain Layer = 600
	LayerCursor  Layer = 700
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
