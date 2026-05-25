package domain

import p "pncdsl/pncdsl"

func defineBeans(g *p.Game) {
	g.ItemManager.Register(p.Item{
		Name:        "beans",
		Sprite:      "spr/beans",
		Description: "őrölt kávébab",
	})
}
