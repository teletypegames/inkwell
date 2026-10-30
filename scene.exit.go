package inkwell

// ExitSide says which edge of the picture an Exit sits on when it declares
// no Area of its own.
type ExitSide int

// The edges an Exit can sit on.
const (
	ExitLeft  ExitSide = iota // off the left edge
	ExitRight                 // off the right edge
	ExitBack                  // into the depth of the picture
	ExitNear                  // out towards the camera
)

// Exit is a way out of a scene. It becomes a Hotspot with the exit cursor
// whose use is the journey; look and take fall back to Game.ExitLook and
// Game.ExitTake when the Exit says nothing of its own.
type Exit struct {
	// To is where the exit leads. Nil means back the way the player came.
	// A scene that leads back to the one that leads here is written with
	// Scene.Leads, after both are declared.
	To *Scene

	Label Text

	// Side places the hotspot on an edge of the scene area when Area is nil.
	Side ExitSide

	Area Shape

	// Needs is a flag the journey requires; Blocked is what happens without
	// it. Nil Blocked simply does nothing.
	Needs   string
	Blocked Action

	OnLook Action
	OnTake Action
}

// ExitName is the hotspot name an Exit gets: "exit:<scene>", or "exit:back".
func ExitName(to *Scene) string {
	if to == nil {
		return "exit:back"
	}
	return "exit:" + to.Name
}

func (e Exit) hotspot(g *Game) Hotspot {
	var travel Action = GoTo(e.To)
	if e.To == nil {
		travel = Back()
	}
	if e.Needs != "" {
		if e.Blocked != nil {
			travel = If(Flag(e.Needs), travel, e.Blocked)
		} else {
			travel = If(Flag(e.Needs), travel)
		}
	}
	area := e.Area
	if area == nil {
		area = e.Side.area(g.SceneArea())
	}
	look, take := e.OnLook, e.OnTake
	if look == nil && g.ExitLook != nil {
		look = g.ExitLook(e)
	}
	if take == nil && g.ExitTake != nil {
		take = g.ExitTake(e)
	}
	return Hotspot{
		Name:   ExitName(e.To),
		Label:  e.Label,
		Area:   area,
		Cursor: CursorExit,
		OnLook: look,
		OnUse:  travel,
		OnTake: take,
	}
}

const (
	exitEdgeWidth  = 0.0875
	exitEdgeInset  = 0.080
	exitEdgeHeight = 0.876
	exitDoorWidth  = 0.275
	exitBackInset  = 0.033
	exitBackHeight = 0.347
	exitNearHeight = 0.198
)

func (s ExitSide) area(r Rectangle) Shape {
	edge := r.W * exitEdgeWidth
	door := r.W * exitDoorWidth
	doorX := r.X + (r.W-door)/2
	switch s {
	case ExitLeft:
		return Rect(r.X, r.Y+r.H*exitEdgeInset, edge, r.H*exitEdgeHeight)
	case ExitRight:
		return Rect(r.X+r.W-edge, r.Y+r.H*exitEdgeInset, edge, r.H*exitEdgeHeight)
	case ExitBack:
		return Rect(doorX, r.Y+r.H*exitBackInset, door, r.H*exitBackHeight)
	default:
		return Rect(doorX, r.Y+r.H*(1-exitNearHeight), door, r.H*exitNearHeight)
	}
}

// resolveDerived asks every scene for its Derived exits once, before
// Validate, and appends them in place.
func (g *Game) resolveDerived() {
	for _, s := range g.SceneManager.All() {
		if s.Derived == nil {
			continue
		}
		if derived := s.Derived(g, s); len(derived) > 0 {
			s.Exits = append(s.Exits, derived...)
		}
	}
}
