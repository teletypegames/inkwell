package domain

import (
	"image/color"

	p "pncdsl/pncdsl"
)

func defineCat(g *p.Game) {
	g.CharacterManager.Register(p.Character{
		Name:        "cat",
		Sprite:      "spr/cat",
		Speed:       40,
		SpeechColor: color.RGBA{200, 200, 255, 255},
		W:           34,
		H:           20,
	})
}
