package inkbase

// Status is the outcome of one tick of a Runner.
type Status int

// The outcomes of a tick.
const (
	StatusRunning Status = iota // not finished; tick it again next frame
	StatusDone                  // finished; a sequence moves on
	StatusFailed                // gave up; a sequence stops here
)

// Ctx is the per-tick context passed to running actions. State is available
// directly — most generic actions only need it. Game is opaque; an action
// needing more (audio, stages, themes) type-asserts it to a narrow host
// interface. Adventure actions assert it to *Game.
type Ctx struct {
	DT    float64
	State *State
	Game  any
}

// Action is an immutable spec of work; Start makes a fresh Runner with
// state. Storing an Action in a struct field is safe because the engine
// always calls Start before ticking, so two invocations can never share
// mutable state.
type Action interface {
	Start() Runner
}

// Runner is one in-flight execution of an Action. Tick advances it one
// frame and returns whether it's still running, done, or failed.
type Runner interface {
	Tick(ctx *Ctx) Status
}
