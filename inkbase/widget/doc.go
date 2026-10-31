// Package widget is the widget contract: the [Widget] interface, the [Layer]
// constants, and the visibility and ordering helpers a host drives its
// per-frame loop with.
//
// A Widget is [inkbase.Named] and has Tick and Draw, both handed a [UICtx]
// carrying the frame's DT and the Game as an opaque any, which the widget
// type-asserts to the host interface it needs. A widget is responsible for
// its own input, hit testing, state and drawing; there is no layout layer —
// widgets carry their own bounds or work them out.
//
// # Optional interfaces
//
//   - [Layered]: the widget declares its Layer. [LayerOf] reports it, and a
//     widget that does not implement Layered sits on LayerHUD.
//   - [Gated]: the widget can be switched off by game state through
//     VisibleWhen. [WidgetVisible] evaluates it; a widget whose condition is
//     false neither ticks nor draws nor blocks a click, so a HUD hides for a
//     cutscene by saying so on its widgets rather than by unregistering
//     them. A widget that does not implement Gated is always live.
//   - [Awake]: the widget ticks even while a script holds the frame.
//     [WidgetAwake] reports it.
//   - [ClickBlocker]: the widget owns part of the screen. A host asks every
//     visible widget's BlocksClickAt before it hands a click to whatever lies
//     behind the widgets (in inkwell the scene), so a click on the HUD does
//     not also walk the player or fire a hotspot.
//
// # Hosts
//
// The Game in a UICtx is whatever the host loop put there. A widget that
// needs more than DT type-asserts it to an interface of its own; the widget
// packages of inkbase each call theirs Host, and most embed [Screen] — the
// screen size and text drawing. A Game that does not implement a widget's
// Host leaves the widget inert rather than panicking. inkwell's *Game
// implements every one of them.
//
// # Z-order and input
//
// Draw runs from the bottom layer up ([OrderedWidgets]), so a widget on a
// higher layer is painted on top; inside one layer, registration order
// decides. Tick runs from the top layer down ([ReversedWidgets]), so the
// widget drawn on top gets the click first. A widget claims a click or a key
// by consuming it on the host's input (see the input package), and the
// widgets below see it as gone.
//
// A widget declares its layer rather than inheriting one from the order it
// was registered in, which is what lets a game register its widgets one file
// at a time, in whatever order the files fall, without the HUD coming out
// shuffled. A wrapper that forwards to an inner widget should return
// LayerOf(inner) from its own Layer, so that wrapping does not move it.
//
// A running script normally has the frame to itself: no widget ticks until
// it is done, which keeps a click from racing a cutscene. An Awake widget
// opts out and ticks on those frames too, before the script's runner, so a
// click or key it consumes is not also read as "skip this line". Two inkwell
// built-ins need it, for the same reason: a menu, which must stay closable
// while a cutscene plays, and the dialog box, because a conversation is the
// running script and its lines would never advance otherwise. Nothing else
// should.
//
// # Custom widgets
//
// Anything that implements Widget plugs in; the host cares about nothing
// beyond the name, and calls Tick and Draw:
//
//	type Clock struct{ Name string }
//
//	func (c *Clock) GetName() string             { return c.Name }
//	func (c *Clock) Layer() widget.Layer         { return widget.LayerHUD }
//	func (c *Clock) Tick(ctx *widget.UICtx)      {}
//	func (c *Clock) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
//		host := ctx.Game.(interface{ StateVar(string) any })
//		font.DrawText(dst, fmt.Sprint(host.StateVar("time")), 4, 4, color.White)
//	}
package widget
