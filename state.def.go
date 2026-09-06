package inkwell

type State struct {
	flags   map[string]bool
	vars    map[string]any
	visited map[string]int
	talked  map[string]int
	known   map[string]int
}

func NewState() *State {
	return &State{
		flags:   make(map[string]bool),
		vars:    make(map[string]any),
		visited: make(map[string]int),
		talked:  make(map[string]int),
		known:   make(map[string]int),
	}
}

// Knowable is an entity the player can come to know — an Item, a Character,
// a Speaker — and so an entry a dossier or a codex screen can list. The kind
// is part of the key, so an item and a character may share a Name.
type Knowable interface {
	Named
	TypeLabel() string
}

func knownKey(entity Knowable) string { return entity.TypeLabel() + "/" + entity.GetName() }

// Known is the order in which the player came to know an entity, counting
// from 1, or 0 for one they have not — so a list sorted by it reads in the
// order of discovery. A character is known once it has spoken, an item once
// it has been in the inventory, and anything else once a Learn says so.
func (s *State) Known(entity Knowable) int { return s.known[knownKey(entity)] }

func (s *State) NoteKnown(entity Knowable) {
	key := knownKey(entity)
	if s.known[key] > 0 {
		return
	}
	s.known[key] = len(s.known) + 1
}

func (s *State) Flag(name string) bool     { return s.flags[name] }
func (s *State) SetFlag(name string)       { s.flags[name] = true }
func (s *State) ClearFlag(name string)     { delete(s.flags, name) }
func (s *State) Var(name string) any       { return s.vars[name] }
func (s *State) SetVar(name string, v any) { s.vars[name] = v }
func (s *State) Visited(name string) int   { return s.visited[name] }
func (s *State) NoteVisit(name string)     { s.visited[name]++ }
func (s *State) Talked(node string) int    { return s.talked[node] }
func (s *State) NoteTalked(node string)    { s.talked[node]++ }
