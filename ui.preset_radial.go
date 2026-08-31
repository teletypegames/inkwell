package inkwell

// RegisterRadialVerbUI installs an alternative HUD that swaps the
// permanent verb-bar for a verb-coin (right-click radial menu). Inventory,
// dialog, speech and cursor stay the same.
func RegisterRadialVerbUI(g *Game) {
	g.WidgetManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	g.WidgetManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 4, Y: 178},
		Slots:    14,
		Cols:     14,
		SlotSize: 22,
		Gap:      0,
	})
	g.WidgetManager.Register(&StatusLine{Name: "status", Y: 168, Align: AlignCenter})
	g.WidgetManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 14})
	g.WidgetManager.Register(&DialogBox{Name: "dialog", Bounds: Rect(0, 140, float64(g.Width), 60), LineHeight: 14, Padding: 6})
	g.WidgetManager.Register(&RadialVerbs{Name: "verbs", Trigger: MouseButtonRight, Radius: 40})
	g.WidgetManager.Register(&EndCard{Name: "endcard"})
	g.WidgetManager.Register(&Cursor{Name: "cursor"})
}
