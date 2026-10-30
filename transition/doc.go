// Package transition holds the [Transition] contract and the fade.
//
// A Transition is how one picture gives way to another, in two halves: the
// old picture is covered, the change lands on the frame it is fully hidden,
// and the new picture is uncovered. Seconds is the length of one half; Draw
// paints the cover over the frame for a cover value that runs from 0 to 1
// while the old picture goes and back to 0 while the new one comes.
//
// [Fade] covers with black and [FadeTo] with any colour. A wipe or an iris is
// another implementation of the same two methods.
//
// In an inkwell game a Stage names its own transition (Stage.Transition; nil
// is a cut), and a scene change uses a quarter-second fade. At the midpoint
// of a scene change the old scene's OnLeave is queued, the scene changes,
// the visit is counted, the actors are placed, the new scene's music starts
// and its OnEnter is queued. Input is frozen while the old picture is going.
// The opening scene is entered without a transition, so the opening script
// runs at once.
package transition
