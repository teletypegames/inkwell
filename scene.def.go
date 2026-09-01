package inkwell

type Scene struct {
	Name       string
	Title      Text   // optional human-readable title for the TopBar widget
	Background string // Asset.Name
	Music      string // Asset.Name (optional)
	Hotspots   []Hotspot
	Exits      []Exit // connections to other scenes; see scene.exit.go
	Walkboxes  []Polygon
	Triggers   []Trigger
	Actors     []SceneActor
	OnEnter    Action
	OnLeave    Action

	// Derived hands back exits worked out from the rest of the game
	// rather than written into this scene — a map screen whose pins are
	// every location that says it leads here, say. Run asks each scene
	// once, before validation, and appends what comes back to Exits, so
	// everything downstream sees one kind of exit and the scene file
	// stays a literal instead of being patched from the wiring.
	Derived func(g *Game) []Exit
}

func (s Scene) GetName() string   { return s.Name }
func (s Scene) TypeLabel() string { return "scene" }

// BackgroundAsset is the image the scene paints: its Background, or its own
// name when it does not say. A game whose backgrounds are named after the
// scenes they belong to therefore writes the pairing down once — in the
// scene's name — instead of in a second constant, a second file and a
// field that repeats both.
func (s Scene) BackgroundAsset() string {
	if s.Background != "" {
		return s.Background
	}
	return s.Name
}

// SceneActor places a registered Character at a starting position inside a scene.
type SceneActor struct {
	CharacterName string
	At            Point
}
