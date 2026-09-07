package widget

import "git.teletypegames.org/engines/inkwell/inkbase"

// WidgetVisible reports whether a widget takes part in this frame at all.
func WidgetVisible(ctx *inkbase.Ctx, w Widget) bool {
	g, ok := w.(Gated)
	if !ok {
		return true
	}
	c := g.VisibleWhen()
	return c == nil || c.Eval(ctx)
}

// WidgetAwake reports whether a widget ticks during a running script.
func WidgetAwake(w Widget) bool {
	a, ok := w.(Awake)
	return ok && a.TicksDuringScript()
}

// LayerOf reports the layer a widget paints on.
func LayerOf(w Widget) Layer {
	if l, ok := w.(Layered); ok {
		return l.Layer()
	}
	return LayerHUD
}
