package inkwell

// changeScene takes the player to another scene through the fade: the work
// happens at the darkest frame, so OnLeave, the new scene's actors, its
// triggers, its music and its OnEnter all land behind black and the player
// never sees the swap. An unregistered name logs and changes nothing.
func (g *Game) changeScene(name string) {
	if !g.SceneManager.Has(name) {
		logf("changeScene: unknown %q", name)
		return
	}
	prev := g.currentScene
	g.transition.start(func() {
		if prev != "" && prev != name {
			g.previousScene = prev
		}
		if prev != "" {
			old := g.SceneManager.MustGet(prev)
			if old.OnLeave != nil {
				g.queueAction(old.OnLeave, "OnLeave")
			}
		}
		g.currentScene = name
		g.State.NoteVisit(name)
		g.placeActors(name)
		g.resetTriggers()
		s := g.SceneManager.MustGet(name)
		if s.Music != "" {
			g.Audio.PlayMusic(s.Music)
		}
		if s.OnEnter != nil {
			g.queueAction(s.OnEnter, "OnEnter")
		}
	})
}
