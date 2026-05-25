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

	// SCUMM-style HUD. Swap for p.RegisterRadialVerbUI(g) to try the
	// verb-coin variant, or build a fully custom widget set by calling
	// g.UIManager.Register(...) directly.
	p.RegisterDefaultUI(g)

	g.StartAt("bedroom").OnStart(p.RunScript("intro"))
	return g
}
