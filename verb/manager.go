package verb

import "git.teletypegames.org/engines/inkwell/inkbase"

type Manager = inkbase.Manager[*Verb]

func NewManager() *Manager { return inkbase.NewManager[*Verb]() }
