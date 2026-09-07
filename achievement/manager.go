package achievement

import "git.teletypegames.org/engines/inkwell/inkbase"

type Manager = inkbase.Manager[*Achievement]

func NewManager() *Manager { return inkbase.NewManager[*Achievement]() }
