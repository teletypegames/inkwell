package inkwell

type runtimeChar struct {
	def *Character
	pos Point
	// path is the queue of remaining waypoints, head-first. Empty when
	// idle. The last entry equals target.
	path []Point
	// target is the final destination — kept separately so saves and the
	// post-Load resume in tickCharacters can rebuild a missing path.
	target Point
	moving bool

	// animation playback state — see character.animation.go for the rules
	// that drive currentClip selection.
	pose        Pose
	posed       bool
	poseElapsed float64
	still       float64
	facing      Facing
}

// sceneActors is the cast a scene actually shows: its own list, or the
// player alone when the scene names nobody.
func (g *Game) sceneActors(s *Scene) []SceneActor {
	if s.Actors != nil || g.Player == nil {
		return s.Actors
	}
	return []SceneActor{{
		Character: g.Player,
		At:        g.Player.Start,
	}}
}

func (g *Game) placeActors(sceneName string) {
	s := g.SceneManager.MustGet(sceneName)
	for _, a := range g.sceneActors(s) {
		def := a.Character
		rc, ok := g.chars[def.Name]
		if !ok {
			rc = &runtimeChar{def: def, pos: a.At}
			g.chars[def.Name] = rc
		} else {
			rc.def = def
			rc.pos = a.At
			rc.moving = false
		}
	}
}

func (g *Game) runtimeChar(name string) (*runtimeChar, bool) {
	c, ok := g.chars[name]
	return c, ok
}

func (g *Game) walkCharacter(name string, to Point) {
	c, ok := g.chars[name]
	if !ok {
		return
	}
	var boxes []Polygon
	if g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		boxes = s.Walkboxes
		if boxes == nil {
			boxes = g.Walkboxes
		}
	}
	path := pathfind(c.pos, to, boxes)
	if len(path) == 0 {
		c.moving = false
		return
	}
	c.path = path
	c.target = path[len(path)-1]
	c.moving = true
	c.posed = false
}

func (g *Game) characterMoving(name string) bool {
	c, ok := g.chars[name]
	return ok && c.moving
}

func (g *Game) tickCharacters(dt float64) {
	for _, c := range g.chars {
		g.tickPose(c, dt)
		if !c.moving {
			continue
		}
		// Resume after Load: moving with no live path means we lost the
		// route but still know the final destination; rebuild a trivial
		// straight-line path so the character finishes its walk.
		if len(c.path) == 0 {
			if c.pos == c.target {
				c.moving = false
				continue
			}
			c.path = []Point{c.target}
		}
		next := c.path[0]
		dx := next.X - c.pos.X
		dy := next.Y - c.pos.Y
		if dx < 0 {
			c.facing = FacingLeft
		} else if dx > 0 {
			c.facing = FacingRight
		}
		d := c.pos.Dist(next)
		speed := c.def.Speed
		if speed <= 0 {
			speed = 60
		}
		step := speed * dt
		if d <= step {
			c.pos = next
			c.path = c.path[1:]
			if len(c.path) == 0 {
				c.moving = false
			}
			continue
		}
		c.pos.X += dx / d * step
		c.pos.Y += dy / d * step
	}
}
