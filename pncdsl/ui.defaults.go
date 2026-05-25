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

// RegisterRichUI installs a "story-rich" HUD that matches the layout of
// the screenshot used as the reference: top bar with title/score/time,
// floating character panels for the player and a single NPC, a chat-log
// strip at the bottom-left, a tiny inventory at the bottom-right and a
// permanent verb-wheel on the right edge.
//
// The HUD assumes the default 320×200 internal resolution and that
// PlayerName / NPCName correspond to registered characters. Pass "" for
// NPCName to skip the NPC panel.
func RegisterRichUI(g *Game, playerName, npcName string) {
	g.UIManager.Register(&TopBar{
		Name:     "topbar",
		Height:   12,
		ScoreVar: "score", ScoreMax: 100,
		TimeVar: "time",
	})
	g.UIManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	if playerName != "" {
		g.UIManager.Register(&CharacterPanel{
			Name:      "panel_player",
			Bounds:    Rect(2, 14, 92, 36),
			Character: playerName,
			Title:     "PLAYER",
			Stats: []CharStat{
				{Label: "State", VarKey: playerName + ".state"},
				{Label: "Mood", VarKey: playerName + ".mood"},
			},
		})
	}
	if npcName != "" {
		g.UIManager.Register(&CharacterPanel{
			Name:      "panel_npc",
			Bounds:    Rect(180, 14, 92, 36),
			Character: npcName,
			Title:     "NPC",
			Stats: []CharStat{
				{Label: "State", VarKey: npcName + ".state"},
				{Label: "Mood", VarKey: npcName + ".mood"},
			},
		})
	}
	g.UIManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 4, Y: 188},
		Slots:    8,
		Cols:     8,
		SlotSize: 10,
		Gap:      1,
	})
	g.UIManager.Register(&ChatLog{
		Name:       "chat",
		Bounds:     Rect(2, 148, 280, 38),
		LineHeight: 9,
		Padding:    3,
		ShowBorder: true,
	})
	g.UIManager.Register(&RadialVerbs{
		Name:          "verbs",
		AlwaysVisible: true,
		Center:        Point{X: 290, Y: 90},
		Radius:        28,
	})
	g.UIManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 18})
	g.UIManager.Register(&DialogBox{Name: "dialog", Bounds: Rect(0, 90, float64(g.Width), 56), LineHeight: 12, Padding: 6})
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
