package dialog

import "git.teletypegames.org/engines/inkwell/inkbase"

type Manager = inkbase.Manager[*Dialog]

func NewManager() *Manager { return inkbase.NewManager[*Dialog]() }
