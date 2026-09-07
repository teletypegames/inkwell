package theme

import "git.teletypegames.org/engines/inkwell/inkbase"

type Manager = inkbase.Manager[*Theme]

func NewManager() *Manager { return inkbase.NewManager[*Theme]() }
