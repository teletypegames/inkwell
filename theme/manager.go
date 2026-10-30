package theme

import "git.teletypegames.org/engines/inkwell/inkbase"

// Manager is the registry of themes, keyed by Theme.Name and kept in
// registration order; see [inkbase.Manager].
type Manager = inkbase.Manager[*Theme]

// NewManager returns an empty [Manager] of themes.
func NewManager() *Manager { return inkbase.NewManager[*Theme]() }
