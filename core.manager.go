package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Named is anything that can be registered in a Manager: it answers its
// name; see [inkbase.Named].
type Named = inkbase.Named

// Knowable is an entity the player can come to know — an Item, a Character,
// a Speaker — and so an entry a codex screen can list; see
// [inkbase.Knowable].
type Knowable = inkbase.Knowable

// Manager is the registry every entity kind is kept in, keyed by name and
// in registration order; see [inkbase.Manager].
type Manager[T Named] = inkbase.Manager[T]

// NewManager returns an empty Manager for entities of type T; see
// [inkbase.NewManager].
func NewManager[T Named]() *Manager[T] {
	return inkbase.NewManager[T]()
}
