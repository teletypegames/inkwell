package inkbase

// State is the game's progress: flags, variables, how often each scene was
// visited and each dialog node talked, and the entities the player has come
// to know. It is what a save file holds; the zero value is not usable, make
// one with [NewState].
type State struct {
	flags   map[string]bool
	vars    map[string]any
	visited map[string]int
	talked  map[string]int
	known   map[string]int
}

// NewState returns an empty State with all of its maps made.
func NewState() *State {
	return &State{
		flags:   make(map[string]bool),
		vars:    make(map[string]any),
		visited: make(map[string]int),
		talked:  make(map[string]int),
		known:   make(map[string]int),
	}
}

func knownKey(entity Knowable) string { return entity.TypeLabel() + "/" + entity.GetName() }

// Known is the order in which the player came to know an entity, counting
// from 1, or 0 for one they have not.
func (s *State) Known(entity Knowable) int { return s.known[knownKey(entity)] }

// NoteKnown records that the player now knows entity, giving it the next
// number in the order [State.Known] reports. An entity already known keeps
// its number.
func (s *State) NoteKnown(entity Knowable) {
	key := knownKey(entity)
	if s.known[key] > 0 {
		return
	}
	s.known[key] = len(s.known) + 1
}

// Flag reports whether the flag name is set.
func (s *State) Flag(name string) bool { return s.flags[name] }

// SetFlag sets the flag name.
func (s *State) SetFlag(name string) { s.flags[name] = true }

// ClearFlag unsets the flag name.
func (s *State) ClearFlag(name string) { delete(s.flags, name) }

// Var returns the variable name, or nil when it was never set.
func (s *State) Var(name string) any { return s.vars[name] }

// SetVar sets the variable name to v.
func (s *State) SetVar(name string, v any) { s.vars[name] = v }

// Visited returns how many times the scene name has been entered.
func (s *State) Visited(name string) int { return s.visited[name] }

// NoteVisit counts one more visit to the scene name.
func (s *State) NoteVisit(name string) { s.visited[name]++ }

// Talked returns how many times the dialog node has been played.
func (s *State) Talked(node string) int { return s.talked[node] }

// NoteTalked counts one more playing of the dialog node.
func (s *State) NoteTalked(node string) { s.talked[node]++ }

// Flags returns a copy of the underlying flag map — for serialization.
func (s *State) Flags() map[string]bool {
	out := make(map[string]bool, len(s.flags))
	for k, v := range s.flags {
		out[k] = v
	}
	return out
}

// Vars returns a copy of the underlying var map — for serialization.
func (s *State) Vars() map[string]any {
	out := make(map[string]any, len(s.vars))
	for k, v := range s.vars {
		out[k] = v
	}
	return out
}

// VisitedMap returns a copy of the visited map — for serialization.
func (s *State) VisitedMap() map[string]int {
	out := make(map[string]int, len(s.visited))
	for k, v := range s.visited {
		out[k] = v
	}
	return out
}

// TalkedMap returns a copy of the talked map — for serialization.
func (s *State) TalkedMap() map[string]int {
	out := make(map[string]int, len(s.talked))
	for k, v := range s.talked {
		out[k] = v
	}
	return out
}

// KnownMap returns a copy of the known map — for serialization.
func (s *State) KnownMap() map[string]int {
	out := make(map[string]int, len(s.known))
	for k, v := range s.known {
		out[k] = v
	}
	return out
}

// LoadFlags replaces the flag map from a saved state.
func (s *State) LoadFlags(flags map[string]bool) {
	s.flags = make(map[string]bool, len(flags))
	for k, v := range flags {
		s.flags[k] = v
	}
}

// LoadVars replaces the var map from a saved state.
func (s *State) LoadVars(vars map[string]any) {
	s.vars = make(map[string]any, len(vars))
	for k, v := range vars {
		s.vars[k] = v
	}
}

// LoadVisited replaces the visited map from a saved state.
func (s *State) LoadVisited(visited map[string]int) {
	s.visited = make(map[string]int, len(visited))
	for k, v := range visited {
		s.visited[k] = v
	}
}

// LoadTalked replaces the talked map from a saved state.
func (s *State) LoadTalked(talked map[string]int) {
	s.talked = make(map[string]int, len(talked))
	for k, v := range talked {
		s.talked[k] = v
	}
}

// LoadKnown replaces the known map from a saved state.
func (s *State) LoadKnown(known map[string]int) {
	s.known = make(map[string]int, len(known))
	for k, v := range known {
		s.known[k] = v
	}
}
