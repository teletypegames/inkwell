package inkbase

// Script is a named Action, registered so that a game can start it by name
// or hand it to the RunScript action (see the action package).
type Script struct {
	// Name is the script's key in its registry.
	Name string
	// Actions is what the script runs.
	Actions Action
}

// GetName returns the script's Name, its key in the registry.
func (s Script) GetName() string { return s.Name }

// TypeLabel returns "script", the kind a Script is registered as and the
// word used in log and panic messages.
func (s Script) TypeLabel() string { return "script" }
