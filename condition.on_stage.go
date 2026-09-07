package inkwell

import "git.teletypegames.org/engines/inkwell/cond"

// OnStage is true while the given Stage is the one showing.
func OnStage(stage *Stage) Condition { return cond.OnStage(stage) }

// OnWorld is true while the Stage showing is the world -- whichever Stage
// that is.
func OnWorld() Condition { return cond.OnWorld() }

// Started is true once a session exists: the world has been begun or loaded.
func Started() Condition { return cond.Started() }
