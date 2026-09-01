package inkwell

// Named is implemented by every entity that can be registered into a Manager.
type Named interface {
	GetName() string
}

// Manager is the single registry shape used by every entity type. Each entity
// kind gets a named alias (ItemManager = Manager[Item], etc.). Insertion
// order is preserved — Names() and All() both hand it back.
type Manager[T Named] struct {
	items map[string]T
	order []string
}

func NewManager[T Named]() *Manager[T] {
	return &Manager[T]{items: make(map[string]T)}
}

// Register adds v to the registry and hands it back. Panics on empty or
// duplicate name — these are construction-time bugs, not runtime conditions.
//
// It returns what it registered so that a registration can be an expression
// rather than a statement, which is what lets an entity file be a
// declaration instead of a procedure:
//
//	var _ = SceneManager.Register(Scene{
//		Name: "alley",
//		…
//	})
//
// That form is not merely shorter than an init(): a package-level variable
// is initialised in dependency order, so the manager is built before the
// file that writes into it no matter which one the compiler reads first.
// The registration order stays the order the files are presented in, the
// same as with init().
//
// The result is safe to ignore, so the plain statement form still compiles.
func (m *Manager[T]) Register(v T) T {
	name := v.GetName()
	if name == "" {
		panic("inkwell: Register: empty Name on " + typeName(v))
	}
	if _, dup := m.items[name]; dup {
		panic("inkwell: Register: duplicate " + typeName(v) + " name: " + name)
	}
	m.items[name] = v
	m.order = append(m.order, name)
	return v
}

// Set replaces a registered entry in place, keeping its position in the
// registration order, and registers the entry when the name is new. Register
// panics on a duplicate; Set is the deliberate overwrite.
func (m *Manager[T]) Set(v T) {
	name := v.GetName()
	if _, ok := m.items[name]; !ok {
		m.Register(v)
		return
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
		panic("inkwell: MustGet: unknown name: " + name)
	}
	return v
}

func (m *Manager[T]) Has(name string) bool {
	_, ok := m.items[name]
	return ok
}

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
