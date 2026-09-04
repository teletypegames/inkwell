package inkwell

// Audible is anything the AudioPlayer can play: a file on the asset
// filesystem, named by Source. Voice, Music and Sound are the three, and
// they are three types rather than one because they are played differently
// — a Voice for as long as a line is spoken, a Music until another replaces
// it, a Sound once, over whatever else is playing.
type Audible interface {
	Source() string
}

// Voice is a recorded line. It sits on a Translation, because the English
// and the Hungarian recording are two files, and Say plays the one the
// resolved language carries — when the player has voices on. A Voice is
// owned by exactly one line, so it needs no name and no registration:
//
//	{Lang: "hu", Text: "Katt.", Voice: &Voice{Path: "assets/voice/alley/click.hu.ogg"}},
type Voice struct {
	Path string
}

func (v *Voice) Source() string { return v.Path }

// Music is a track a Scene or a Stage plays while it is showing. Loop says
// whether it starts over when it ends; a title theme loops, a stinger does
// not. Entering a Scene or Stage whose Music is the one already playing
// changes nothing.
type Music struct {
	Name string
	Path string
	Loop bool
}

func (m Music) GetName() string   { return m.Name }
func (m Music) TypeLabel() string { return "music" }
func (m *Music) Source() string   { return m.Path }

// Sound is a noise: a door, a latch, a cassette clicking home. PlaySound
// plays it once, over the music and over another Sound.
type Sound struct {
	Name string
	Path string
}

func (s Sound) GetName() string   { return s.Name }
func (s Sound) TypeLabel() string { return "sound" }
func (s *Sound) Source() string   { return s.Path }
