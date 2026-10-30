package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Audible is anything the AudioPlayer can play, named by its Source path;
// see [inkbase.Audible].
type Audible = inkbase.Audible

// Voice is the recording of one line in one language, carried by a
// Translation; see [inkbase.Voice].
type Voice = inkbase.Voice

// Music is a track a Scene or a Stage plays while it shows, looped when its
// Loop is set; see [inkbase.Music].
type Music = inkbase.Music

// Sound is a one-shot noise, played by [PlaySound]; see [inkbase.Sound].
type Sound = inkbase.Sound
