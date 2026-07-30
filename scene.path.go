package inkwell

import "math"

// pathEps is the tolerance in world units (≈ pixels) used when comparing
// polygon vertices for adjacency. Authors typically snap walkbox corners
// to integer or half-pixel coordinates, so 0.5 is generous.
const pathEps = 0.5

// pathfind returns a sequence of waypoints from start to end that stays
// within the union of the given walkboxes. The last entry is always the
// effective destination (which may differ from end if end was outside
// every walkbox — in that case the closest boundary point is used).
//
// When boxes is empty, walkbox routing is off and pathfind degenerates
// to a straight line ([]Point{end}). The caller (walkCharacter) treats
// this as the legacy single-target behaviour.
func pathfind(start, end Point, boxes []Polygon) []Point {
	if len(boxes) == 0 {
		return []Point{end}
	}

	startIdx := containingPolygon(start, boxes)
	if startIdx < 0 {
		start, startIdx = nearestPolygonPoint(start, boxes)
		if startIdx < 0 {
			return []Point{end}
		}
	}
	endIdx := containingPolygon(end, boxes)
	if endIdx < 0 {
		end, endIdx = nearestPolygonPoint(end, boxes)
		if endIdx < 0 {
			return []Point{end}
		}
	}
	if startIdx == endIdx {
		return []Point{end}
	}

	adj := buildAdjacency(boxes)
	seq := bfsPolygons(startIdx, endIdx, adj, len(boxes))
	if seq == nil {
		// Disconnected components — give up and walk straight; the
		// caller's tickCharacters will at least move toward end.
		return []Point{end}
	}

	out := make([]Point, 0, len(seq))
	for i := 0; i < len(seq)-1; i++ {
		mid, ok := sharedEdgeMidpoint(boxes[seq[i]], boxes[seq[i+1]])
		if !ok {
			continue
		}
		out = append(out, mid)
	}
	out = append(out, end)
	return out
}

func containingPolygon(p Point, boxes []Polygon) int {
	for i, b := range boxes {
		if b.Contains(p) {
			return i
		}
	}
	return -1
}

// nearestPolygonPoint projects p onto the boundary of the nearest
// walkbox and returns the projection plus its polygon index.
func nearestPolygonPoint(p Point, boxes []Polygon) (Point, int) {
	bestIdx := -1
	bestPt := p
	bestD := math.Inf(1)
	for i, b := range boxes {
		pt, d := closestOnPolygonBoundary(p, b)
		if d < bestD {
			bestD = d
			bestPt = pt
			bestIdx = i
		}
	}
	return bestPt, bestIdx
}

func closestOnPolygonBoundary(p Point, b Polygon) (Point, float64) {
	if len(b.Points) < 2 {
		return p, math.Inf(1)
	}
	bestPt := b.Points[0]
	bestD := p.Dist(bestPt)
	n := len(b.Points)
	for i := 0; i < n; i++ {
		a := b.Points[i]
		c := b.Points[(i+1)%n]
		q := closestOnSegment(p, a, c)
		d := p.Dist(q)
		if d < bestD {
			bestD = d
			bestPt = q
		}
	}
	return bestPt, bestD
}

func closestOnSegment(p, a, b Point) Point {
	dx, dy := b.X-a.X, b.Y-a.Y
	denom := dx*dx + dy*dy
	if denom == 0 {
		return a
	}
	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / denom
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return Point{X: a.X + t*dx, Y: a.Y + t*dy}
}

// buildAdjacency maps each polygon index to the indices of polygons it
// shares an edge with. O(n²·k²) where k is avg edge count — fine for
// the handful of walkboxes per scene that adventure games use.
func buildAdjacency(boxes []Polygon) map[int][]int {
	adj := make(map[int][]int, len(boxes))
	for i := range boxes {
		for j := i + 1; j < len(boxes); j++ {
			if _, ok := sharedEdgeMidpoint(boxes[i], boxes[j]); ok {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}
	return adj
}

// sharedEdgeMidpoint returns the midpoint of an edge shared between two
// polygons, if any. Two edges are considered shared when their endpoints
// match (within pathEps) in either orientation — the typical case where
// neighbouring walkboxes are authored with a common boundary.
func sharedEdgeMidpoint(a, b Polygon) (Point, bool) {
	na, nb := len(a.Points), len(b.Points)
	if na < 2 || nb < 2 {
		return Point{}, false
	}
	for i := 0; i < na; i++ {
		ai, ai1 := a.Points[i], a.Points[(i+1)%na]
		for j := 0; j < nb; j++ {
			bj, bj1 := b.Points[j], b.Points[(j+1)%nb]
			if (pointsEq(ai, bj1) && pointsEq(ai1, bj)) ||
				(pointsEq(ai, bj) && pointsEq(ai1, bj1)) {
				return Point{X: (ai.X + ai1.X) / 2, Y: (ai.Y + ai1.Y) / 2}, true
			}
		}
	}
	return Point{}, false
}

func pointsEq(a, b Point) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx <= pathEps && dy <= pathEps
}

// bfsPolygons returns the shortest polygon-index sequence from src to
// dst over the adjacency graph, inclusive on both ends. nil if no
// path exists.
func bfsPolygons(src, dst int, adj map[int][]int, n int) []int {
	if src == dst {
		return []int{src}
	}
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}
	prev[src] = src
	queue := []int{src}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == dst {
			break
		}
		for _, nxt := range adj[cur] {
			if prev[nxt] != -1 {
				continue
			}
			prev[nxt] = cur
			queue = append(queue, nxt)
		}
	}
	if prev[dst] == -1 {
		return nil
	}
	// Walk back from dst.
	out := []int{dst}
	for cur := dst; cur != src; cur = prev[cur] {
		out = append([]int{prev[cur]}, out...)
	}
	return out
}
