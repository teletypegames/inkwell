package domain

import p "pncdsl/pncdsl"

// Build assembles the full Reggeli Kávé game by registering every entity
// into its matching manager on a fresh *Game. The order matters only for
// readability — Validate() runs after Build to cross-check references.
func Build() *p.Game {
	g := p.NewGame("Reggeli Kávé", 320, 200)

	registerAssets(g)

	definePlayer(g)
	defineCat(g)

	defineKey(g)
	defineBeans(g)
	defineMug(g)

	defineCatMorning(g)

	defineIntroScript(g)
	defineVictoryScript(g)

	defineBedroom(g)
	defineKitchen(g)

	// Story-rich HUD: top bar + NPC character panel + chat log + permanent
	// verb wheel + small inventory. The player panel is skipped — the
	// player avatar is always visible in the scene, so the floating card
	// in the corner just gets in the way. Pass "player" as the second
	// arg to put it back.
	p.RegisterRichUI(g, "", "cat")
	g.MaxLogLines = 64

	// initial HUD vars
	g.State.SetVar("score", 0)
	g.State.SetVar("time", "Nap 1 - 07:23")
	g.State.SetVar("player.state", "Ébren")
	g.State.SetVar("player.mood", "Álmos")
	g.State.SetVar("cat.state", "Várakozik")
	g.State.SetVar("cat.mood", "Türelmetlen")

	g.StartAt("bedroom").OnStart(p.RunScript("intro"))
	return g
}
