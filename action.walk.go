package inkwell

type walkAction struct {
	character string
	to        Point
}

// Walk sends a character to a point, routed around the scene's walkboxes,
// and holds the frame until they arrive.
func Walk(character string, to Point) Action { return &walkAction{character: character, to: to} }

func (a *walkAction) Start() Runner { return &walkRunner{spec: a} }

type walkRunner struct {
	spec    *walkAction
	started bool
}

func (r *walkRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.walkCharacter(r.spec.character, r.spec.to)
	}
	if ctx.Game.characterMoving(r.spec.character) {
		return StatusRunning
	}
	return StatusDone
}
