package inkwell

type runtimeChar struct {
	def Character
	pos Point
	// path is the queue of remaining waypoints, head-first. Empty when
	// idle. The last entry equals target.
	path []Point
	// target is the final destination — kept separately so saves and the
	// post-Load resume in tickCharacters can rebuild a missing path.
	target Point
	moving bool

	// animation playback state — see actor.animation.go for the rules
	// that drive currentClip selection.
	currentClip string
	clipElapsed float64
	clipFrame   int
}

// sceneActors is the cast a scene actually shows: its own list, or the
// player alone when the scene names nobody.
func (g *Game) sceneActors(s Scene) []SceneActor {
	if s.Actors != nil || g.Player == "" {
		return s.Actors
	}
	def, ok := g.CharacterManager.Get(g.Player)
	if !ok {
		return nil
	}
	return []SceneActor{{
		CharacterName: def.Name,
		At:            def.Start,
	}}
}

func (g *Game) placeActors(sceneName string) {
	s := g.SceneManager.MustGet(sceneName)
	for _, a := range g.sceneActors(s) {
		def := g.CharacterManager.MustGet(a.CharacterName)
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
}

func (g *Game) characterMoving(name string) bool {
	c, ok := g.chars[name]
	return ok && c.moving
}

func (g *Game) tickCharacters(dt float64) {
	for _, c := range g.chars {
		g.tickAnimation(c, dt)
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
