package inkwell

// Script is a named Action registered in the ScriptManager, so a game can
// start it by name or hand it to [RunScript], [Game.OnStart] or
// [Game.OnFinale].
type Script struct {
	Name    string
	Actions Action
}

// GetName returns the script's Name, its key in the ScriptManager.
func (s Script) GetName() string { return s.Name }

// TypeLabel returns "script", the kind a Script is registered as.
func (s Script) TypeLabel() string { return "script" }
