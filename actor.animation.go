package pncdsl

// AnimationClip is a placeholder for sprite-sheet animation data. Not used
// for rendering yet — characters draw as flat colored rectangles in this
// milestone — but the field exists so domain code can reference it.
type AnimationClip struct {
	Frames     []Rectangle // source rects on the sprite sheet
	FrameTime  float64
	Loop       bool
}
