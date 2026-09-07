package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type runScriptAction struct{ script *inkbase.Script }

// RunScript runs a named action sequence in place.
func RunScript(script *inkbase.Script) inkbase.Action { return &runScriptAction{script: script} }
func (a *runScriptAction) Start() inkbase.Runner      { return &runScriptRunner{spec: a} }

type runScriptRunner struct {
	spec  *runScriptAction
	inner inkbase.Runner
}

func (r *runScriptRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if r.inner == nil {
		if r.spec.script == nil || r.spec.script.Actions == nil {
			return inkbase.StatusFailed
		}
		r.inner = r.spec.script.Actions.Start()
	}
	return r.inner.Tick(ctx)
}
