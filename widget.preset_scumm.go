package inkwell

// RegisterDefaultUI installs the SCUMM-style HUD widgets into g.WidgetManager
// in the conventional Z-order (HotspotDebug at the back, Cursor on top).
// Domains that want a different layout call this and then either tweak
// the registered widgets in place or replace them entirely.
//
// Layout assumes the default 320×200 internal resolution.
func RegisterDefaultUI(g *Game) {
	g.WidgetManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	g.WidgetManager.Register(&VerbBar{
		Name:   "verbs",
		Origin: Point{X: 4, Y: 152},
		Cols:   2,
		Button: Size{W: 60, H: 14},
	})
	g.WidgetManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 132, Y: 152},
		Slots:    8,
		Cols:     8,
		SlotSize: 22,
		Gap:      2,
	})
	g.WidgetManager.Register(&StatusLine{Name: "status", Y: 142, Align: AlignCenter})
	g.WidgetManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 14})
	g.WidgetManager.Register(&DialogueBox{Name: "dialogue", Bounds: Rect(0, 140, float64(g.Width), 60), LineHeight: 14, Padding: 6})
	g.WidgetManager.Register(&EndCard{Name: "endcard"})
	g.WidgetManager.Register(&Cursor{Name: "cursor"})
}
