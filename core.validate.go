package inkwell

import "fmt"

func (g *Game) Validate() error {
	if !g.SceneManager.Has(g.startID) {
		return fmt.Errorf("%w: %q", ErrNoStartScene, g.startID)
	}
	for _, name := range g.SceneManager.Names() {
		s := g.SceneManager.MustGet(name)
		if bg := s.BackgroundAsset(); !g.AssetManager.Has(bg) {
			return fmt.Errorf("%w: scene %q background %q", ErrUnknownAsset, name, bg)
		}
		if s.Music != "" && !g.AssetManager.Has(s.Music) {
			return fmt.Errorf("%w: scene %q music %q", ErrUnknownAsset, name, s.Music)
		}
		for _, a := range g.sceneActors(s) {
			if !g.CharacterManager.Has(a.CharacterName) {
				return fmt.Errorf("%w: scene %q actor %q", ErrUnknownCharacter, name, a.CharacterName)
			}
		}
		for _, e := range s.Exits {
			if e.To != "" && !g.SceneManager.Has(e.To) {
				return fmt.Errorf("%w: scene %q exit to %q", ErrUnknownScene, name, e.To)
			}
		}
	}
	for _, name := range []string{g.onStart, g.onFinale} {
		if name != "" && !g.ScriptManager.Has(name) {
			return fmt.Errorf("%w: %q", ErrUnknownScript, name)
		}
	}
	if g.Player != "" && !g.CharacterManager.Has(g.Player) {
		return fmt.Errorf("%w: player %q", ErrUnknownCharacter, g.Player)
	}
	if g.activeTheme == "" || !g.ThemeManager.Has(g.activeTheme) {
		return fmt.Errorf("inkwell: no active theme (got %q)", g.activeTheme)
	}
	return nil
}
