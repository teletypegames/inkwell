package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Script is a named Action registered in the ScriptManager, so a game can
// start it by name or hand it to [RunScript], [Game.OnStart] or
// [Game.OnFinale]; see [inkbase.Script].
type Script = inkbase.Script
