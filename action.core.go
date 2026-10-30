package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Status is what one tick of a Runner says about it; see [inkbase.Status].
type Status = inkbase.Status

// The outcomes of a tick.
const (
	StatusRunning = inkbase.StatusRunning // not finished; tick it again next frame
	StatusDone    = inkbase.StatusDone    // finished
	StatusFailed  = inkbase.StatusFailed  // gave up; a sequence stops here
)

// Ctx is what a Runner is handed each tick: the frame's time step, the
// State, and the *Game; see [inkbase.Ctx].
type Ctx = inkbase.Ctx

// Action is an immutable spec of work, started fresh into a Runner every
// time it runs; see [inkbase.Action].
type Action = inkbase.Action

// Runner is one running execution of an Action, ticked once a frame until
// it is done or failed; see [inkbase.Runner].
type Runner = inkbase.Runner

// Condition is a question about the game as it stands: a flag, an item in
// hand, the scene the player is in. Actions branch on one, widgets are shown
// by one, and a dialog choice is offered by one, so the same eight little
// types answer for all three.
type Condition = inkbase.Condition

// gameOf extracts the *Game from a Ctx whose Game field is set by the
// engine. Every adventure action calls this at the top of its Tick.
func gameOf(ctx *Ctx) *Game { return ctx.Game.(*Game) }

// ----- Custom action ------------------------------------------------------

type immediateAction struct {
	fn func(*Ctx) Status
}

func (a *immediateAction) Start() Runner { return a }
func (a *immediateAction) Tick(ctx *Ctx) Status {
	if a.fn == nil {
		return StatusDone
	}
	return a.fn(ctx)
}

// Custom wraps a user function as an Action. Returns Done on first tick
// unless the function itself returns StatusRunning.
func Custom(fn func(*Ctx) Status) Action { return &immediateAction{fn: fn} }

// ----- queue --------------------------------------------------------------

// Do queues an action. Queued actions run one at a time and in order —
// the engine starts the next one on the first frame the previous one is
// done, so a widget or a domain pump can hand work in at any moment.
func (g *Game) Do(a Action) {
	if a == nil {
		return
	}
	g.queue = append(g.queue, a)
}

func (g *Game) queueAction(a Action, label string) {
	if a == nil {
		return
	}
	logf("queueAction %s", label)
	g.Do(a)
}

// pumpQueue starts the next queued action when nothing else is running.
func (g *Game) pumpQueue() {
	if g.scriptRunner != nil || len(g.queue) == 0 {
		return
	}
	a := g.queue[0]
	g.queue = g.queue[1:]
	g.scriptRunner = a.Start()
	g.scriptCtx = g.makeCtx()
}
