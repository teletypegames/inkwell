package inkwell

type learnAction struct{ entity Knowable }

// Learn marks an entity Known — the one a script has to say out loud, because
// the two rules the engine applies on its own do not reach it: a character
// is known once it has spoken, an item once it has been in the inventory. A
// name read off a list, a face seen through glass, is a Learn.
func Learn(entity Knowable) Action { return &learnAction{entity: entity} }

func (a *learnAction) Start() Runner { return a }

func (a *learnAction) Tick(ctx *Ctx) Status {
	ctx.Game.State.NoteKnown(a.entity)
	return StatusDone
}
