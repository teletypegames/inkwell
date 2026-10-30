package inkbase

// Script is a named action sequence.
type Script struct {
	Name    string
	Actions Action
}

// GetName returns the script's Name, its key in the registry.
func (s Script) GetName() string { return s.Name }

// TypeLabel returns "script", the kind used in log and panic messages.
func (s Script) TypeLabel() string { return "script" }
