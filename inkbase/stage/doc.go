// Package stage holds [Stage]: one full-screen thing a game can be showing —
// a main menu, a settings screen, a splash card, or the world itself.
//
// Exactly one Stage shows at a time. A Stage can hold for some seconds and
// move on to its Next (a splash, a title sequence), optionally Skippable, and
// gives way through its Transition (nil cuts). OnEnter and OnLeave run as it
// comes and goes, and Music, when set, plays while it shows.
//
// A Stage is data; the host owns which one is showing. The action package's
// GoToStage and GoBack ask a StageHost to change it, and the condition
// package's OnStage and OnWorld ask a StageReader which one shows. The root
// inkwell package re-exports Stage, and its *Game is both hosts.
package stage
