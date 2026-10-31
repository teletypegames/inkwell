package inkwell

import ibstage "git.teletypegames.org/engines/inkwell/inkbase/stage"

// Stage is one full-screen thing the game can be showing — a main menu, a
// settings screen, a splash card, or the world itself; see [ibstage.Stage].
//
// That is the difference between a Stage and a Scene. A Scene is a place in
// the fiction, and the engine always has one; walking from a room to an alley
// is a Scene change inside the world. A Stage is what owns the frame. Off the
// world Stage the scene is not drawn, characters do not walk, triggers do not
// fire and hotspots answer nothing — which is what a game showing its main
// menu wants. Widgets need no new field to follow this: [OnStage] is a
// Condition like any other, so a HUD hides itself off the world the same way
// it hides itself for a cutscene.
type Stage = ibstage.Stage
