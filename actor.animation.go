package pncdsl

import "image"

// AnimationClip describes one named sprite-sheet animation. Frames are
// source rectangles into the character's Sprite asset; the renderer
// projects them onto the character's screen-space W×H footprint.
//
// FrameTime is the time (seconds) each frame is shown; Loop=true rewinds
// to the first frame after the last, Loop=false freezes on the final
// frame. A clip with zero frames or zero FrameTime is treated as static
// and the renderer falls back to the whole-sprite draw path.
type AnimationClip struct {
	Frames    []Rectangle
	FrameTime float64
	Loop      bool
}

// Reserved clip names the engine picks automatically:
//
//	"idle"  — selected when the character is not moving.
//	"walk"  — selected while the character has an active path.
//
// A character with neither name falls back to the static whole-sprite
// (or placeholder) draw.
const (
	clipIdle = "idle"
	clipWalk = "walk"
)

// tickAnimation advances a character's current animation clip and
// switches between idle/walk based on the moving flag.
func (g *Game) tickAnimation(c *runtimeChar, dt float64) {
	desired := clipIdle
	if c.moving {
		desired = clipWalk
	}
	clip, has := c.def.Animations[desired]
	if !has {
		// Fall back to the other clip if the desired one is missing —
		// a character with only "idle" defined keeps idling during walks.
		other := clipIdle
		if desired == clipIdle {
			other = clipWalk
		}
		if cl, ok := c.def.Animations[other]; ok {
			desired = other
			clip = cl
			has = true
		}
	}
	if !has || len(clip.Frames) == 0 || clip.FrameTime <= 0 {
		c.currentClip = ""
		c.clipElapsed = 0
		c.clipFrame = 0
		return
	}
	if c.currentClip != desired {
		c.currentClip = desired
		c.clipElapsed = 0
		c.clipFrame = 0
	}
	c.clipElapsed += dt
	for c.clipElapsed >= clip.FrameTime {
		c.clipElapsed -= clip.FrameTime
		c.clipFrame++
		if c.clipFrame >= len(clip.Frames) {
			if clip.Loop {
				c.clipFrame = 0
			} else {
				c.clipFrame = len(clip.Frames) - 1
				c.clipElapsed = 0
				break
			}
		}
	}
}

// currentFrameRect returns the source rectangle of the character's
// current animation frame, in image.Rectangle form so it can feed
// ebiten.Image.SubImage directly. ok=false when the character should be
// drawn without animation (no clip / empty clip / out-of-range frame).
func (c *runtimeChar) currentFrameRect() (image.Rectangle, bool) {
	if c.currentClip == "" {
		return image.Rectangle{}, false
	}
	clip, ok := c.def.Animations[c.currentClip]
	if !ok || c.clipFrame < 0 || c.clipFrame >= len(clip.Frames) {
		return image.Rectangle{}, false
	}
	r := clip.Frames[c.clipFrame]
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H)), true
}
