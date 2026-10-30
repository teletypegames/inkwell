package achievement

import "git.teletypegames.org/engines/inkwell/inkbase"

// Manager is the registry of achievements, keyed by Achievement.Name and kept in
// registration order; see [inkbase.Manager].
type Manager = inkbase.Manager[*Achievement]

// NewManager returns an empty [Manager] of achievements.
func NewManager() *Manager { return inkbase.NewManager[*Achievement]() }
