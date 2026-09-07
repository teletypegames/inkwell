package inkbase

import "math"

type Point struct {
	X, Y float64
}

func (p Point) Add(q Point) Point { return Point{X: p.X + q.X, Y: p.Y + q.Y} }
func (p Point) Dist(q Point) float64 {
	dx, dy := p.X-q.X, p.Y-q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Shape is a hotspot area. Rectangle and Polygon both satisfy it.
type Shape interface {
	Contains(p Point) bool
	Bounds() Rectangle
}

type Rectangle struct {
	X, Y, W, H float64
}

func Rect(x, y, w, h float64) Rectangle { return Rectangle{X: x, Y: y, W: w, H: h} }

func (r Rectangle) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

func (r Rectangle) Bounds() Rectangle { return r }

type Polygon struct {
	Points []Point
}

func Poly(pts ...Point) Polygon { return Polygon{Points: pts} }

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
