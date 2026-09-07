package action

import "git.teletypegames.org/engines/inkwell/inkbase"

type learnAction struct{ entity inkbase.Knowable }

// Learn marks an entity Known.
func Learn(entity inkbase.Knowable) inkbase.Action { return &learnAction{entity: entity} }

func (a *learnAction) Start() inkbase.Runner { return a }

func (a *learnAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	ctx.State.NoteKnown(a.entity)
	return inkbase.StatusDone
}
