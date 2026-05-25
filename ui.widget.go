package pncdsl

import "github.com/hajimehoshi/ebiten/v2"

// Widget is the minimal surface a HUD element must implement to participate
// in the per-frame loop. The library ships several built-in widgets
// (VerbBar, InventoryBar, DialogBox, etc.) and the domain can register
// arbitrary new ones — chat panels, minimaps, hotbars — as long as they
// satisfy this interface.
//
// Lifecycle:
//   - Tick runs once per frame in REVERSE registration order so the
//     top-most widget gets a chance to consume input first via
//     ctx.Game.Input.ConsumeLeft / ConsumeRight.
//   - Draw runs once per frame in REGISTRATION order, so widgets
//     registered later are painted on top.
type Widget interface {
	Named
	Tick(ctx *UICtx)
	Draw(dst *ebiten.Image, ctx *UICtx)
}

// UICtx is the per-tick context handed to widgets. Game is the root
// aggregate; DT is seconds since the previous frame.
type UICtx struct {
	Game *Game
	DT   float64
}

// MouseButton identifies which mouse button a widget binds to (e.g. the
// RadialVerbs widget uses MouseButtonRight by default).
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
)

// Size is a {Width, Height} pair in screen-space pixels.
type Size struct {
	W, H int
}

// Align controls horizontal text alignment for widgets that draw a single line.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)
