package pncdsl

import "sort"

// Named is implemented by every entity that can be registered into a Manager.
type Named interface {
	GetName() string
}

// Manager is the single registry shape used by every entity type. Each entity
// kind gets a named alias (ItemManager = Manager[Item], etc.).
type Manager[T Named] struct {
	items map[string]T
}

func NewManager[T Named]() *Manager[T] {
	return &Manager[T]{items: make(map[string]T)}
}

// Register adds v to the registry. Panics on empty or duplicate name —
// these are construction-time bugs, not runtime conditions.
func (m *Manager[T]) Register(v T) {
	name := v.GetName()
	if name == "" {
		panic("pncdsl: Register: empty Name on " + typeName(v))
	}
	if _, dup := m.items[name]; dup {
		panic("pncdsl: Register: duplicate " + typeName(v) + " name: " + name)
	}
	m.items[name] = v
}

func (m *Manager[T]) Get(name string) (T, bool) {
	v, ok := m.items[name]
	return v, ok
}

func (m *Manager[T]) MustGet(name string) T {
	v, ok := m.items[name]
	if !ok {
		panic("pncdsl: MustGet: unknown name: " + name)
	}
	return v
}

func (m *Manager[T]) Has(name string) bool {
	_, ok := m.items[name]
	return ok
}

func (m *Manager[T]) Len() int { return len(m.items) }

func (m *Manager[T]) Names() []string {
	names := make([]string, 0, len(m.items))
	for n := range m.items {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (m *Manager[T]) Each(fn func(T)) {
	for _, n := range m.Names() {
		fn(m.items[n])
	}
}

func typeName(v any) string {
	type namer interface{ TypeLabel() string }
	if n, ok := v.(namer); ok {
		return n.TypeLabel()
	}
	return "entity"
}
