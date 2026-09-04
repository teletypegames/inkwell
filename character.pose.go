package inkwell

// tickPose decides what a character is doing this frame and advances its
// sprite. A pose held by SetPose wins until the character moves or speaks;
// otherwise Walk while moving, Talk while its line is up, Wait once it has
// stood still for waitAfter seconds, and Idle in between.
func (g *Game) tickPose(c *runtimeChar, dt float64) {
	if c.moving || g.talking == c.def.Name {
		c.posed = false
		c.still = 0
	} else {
		c.still += dt
	}
	desired := PoseIdle
	switch {
	case c.posed:
		desired = c.pose
	case c.moving:
		desired = PoseWalk
	case g.talking == c.def.Name:
		desired = PoseTalk
	case c.still >= waitAfter:
		desired = PoseWait
	}
	if _, has := c.def.Poses[desired]; !has {
		desired = PoseIdle
	}
	if desired != c.pose {
		c.pose = desired
		c.poseElapsed = 0
	} else {
		c.poseElapsed += dt
	}
}

// sprite is the Sprite for the character's current pose, or nil.
func (c *runtimeChar) sprite() *Sprite {
	return c.def.Poses[c.pose]
}
