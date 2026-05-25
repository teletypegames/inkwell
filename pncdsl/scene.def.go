package pncdsl

type Scene struct {
	Name       string
	Title      string // optional human-readable title for the TopBar widget
	Background string // Asset.Name
	Music      string // Asset.Name (optional)
	Hotspots   []Hotspot
	Walkboxes  []Polygon
	Triggers   []Trigger
	Actors     []SceneActor
	OnEnter    Action
	OnLeave    Action
}

func (s Scene) GetName() string   { return s.Name }
func (s Scene) TypeLabel() string { return "scene" }

// SceneActor places a registered Character at a starting position inside a scene.
type SceneActor struct {
	CharacterName string
	At            Point
}
