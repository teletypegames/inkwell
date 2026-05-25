package pncdsl

// RegisterDefaultUI installs the SCUMM-style HUD widgets into g.UIManager
// in the conventional Z-order (HotspotDebug at the back, Cursor on top).
// Domains that want a different layout call this and then either tweak
// the registered widgets in place or replace them entirely.
//
// Layout assumes the default 320×200 internal resolution.
func RegisterDefaultUI(g *Game) {
	g.UIManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	g.UIManager.Register(&VerbBar{
		Name:   "verbs",
		Origin: Point{X: 4, Y: 152},
		Cols:   2,
		Button: Size{W: 60, H: 14},
	})
	g.UIManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 132, Y: 152},
		Slots:    8,
		Cols:     8,
		SlotSize: 22,
		Gap:      2,
	})
	g.UIManager.Register(&StatusLine{Name: "status", Y: 142, Align: AlignCenter})
	g.UIManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 14})
	g.UIManager.Register(&DialogBox{Name: "dialog", Bounds: Rect(0, 140, float64(g.Width), 60), LineHeight: 14, Padding: 6})
	g.UIManager.Register(&EndCard{Name: "endcard"})
	g.UIManager.Register(&Cursor{Name: "cursor"})
}

// RegisterRadialVerbUI installs an alternative HUD that swaps the
// permanent verb-bar for a verb-coin (right-click radial menu). Inventory,
// dialog, speech and cursor stay the same.
func RegisterRadialVerbUI(g *Game) {
	g.UIManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	g.UIManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 4, Y: 178},
		Slots:    14,
		Cols:     14,
		SlotSize: 22,
		Gap:      0,
	})
	g.UIManager.Register(&StatusLine{Name: "status", Y: 168, Align: AlignCenter})
	g.UIManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 14})
	g.UIManager.Register(&DialogBox{Name: "dialog", Bounds: Rect(0, 140, float64(g.Width), 60), LineHeight: 14, Padding: 6})
	g.UIManager.Register(&RadialVerbs{Name: "verbs", Trigger: MouseButtonRight, Radius: 40})
	g.UIManager.Register(&EndCard{Name: "endcard"})
	g.UIManager.Register(&Cursor{Name: "cursor"})
}
