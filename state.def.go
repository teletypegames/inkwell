package inkwell

type State struct {
	flags   map[string]bool
	vars    map[string]any
	visited map[string]int
	talked  map[string]int
}

func NewState() *State {
	return &State{
		flags:   make(map[string]bool),
		vars:    make(map[string]any),
		visited: make(map[string]int),
		talked:  make(map[string]int),
	}
}

func (s *State) Flag(name string) bool    { return s.flags[name] }
func (s *State) SetFlag(name string)      { s.flags[name] = true }
func (s *State) ClearFlag(name string)    { delete(s.flags, name) }
func (s *State) Var(name string) any      { return s.vars[name] }
func (s *State) SetVar(name string, v any) { s.vars[name] = v }
func (s *State) Visited(name string) int  { return s.visited[name] }
func (s *State) NoteVisit(name string)    { s.visited[name]++ }
func (s *State) Talked(node string) int   { return s.talked[node] }
func (s *State) NoteTalked(node string)   { s.talked[node]++ }
