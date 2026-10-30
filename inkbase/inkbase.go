package inkbase

// Named is implemented by every entity that can be registered into a Manager.
type Named interface {
	GetName() string
}

// Knowable is an entity the player can come to know — an Item, a Character,
// a Speaker — and so an entry a dossier or a codex screen can list. The kind
// is part of the key, so an item and a character may share a Name.
type Knowable interface {
	Named
	TypeLabel() string
}

// Manager is the single registry shape used by every entity type. Insertion
// order is preserved — Names() and All() both hand it back.
type Manager[T Named] struct {
	items map[string]T
	order []string
}

// NewManager returns an empty registry for entities of type T.
func NewManager[T Named]() *Manager[T] {
	return &Manager[T]{items: make(map[string]T)}
}

// Register adds v to the registry and hands it back. Panics on empty or
// duplicate name — these are construction-time bugs, not runtime conditions.
//
// It returns what it registered so that a registration can be an expression
// rather than a statement, which is what lets an entity file be a
// declaration instead of a procedure.
func (m *Manager[T]) Register(v T) T {
	name := v.GetName()
	if name == "" {
		panic("inkbase: Register: empty Name on " + typeName(v))
	}
	if _, dup := m.items[name]; dup {
		panic("inkbase: Register: duplicate " + typeName(v) + " name: " + name)
	}
	m.items[name] = v
	m.order = append(m.order, name)
	return v
}

// Set replaces a registered entry in place, keeping its position in the
// registration order, and registers the entry when the name is new.
func (m *Manager[T]) Set(v T) {
	name := v.GetName()
	if _, ok := m.items[name]; !ok {
		m.Register(v)
		return
	}
	m.items[name] = v
}

// Get returns the entry registered under name and true, or the zero T and
// false when there is none.
func (m *Manager[T]) Get(name string) (T, bool) {
	v, ok := m.items[name]
	return v, ok
}

// MustGet returns the entry registered under name. It panics when there is
// none, so use it only for names the game itself registered.
func (m *Manager[T]) MustGet(name string) T {
	v, ok := m.items[name]
	if !ok {
		panic("inkbase: MustGet: unknown name: " + name)
	}
	return v
}

// Has reports whether an entry is registered under name.
func (m *Manager[T]) Has(name string) bool {
	_, ok := m.items[name]
	return ok
}

// Len returns the number of registered entries.
func (m *Manager[T]) Len() int { return len(m.items) }

// Names returns the registered names in insertion order.
func (m *Manager[T]) Names() []string {
	return append([]string(nil), m.order...)
}

// All returns every registered entry in insertion order.
func (m *Manager[T]) All() []T {
	out := make([]T, len(m.order))
	for i, n := range m.order {
		out[i] = m.items[n]
	}
	return out
}

// Each calls fn on every registered entry in insertion order.
func (m *Manager[T]) Each(fn func(T)) {
	for _, n := range m.order {
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
