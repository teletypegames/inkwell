package inkwell

import (
	"git.teletypegames.org/engines/inkwell/cond"
	ibcond "git.teletypegames.org/engines/inkwell/inkbase/condition"
)

// --- delegated to cond/ (adventure conditions) ----------------------------

// Given is a Condition written as a literal: every field that is set has
// to hold, and an empty Given holds always.
type Given = cond.Given

// HasItem is true while the item is in the inventory.
func HasItem(item *Item) Condition { return cond.HasItem(item.Name) }

// SelectedItem is true while the item is the one in hand.
func SelectedItem(item *Item) Condition { return cond.SelectedItem(item.Name) }

// InScene is true while the given scene is the current one.
func InScene(scene *Scene) Condition { return cond.InScene(scene.Name) }

// OnStage is true while the given Stage is the one showing.
func OnStage(stage *Stage) Condition { return cond.OnStage(stage) }

// OnWorld is true while the Stage showing is the world.
func OnWorld() Condition { return cond.OnWorld() }

// Started is true once a session exists: the world has been begun or loaded.
func Started() Condition { return cond.Started() }

// --- delegated to inkbase/condition/ (generic conditions) -----------------

// Flag holds while a named world-state flag is raised.
func Flag(name string) Condition { return ibcond.Flag(name) }

// VarEq holds while a world-state variable compares equal to a value.
func VarEq(name string, v any) Condition { return ibcond.VarEq(name, v) }

// Knows is true once the player has come to know an entity.
func Knows(entity Knowable) Condition { return ibcond.Knows(entity) }

// --- combinators (delegated to inkbase/condition/) ------------------------

// And holds when all of its conditions hold.
func And(conditions ...Condition) Condition { return ibcond.And(conditions...) }

// Or holds when at least one of its conditions holds.
func Or(conditions ...Condition) Condition { return ibcond.Or(conditions...) }

// Not inverts a condition.
func Not(c Condition) Condition { return ibcond.Not(c) }
