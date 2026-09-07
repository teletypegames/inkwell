package inkbase

// Asset is a picture on the asset filesystem: a scene's background, a
// stage's, a sprite.
type Asset struct {
	Name string
	Path string
}

func (a Asset) GetName() string   { return a.Name }
func (a Asset) TypeLabel() string { return "asset" }

// SceneLayer is one plane of a scene's picture at a depth Z. The cast
// stands at 0, so negative Z is behind every character and positive Z is in
// front.
type SceneLayer struct {
	Asset *Asset
	Z     int
}
