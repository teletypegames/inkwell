package achievement

import "git.teletypegames.org/engines/inkwell/inkbase"

// Achievement is something the player is recognised for: a name a platform
// knows it by, a title and a line for a screen that lists it, and the
// Condition that earns it. The engine evaluates When every frame the world is
// showing and reports the rising edge — a flag raised, an item picked up, a
// scene entered — to every extension that can take one (see
// inkwell.Unlocker). The engine keeps no record of its own: a platform that
// has one, Steam, answers idempotently, and a game without one has nothing
// to persist.
type Achievement struct {
	// Name is the key a platform knows the achievement by — the API name
	// entered on the Steamworks partner site, for one.
	Name string
	// Title and Text are what a screen shows for it, in every language the
	// game speaks; a platform shows its own copy of the same words.
	Title inkbase.Text
	Text  inkbase.Text
	// When earns it. Evaluated on the world stage, every frame.
	When inkbase.Condition
}

func (a Achievement) GetName() string   { return a.Name }
func (a Achievement) TypeLabel() string { return "achievement" }
