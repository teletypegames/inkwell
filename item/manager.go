package item

import "git.teletypegames.org/engines/inkwell/inkbase"

// Manager is the registry of items, keyed by Item.Name and kept in
// registration order; see [inkbase.Manager].
type Manager = inkbase.Manager[*Item]

// NewManager returns an empty [Manager] of items.
func NewManager() *Manager { return inkbase.NewManager[*Item]() }
