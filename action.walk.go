package inkwell

type walkAction struct {
	character *Character
	to        Point
}

// Walk sends a character to a point through the walkboxes and holds the
// script until it arrives.
func Walk(character *Character, to Point) Action {
	return &walkAction{character: character, to: to}
}

func (a *walkAction) Start() Runner { return &walkRunner{spec: a} }

type walkRunner struct {
	spec    *walkAction
	started bool
}

func (r *walkRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.walkCharacter(r.spec.character.Name, r.spec.to)
	}
	if ctx.Game.characterMoving(r.spec.character.Name) {
		return StatusRunning
	}
	return StatusDone
}
