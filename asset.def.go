package inkwell

// Asset is a picture on the asset filesystem: a scene's background, a
// Stage's, a sprite. Sound has its own three types — Voice, Music, Sound —
// under Audible.
type Asset struct {
	Name string
	Path string
}

func (a Asset) GetName() string   { return a.Name }
func (a Asset) TypeLabel() string { return "asset" }
