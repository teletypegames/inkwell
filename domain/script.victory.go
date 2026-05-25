package domain

import p "pncdsl/pncdsl"

func defineVictoryScript(g *p.Game) {
	g.ScriptManager.Register(p.Script{
		Name: "victory",
		Actions: p.Seq(
			p.Say("player", "Készen is van!"),
			p.PlaySound("snd/coffee_done"),
			p.Wait(0.6),
			p.Say("cat", "Mióóóóu. *boldog macska*"),
			p.Say("player", "Reggeli kávé: megmentve."),
			p.ShowEnd("Vége — kösz, hogy játszottál!"),
		),
	})
}
