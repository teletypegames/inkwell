package item

import "git.teletypegames.org/engines/inkwell/inkbase"

type Manager = inkbase.Manager[*Item]

func NewManager() *Manager { return inkbase.NewManager[*Item]() }
