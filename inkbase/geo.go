package inkbase

import "math"

// Point is a position in screen-space pixels.
type Point struct {
	X, Y float64
}

// Add returns the sum of p and q, p moved by q.
func (p Point) Add(q Point) Point { return Point{X: p.X + q.X, Y: p.Y + q.Y} }

// Dist returns the straight-line distance between p and q.
func (p Point) Dist(q Point) float64 {
	dx, dy := p.X-q.X, p.Y-q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Shape is a hotspot area. Rectangle and Polygon both satisfy it.
type Shape interface {
	Contains(p Point) bool
	Bounds() Rectangle
}

// Rectangle is an axis-aligned box: its top-left corner X, Y and its size
// W, H, in screen-space pixels.
type Rectangle struct {
	X, Y, W, H float64
}

// Rect returns the Rectangle with top-left corner x, y and size w, h.
func Rect(x, y, w, h float64) Rectangle { return Rectangle{X: x, Y: y, W: w, H: h} }

// Contains reports whether p lies in r. The left and top edges belong to
// the rectangle, the right and bottom ones do not.
func (r Rectangle) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

// Bounds returns r itself.
func (r Rectangle) Bounds() Rectangle { return r }

// Polygon is a closed shape through Points in order; the last point joins
// back to the first.
type Polygon struct {
	Points []Point
}

// Poly returns the Polygon through pts.
func Poly(pts ...Point) Polygon { return Polygon{Points: pts} }

// Contains reports whether pt lies inside g, by the even-odd rule. A
// polygon of fewer than three points contains nothing.
func (g Polygon) Contains(pt Point) bool {
	n := len(g.Points)
	if n < 3 {
		return false
	}
	inside := false
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		pi, pj := g.Points[i], g.Points[j]
		if (pi.Y > pt.Y) != (pj.Y > pt.Y) &&
			pt.X < (pj.X-pi.X)*(pt.Y-pi.Y)/(pj.Y-pi.Y)+pi.X {
			inside = !inside
		}
	}
	return inside
}

// Bounds returns the smallest Rectangle around every point of g, or the
// zero Rectangle when g has no points.
func (g Polygon) Bounds() Rectangle {
	if len(g.Points) == 0 {
		return Rectangle{}
	}
	minX, maxX := g.Points[0].X, g.Points[0].X
	minY, maxY := g.Points[0].Y, g.Points[0].Y
	for _, p := range g.Points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return Rectangle{X: minX, Y: minY, W: maxX - minX, H: maxY - minY}
}

// PolyRect is the rectangle x,y,w,h as a four-cornered Polygon.
func PolyRect(x, y, w, h float64) Polygon {
	return Poly(
		Point{X: x, Y: y},
		Point{X: x + w, Y: y},
		Point{X: x + w, Y: y + h},
		Point{X: x, Y: y + h},
	)
}
