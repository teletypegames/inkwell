package verb

import "git.teletypegames.org/engines/inkwell/inkbase"

// Manager is the registry of verbs, keyed by Verb.Name and kept in
// registration order; see [inkbase.Manager].
type Manager = inkbase.Manager[*Verb]

// NewManager returns an empty [Manager] of verbs.
func NewManager() *Manager { return inkbase.NewManager[*Verb]() }
