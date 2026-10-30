package inkwell

// Store is where a game keeps what has to outlive the process: its settings
// and its save slots. A key is a short word — "settings", "slots" — and a
// value is the bytes of one document; the Store decides what that means on
// the platform it runs on. Natively a key is a JSON file in the user's
// configuration directory, under the game's App name; in a browser it is
// one localStorage entry. Read returns ErrNotStored for a key nothing has
// written yet, so a first run is not an error.
type Store interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}

// MemoryStore keeps documents in a map, for a test or a run that must not
// touch the disk.
type MemoryStore struct {
	docs map[string][]byte
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{docs: make(map[string][]byte)}
}

// Read returns a copy of the document under key, or ErrNotStored if
// nothing has been written there.
func (m *MemoryStore) Read(key string) ([]byte, error) {
	data, ok := m.docs[key]
	if !ok {
		return nil, ErrNotStored
	}
	return append([]byte(nil), data...), nil
}

// Write keeps a copy of data under key, replacing what was there. It never
// fails.
func (m *MemoryStore) Write(key string, data []byte) error {
	m.docs[key] = append([]byte(nil), data...)
	return nil
}

func (g *Game) store() Store {
	if g.Store == nil {
		g.Store = defaultStore(g.AppName())
	}
	return g.Store
}

// AppName is App, or "inkwell" for a game that named none: the word the
// platform's default Store files the documents under, exported so that a
// Store an extension builds files them under the same one.
func (g *Game) AppName() string {
	if g.App != "" {
		return g.App
	}
	return "inkwell"
}
