package domain

import p "pncdsl/pncdsl"

func defineIntroScript(g *p.Game) {
	g.ScriptManager.Register(p.Script{
		Name: "intro",
		Actions: p.Seq(
			p.Wait(0.4),
			p.Say("player", "Brr, hideg van."),
			p.Say("player", "Egy kávé kéne, mielőtt szétszakad a fejem."),
			p.Walk("player", p.Point{X: 200, Y: 145}),
			p.SetVar("player.mood", "Eltökélt"),
		),
	})
}
