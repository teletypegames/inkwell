package inkbase

// Audible is anything an audio player can play: a file on the asset
// filesystem, named by Source.
type Audible interface {
	Source() string
}

// Voice is a recorded line. It sits on a Translation, because the English
// and the Hungarian recording are two files, and Say plays the one the
// resolved language carries.
type Voice struct {
	Path string
}

// Source returns the recording's Path, which makes a Voice [Audible].
func (v *Voice) Source() string { return v.Path }

// Music is a track a Scene or a Stage plays while it is showing.
type Music struct {
	Name string
	Path string
	Loop bool
}

// GetName returns the track's Name, its key in the registry.
func (m Music) GetName() string { return m.Name }

// TypeLabel returns "music", the kind used in log and panic messages.
func (m Music) TypeLabel() string { return "music" }

// Source returns the track's Path, which makes a *Music [Audible].
func (m *Music) Source() string { return m.Path }

// Sound is a noise: a door, a latch, a cassette clicking home.
type Sound struct {
	Name string
	Path string
}

// GetName returns the sound's Name, its key in the registry.
func (s Sound) GetName() string { return s.Name }

// TypeLabel returns "sound", the kind used in log and panic messages.
func (s Sound) TypeLabel() string { return "sound" }

// Source returns the sound's Path, which makes a *Sound [Audible].
func (s *Sound) Source() string { return s.Path }
