package inkwell

type runScriptAction struct{ name string }

// RunScript runs a registered script as a single action. An unknown name,
// or a script with nothing in it, fails the action rather than the game.
func RunScript(name string) Action       { return &runScriptAction{name: name} }
func (a *runScriptAction) Start() Runner { return &runScriptRunner{spec: a} }

type runScriptRunner struct {
	spec  *runScriptAction
	inner Runner
}

func (r *runScriptRunner) Tick(ctx *Ctx) Status {
	if r.inner == nil {
		s, ok := ctx.Game.ScriptManager.Get(r.spec.name)
		if !ok || s.Actions == nil {
			return StatusFailed
		}
		r.inner = s.Actions.Start()
	}
	return r.inner.Tick(ctx)
}
