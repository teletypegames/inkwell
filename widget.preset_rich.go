package inkwell

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
	g.WidgetManager.Register(&TopBar{
		Name:     "topbar",
		Height:   12,
		ScoreVar: "score", ScoreMax: 100,
		TimeVar: "time",
	})
	g.WidgetManager.Register(&HotspotDebug{Name: "hotspot_debug"})
	if playerName != "" {
		g.WidgetManager.Register(&CharacterPanel{
			Name:      "panel_player",
			Bounds:    Rect(2, 14, 92, 36),
			Character: playerName,
			Title: Text{
				{Lang: "en", Text: "PLAYER"},
				{Lang: "hu", Text: "JÁTÉKOS"},
			},
			Stats: []CharStat{
				{
					Label: Text{
						{Lang: "en", Text: "State"},
						{Lang: "hu", Text: "Állapot"},
					},
					VarKey: playerName + ".state",
				},
				{
					Label: Text{
						{Lang: "en", Text: "Mood"},
						{Lang: "hu", Text: "Kedv"},
					},
					VarKey: playerName + ".mood",
				},
			},
		})
	}
	if npcName != "" {
		g.WidgetManager.Register(&CharacterPanel{
			Name:      "panel_npc",
			Bounds:    Rect(180, 14, 92, 36),
			Character: npcName,
			Title: Text{
				{Lang: "en", Text: "NPC"},
				{Lang: "hu", Text: "NPC"},
			},
			Stats: []CharStat{
				{
					Label: Text{
						{Lang: "en", Text: "State"},
						{Lang: "hu", Text: "Állapot"},
					},
					VarKey: npcName + ".state",
				},
				{
					Label: Text{
						{Lang: "en", Text: "Mood"},
						{Lang: "hu", Text: "Kedv"},
					},
					VarKey: npcName + ".mood",
				},
			},
		})
	}
	g.WidgetManager.Register(&InventoryBar{
		Name:     "inventory",
		Origin:   Point{X: 4, Y: 188},
		Slots:    8,
		Cols:     8,
		SlotSize: 10,
		Gap:      1,
	})
	g.WidgetManager.Register(&ChatLog{
		Name:       "chat",
		Bounds:     Rect(2, 148, 280, 38),
		LineHeight: 9,
		Padding:    3,
		ShowBorder: true,
	})
	g.WidgetManager.Register(&RadialVerbs{
		Name:          "verbs",
		AlwaysVisible: true,
		Center:        Point{X: 282, Y: 90},
		Radius:        36,
		Labels: map[string]Text{
			"look": {
				{Lang: "en", Text: "Look"},
				{Lang: "hu", Text: "Néz"},
			},
			"use": {
				{Lang: "en", Text: "Use"},
				{Lang: "hu", Text: "Tedd"},
			},
			"talk": {
				{Lang: "en", Text: "Talk"},
				{Lang: "hu", Text: "Szól"},
			},
			"take": {
				{Lang: "en", Text: "Take"},
				{Lang: "hu", Text: "Vedd"},
			},
		},
	})
	g.WidgetManager.Register(&SpeechBubble{Name: "speech", MaxWidth: 200, Padding: 3, OffsetY: 4, FallbackY: 18})
	g.WidgetManager.Register(&DialogBox{Name: "dialog", Bounds: Rect(0, 90, float64(g.Width), 56), LineHeight: 12, Padding: 6})
	g.WidgetManager.Register(&EndCard{Name: "endcard"})
	g.WidgetManager.Register(&Cursor{Name: "cursor"})
}
