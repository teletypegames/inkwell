package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// Widget is the minimal surface a HUD element must implement to participate
// in the per-frame loop. The library ships several built-in widgets
// (VerbBar, InventoryBar, DialogBox, etc.) and the domain can register
// arbitrary new ones — chat panels, minimaps, hotbars — as long as they
// satisfy this interface.
//
// Lifecycle:
//   - Tick runs once per frame from the TOP layer down so the top-most
//     widget gets a chance to consume input first via
//     ctx.Game.Input.ConsumeLeft / ConsumeRight.
//   - Draw runs once per frame from the BOTTOM layer up, so a widget on a
//     higher layer is painted on top.
//
// Registration order only decides ties inside one layer — see Layer.
type Widget interface {
	Named
	Tick(ctx *UICtx)
	Draw(dst *ebiten.Image, ctx *UICtx)
}

// Layer is a widget's place in the paint order. A widget says where it
// belongs rather than depending on the order it happened to be registered
// in, which frees a domain to register its widgets one file at a time.
type Layer int

const (
	// LayerScene is over the picture and under the HUD: hotspot outlines,
	// cutscene bars, and the invisible widgets that only tick.
	LayerScene Layer = 0
	// LayerPanel is the plate the HUD is drawn on.
	LayerPanel Layer = 100
	// LayerHUD is everything mounted on that plate: verbs, inventory,
	// status line, top bar.
	LayerHUD Layer = 200
	// LayerSpeech is speech bubbles, which float over the scene.
	LayerSpeech Layer = 300
	// LayerDialog is the dialogue box.
	LayerDialog Layer = 400
	// LayerMenu is what opens on top of the game: radial verbs, menus.
	LayerMenu Layer = 500
	// LayerCurtain is what covers the game: end cards, fades.
	LayerCurtain Layer = 600
	// LayerCursor is the pointer, and nothing else belongs above it.
	LayerCursor Layer = 700
)

// Layered is implemented by a widget that declares its own layer. Every
// built-in widget does. A widget that does not sits on LayerHUD.
type Layered interface {
	Layer() Layer
}

// LayerOf reports the layer a widget paints on. A wrapper widget that
// forwards to an inner one — a visibility gate, say — should return
// LayerOf(inner) so that wrapping does not move the widget.
func LayerOf(w Widget) Layer {
	if l, ok := w.(Layered); ok {
		return l.Layer()
	}
	return LayerHUD
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
