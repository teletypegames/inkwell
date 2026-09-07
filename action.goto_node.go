package inkwell

import "git.teletypegames.org/engines/inkwell/act"

// GotoNode moves the open conversation to another of its nodes.
func GotoNode(node string) Action { return act.GotoNode(node) }
