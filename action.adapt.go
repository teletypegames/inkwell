package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// adaptAction wraps an inkbase.Action so it can run in inkwell's action
// queue — which passes inkwell.Ctx, not inkbase.Ctx. The wrapper builds
// an inkbase.Ctx from the fields the two share (DT, State) and passes
// the *Game as the opaque Game field, so an inkbase action that
// type-asserts to a host interface finds the methods it needs.
func adaptAction(inner inkbase.Action) Action {
	if inner == nil {
		return nil
	}
	return &adaptedAction{inner: inner}
}

type adaptedAction struct{ inner inkbase.Action }

func (a *adaptedAction) Start() Runner {
	return &adaptedRunner{inner: a.inner.Start()}
}

type adaptedRunner struct{ inner inkbase.Runner }

func (r *adaptedRunner) Tick(ctx *Ctx) Status {
	ibctx := &inkbase.Ctx{
		DT:    ctx.DT,
		State: ctx.State,
		Game:  ctx.Game,
	}
	return r.inner.Tick(ibctx)
}
