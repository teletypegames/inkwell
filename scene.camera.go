package inkwell

// Camera is a stub identity transform — no scrolling in this milestone.
type Camera struct {
	Offset Point
}

func NewCamera() *Camera { return &Camera{} }

func (c *Camera) Apply(p Point) Point { return Point{p.X - c.Offset.X, p.Y - c.Offset.Y} }
