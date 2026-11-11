package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type visibilityAction struct {
	character string
	hidden    bool
}

// Hide takes a character out of the picture until [Show] brings it back:
// it is not drawn and does not count as standing in the scene, in this
// scene and in every scene entered meanwhile. Its place, its walk and its
// lines go on as before. The character is identified by name; the root
// wrapper converts *Character to name.
func Hide(character string) inkbase.Action {
	return visibilityAction{character: character, hidden: true}
}

// Show puts a character taken out by [Hide] back into the picture. It does
// nothing to a character that is not hidden.
func Show(character string) inkbase.Action {
	return visibilityAction{character: character}
}

// Start returns the action itself; it is done on its first tick.
func (a visibilityAction) Start() inkbase.Runner { return a }

// Tick hides or shows the character and is done.
func (a visibilityAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	ctx.Game.(VisibilityHost).SetCharacterHidden(a.character, a.hidden)
	return inkbase.StatusDone
}
