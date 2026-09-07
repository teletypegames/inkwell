package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

type Named = inkbase.Named
type Knowable = inkbase.Knowable
type Manager[T Named] = inkbase.Manager[T]

func NewManager[T Named]() *Manager[T] {
	return inkbase.NewManager[T]()
}
