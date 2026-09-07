package inkwell

import "fmt"

// Validate cross-checks what the types cannot: that the game has somewhere
// to open, that a Stage which holds has somewhere to go, and that the
// active theme is registered. A reference to a scene, an item or a dialog
// is a pointer now, so an unknown name is a compile error rather than one
// of these.
func (g *Game) Validate() error {
	if g.StartScene == nil || !g.SceneManager.Has(g.StartScene.Name) {
		return fmt.Errorf("%w: %v", ErrNoStartScene, entityName(g.StartScene))
	}
	if g.StartStage == nil || !g.StageManager.Has(g.StartStage.Name) {
		return fmt.Errorf("%w: %v", ErrNoStartStage, entityName(g.StartStage))
	}
	for _, stage := range g.StageManager.All() {
		if stage.Hold > 0 && stage.Next == nil {
			return fmt.Errorf("%w: stage %q holds %gs but names no Next", ErrUnknownStage, stage.Name, stage.Hold)
		}
		if stage.Next != nil && !g.StageManager.Has(stage.Next.Name) {
			return fmt.Errorf("%w: stage %q next %q is not registered", ErrUnknownStage, stage.Name, stage.Next.Name)
		}
	}
	for _, s := range g.SceneManager.All() {
		for _, a := range g.sceneActors(s) {
			if a.Character == nil || !g.CharacterManager.Has(a.Character.Name) {
				return fmt.Errorf("%w: scene %q actor %v", ErrUnknownCharacter, s.Name, entityName(a.Character))
			}
		}
		for _, e := range s.Exits {
			if e.To != nil && !g.SceneManager.Has(e.To.Name) {
				return fmt.Errorf("%w: scene %q exit to an unregistered scene", ErrUnknownScene, s.Name)
			}
		}
	}
	if g.Player != nil && !g.CharacterManager.Has(g.Player.Name) {
		return fmt.Errorf("%w: player %q", ErrUnknownCharacter, g.Player.Name)
	}
	for _, a := range g.AchievementManager.All() {
		if a.When == nil {
			return fmt.Errorf("inkwell: achievement %q names no When", a.Name)
		}
	}
	if g.activeTheme == "" || !g.ThemeManager.Has(g.activeTheme) {
		return fmt.Errorf("inkwell: no active theme (got %q)", g.activeTheme)
	}
	return nil
}

func entityName(named Named) string {
	if named == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", named.GetName())
}
