package inkwell

import "git.teletypegames.org/engines/inkwell/cond"

// Given is a Condition written as a literal: every field that is set has
// to hold, and an empty Given holds always.
type Given = cond.Given
