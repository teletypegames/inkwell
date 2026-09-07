package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

type Status = inkbase.Status

const (
	StatusRunning = inkbase.StatusRunning
	StatusDone    = inkbase.StatusDone
	StatusFailed  = inkbase.StatusFailed
)

// Ctx is the per-tick context passed to running actions.
type Ctx struct {
	Game    *Game
	DT      float64
	State   *State
	Scene   *Scene
	Hotspot *Hotspot
	Item    *Item
}

// Action is an immutable spec of work; Start makes a fresh Runner with
// state. Storing an Action in a struct field (e.g. Hotspot.OnUse) is safe
// because the engine always calls Start before ticking, so two invocations
// can never share mutable state.
type Action interface {
	Start() Runner
}

// Runner is one in-flight execution of an Action. Tick advances it one
// frame and returns whether it's still running, done, or failed.
type Runner interface {
	Tick(ctx *Ctx) Status
}
