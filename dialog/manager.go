package dialog

import "git.teletypegames.org/engines/inkwell/inkbase"

// Manager is the registry of dialogs, keyed by Dialog.Name and kept in
// registration order; see [inkbase.Manager].
type Manager = inkbase.Manager[*Dialog]

// NewManager returns an empty [Manager] of dialogs.
func NewManager() *Manager { return inkbase.NewManager[*Dialog]() }
