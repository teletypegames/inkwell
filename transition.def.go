package inkwell

import iwtransition "git.teletypegames.org/engines/inkwell/inkbase/transition"

// Transition is how one Stage gives way to the next, such as [Fade]; nil
// is a cut. See [iwtransition.Transition].
type Transition = iwtransition.Transition
