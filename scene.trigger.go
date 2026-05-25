package pncdsl

// Trigger fires Do when its When condition first becomes true after the
// trigger arms (on scene enter). Rising-edge semantics: the Do action is
// queued only on a false→true transition, so a condition that becomes
// true and stays true fires exactly once per arming. Once=true makes the
// trigger fire at most once per arming regardless of further toggles.
type Trigger struct {
	Name string
	When Condition
	Do   Action
	Once bool
}

// triggerState is the per-trigger bookkeeping for the active scene.
// Cleared on scene enter and on Load.
type triggerState struct {
	lastTrue bool
	fired    bool
}

// resetTriggers re-arms every trigger on the current scene — called on
// scene change and after Load. Triggers are not pre-armed: the first
// frame samples their condition fresh, so a Flag that's already true
// will fire the trigger on the rising-edge of the FIRST sample (we treat
// "no prior sample" as false).
func (g *Game) resetTriggers() {
	g.triggerStates = make(map[string]*triggerState)
	if g.currentScene == "" {
		return
	}
	s, ok := g.SceneManager.Get(g.currentScene)
	if !ok {
		return
	}
	for _, t := range s.Triggers {
		if t.Name == "" {
			continue
		}
		g.triggerStates[t.Name] = &triggerState{}
	}
}

// tickTriggers evaluates every armed trigger on the current scene. The
// engine calls this only when the script slot is idle (so we never
// preempt a running cutscene). The first matching rising-edge fires the
// trigger's Do action and consumes the script slot for this frame.
func (g *Game) tickTriggers() {
	if g.scriptRunner != nil || g.currentScene == "" {
		return
	}
	s, ok := g.SceneManager.Get(g.currentScene)
	if !ok {
		return
	}
	ctx := g.makeCtx()
	for _, t := range s.Triggers {
		if t.Name == "" || t.When == nil {
			continue
		}
		st := g.triggerStates[t.Name]
		if st == nil {
			st = &triggerState{}
			g.triggerStates[t.Name] = st
		}
		if st.fired && t.Once {
			continue
		}
		now := t.When.Eval(ctx)
		// Rising edge: was false, now true.
		if now && !st.lastTrue {
			st.fired = true
			st.lastTrue = true
			if t.Do != nil {
				g.queueAction(t.Do, "trigger "+t.Name)
				return // one trigger per frame keeps ordering deterministic
			}
			continue
		}
		st.lastTrue = now
	}
}
