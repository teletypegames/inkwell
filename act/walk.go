package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type walkAction struct {
	character string
	to        inkbase.Point
}

// Walk sends a character to a point through the walkboxes and holds the
// script until it arrives. The character is identified by name; the root
// wrapper converts *Character to name.
func Walk(character string, to inkbase.Point) inkbase.Action {
	return &walkAction{character: character, to: to}
}

func (a *walkAction) Start() inkbase.Runner { return &walkRunner{spec: a} }

type walkRunner struct {
	spec    *walkAction
	started bool
}

func (r *walkRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(MovementHost)
	if !r.started {
		r.started = true
		host.WalkCharacter(r.spec.character, r.spec.to)
	}
	if host.CharacterMoving(r.spec.character) {
		return inkbase.StatusRunning
	}
	return inkbase.StatusDone
}
