// Package verb holds the [Verb] type, its [Manager] and the SCUMM verb set.
//
// A verb is a canonical Name ("look"), a Label in every language the game
// speaks, and a Default action run when a hotspot has no handler for it.
// [Defaults] is the built-in set — look, use, talk and take, in English and
// Hungarian — which inkwell.New registers unless the game hands over a
// VerbManager that already holds verbs. A game adds its own verbs by
// registering them, answers them on a hotspot through OnVerb, and changes a
// preset (a third language, a shorter label) with the manager's Set. The
// verb bar and the verb wheel read the manager every frame, so a verb added
// at runtime shows at once.
package verb
