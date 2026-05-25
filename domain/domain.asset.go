package domain

import p "pncdsl/pncdsl"

// registerAssets declares every asset the game references. Files don't
// need to exist on disk — the library will substitute deterministic
// colored placeholders for anything missing.
func registerAssets(g *p.Game) {
	imgs := []string{
		"bg/bedroom", "bg/kitchen",
		"spr/key", "spr/beans", "spr/mug",
		"spr/player", "spr/cat",
	}
	for _, name := range imgs {
		g.AssetManager.Register(p.Asset{
			Name: name,
			Path: "assets/" + name + ".png",
			Kind: p.AssetImage,
		})
	}
	audios := []string{
		"mus/wakeup", "mus/calm", "snd/coffee_done",
	}
	for _, name := range audios {
		g.AssetManager.Register(p.Asset{
			Name: name,
			Path: "assets/" + name + ".ogg",
			Kind: p.AssetAudio,
		})
	}
}
