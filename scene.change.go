package inkwell

// sceneChange is one scene change in flight, as the action that asked for it
// sees it: whether the swap has landed, and the OnLeave and OnEnter it left
// for that action to run.
type sceneChange struct {
	landed bool
	hooks  []Action
}

// changeScene takes the player to another scene through the fade: the work
// happens at the darkest frame, so OnLeave, the new scene's actors, its
// triggers, its music and its OnEnter all land behind black and the player
// never sees the swap. With a nil change OnLeave and OnEnter are queued;
// otherwise they are handed to change for its action to run. An
// unregistered name logs, changes nothing and returns false.
func (g *Game) changeScene(name string, change *sceneChange) bool {
	if !g.SceneManager.Has(name) {
		logf("changeScene: unknown %q", name)
		return false
	}
	prev := g.currentScene
	g.transition.start(sceneTransition, func() {
		if prev != "" && prev != name {
			g.previousScene = prev
		}
		var hooks []Action
		if prev != "" {
			old := g.SceneManager.MustGet(prev)
			if old.OnLeave != nil {
				hooks = append(hooks, old.OnLeave)
			}
		}
		g.currentScene = name
		g.State.NoteVisit(name)
		g.placeActors(name)
		g.resetTriggers()
		s := g.SceneManager.MustGet(name)
		if s.Music != nil {
			g.Audio.PlayMusic(s.Music)
		}
		if s.OnEnter != nil {
			hooks = append(hooks, s.OnEnter)
		}
		if change == nil {
			for _, h := range hooks {
				g.queueAction(h, "scene hook")
			}
			return
		}
		change.landed = true
		change.hooks = hooks
	})
	return true
}

// sceneChangeRunner holds the script through a scene change: it waits for
// the swap, runs the old scene's OnLeave and the new scene's OnEnter in its
// place, and finishes when the new picture is fully uncovered.
type sceneChangeRunner struct {
	g      *Game
	change *sceneChange
	hooks  Runner
}

// Tick waits for the swap, then ticks the hooks one by one, then waits for
// the fade to finish. A change whose transition was replaced before it
// landed finishes at once. It never fails: a failing hook only ends the
// hooks.
func (r *sceneChangeRunner) Tick(ctx *Ctx) Status {
	if !r.change.landed {
		if !r.g.transition.active {
			return StatusDone
		}
		return StatusRunning
	}
	for r.hooks != nil || len(r.change.hooks) > 0 {
		if r.hooks == nil {
			r.hooks = r.change.hooks[0].Start()
			r.change.hooks = r.change.hooks[1:]
		}
		st := r.hooks.Tick(ctx)
		if st == StatusRunning {
			return StatusRunning
		}
		r.hooks = nil
		if st == StatusFailed {
			r.change.hooks = nil
		}
	}
	if r.g.transition.active {
		return StatusRunning
	}
	return StatusDone
}
