package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

type Status = inkbase.Status

const (
	StatusRunning = inkbase.StatusRunning
	StatusDone    = inkbase.StatusDone
	StatusFailed  = inkbase.StatusFailed
)

type Ctx = inkbase.Ctx
type Action = inkbase.Action
type Runner = inkbase.Runner

// gameOf extracts the *Game from a Ctx whose Game field is set by the
// engine. Every adventure action calls this at the top of its Tick.
func gameOf(ctx *Ctx) *Game { return ctx.Game.(*Game) }
