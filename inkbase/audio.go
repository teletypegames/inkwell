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

func (v *Voice) Source() string { return v.Path }

// Music is a track a Scene or a Stage plays while it is showing.
type Music struct {
	Name string
	Path string
	Loop bool
}

func (m Music) GetName() string   { return m.Name }
func (m Music) TypeLabel() string { return "music" }
func (m *Music) Source() string   { return m.Path }

// Sound is a noise: a door, a latch, a cassette clicking home.
type Sound struct {
	Name string
	Path string
}

func (s Sound) GetName() string   { return s.Name }
func (s Sound) TypeLabel() string { return "sound" }
func (s *Sound) Source() string   { return s.Path }
