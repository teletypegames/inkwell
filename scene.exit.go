package inkwell

// ExitSide is where on the picture an exit sits when it has no Area of its
// own: off one edge, into the depth of the picture, or out towards the
// camera. Until the doors are measured on the artwork, an edge strip is the
// convention every 1990s point & click used, and re-aiming an exit at a real
// door later is a one-field change.
type ExitSide int

const (
	ExitLeft  ExitSide = iota // off the left edge
	ExitRight                 // off the right edge
	ExitBack                  // into the depth of the picture
	ExitNear                  // out towards the camera
)

// Exit is one way out of a scene. The engine turns it into a Hotspot, so a
// connection can be declared as data — where it leads, what the status line
// calls it, and what has to be true before it opens — instead of being spelled
// out as another hand-built hotspot.
//
// A connection stays machine-readable: the generated hotspot is named
// ExitName(To), so the location graph can be read straight back out of the
// registered scenes.
type Exit struct {
	// To is the scene the exit leads to. Empty means "back the way you came" —
	// for a scene reachable from several others, where the way out is a
	// direction rather than a place.
	To string

	// Label is what the status line calls it.
	Label Text

	// Side places the exit when Area is nil.
	Side ExitSide

	// Area overrides the edge strip Side would give it.
	Area Shape

	// Needs is a flag that has to be set before the exit opens; Blocked is
	// what happens while it is not. The exit is visible either way — a locked
	// door still says it is a door.
	Needs   string
	Blocked Action

	// OnLook and OnTake override Game.ExitLook and Game.ExitTake for this one
	// exit.
	OnLook Action
	OnTake Action
}

// ExitName is the hotspot name generated for an exit to the named scene.
func ExitName(to string) string {
	if to == "" {
		return "exit:back"
	}
	return "exit:" + to
}

func (e Exit) hotspot(g *Game) Hotspot {
	var travel Action = GoTo(e.To)
	if e.To == "" {
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

// The edge strips, as fractions of the scene area, so a game that hands over
// only part of the window to the picture (a HUD along the foot, letterbox
// bars) gets strips that sit inside the picture rather than under the HUD.
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
