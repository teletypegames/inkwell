package inkwell

type runScriptAction struct{ script *Script }

// RunScript runs a named action sequence in place, as if its Actions were
// written here.
func RunScript(script *Script) Action    { return &runScriptAction{script: script} }
func (a *runScriptAction) Start() Runner { return &runScriptRunner{spec: a} }

type runScriptRunner struct {
	spec  *runScriptAction
	inner Runner
}

func (r *runScriptRunner) Tick(ctx *Ctx) Status {
	if r.inner == nil {
		if r.spec.script == nil || r.spec.script.Actions == nil {
			return StatusFailed
		}
		r.inner = r.spec.script.Actions.Start()
	}
	return r.inner.Tick(ctx)
}
