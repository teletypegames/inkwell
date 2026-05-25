package domain

import (
	"image/color"

	p "pncdsl/pncdsl"
)

func definePlayer(g *p.Game) {
	g.CharacterManager.Register(p.Character{
		Name:        "player",
		Sprite:      "spr/player",
		Speed:       80,
		SpeechColor: color.RGBA{255, 230, 160, 255},
		W:           14,
		H:           28,
	})
}
