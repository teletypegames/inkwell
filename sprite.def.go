package inkwell

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is a picture or a sequence of them: a sheet, and the frames cut
// from it. No Frames means the whole sheet is the one frame, so a still
// item and a walking character are the same type — one has one frame, the
// other has eight and a FrameTime.
type Sprite struct {
	Sheet *Asset
	// Frames are rectangles on the sheet, in playing order. Empty means one
	// frame, the sheet itself.
	Frames []Rectangle
	// FrameTime is how long one frame shows, in seconds. 0 with several
	// frames shows the first and never advances.
	FrameTime float64
	// Loop starts over after the last frame; off, the last frame holds.
	Loop bool
	// Faces is the way the artwork looks. A character walking the other way
	// is drawn mirrored, so one sheet serves both directions.
	Faces Facing
}

// Facing is left or right; the sheet is drawn as it is when the character
// faces the way the sheet does, and mirrored otherwise.
type Facing int

const (
	FacingRight Facing = iota
	FacingLeft
)

// Pose is what a character is doing, as far as its sprite is concerned. The
// engine picks it: Walk while moving, Talk while its line is up, Wait after
// standing still a while, Idle otherwise — or whatever SetPose held. A pose
// a character has no Sprite for falls back to Idle, and a character with no
// Idle is drawn as the placeholder figure.
type Pose int

const (
	PoseIdle Pose = iota
	PoseWalk
	PoseTalk
	PoseWait
)

// waitAfter is how long a character stands still before Wait replaces Idle.
const waitAfter = 8.0

func (s *Sprite) frameCount() int {
	if s == nil {
		return 0
	}
	if len(s.Frames) == 0 {
		return 1
	}
	return len(s.Frames)
}

// frameAt is the frame showing after elapsed seconds from the start.
func (s *Sprite) frameAt(elapsed float64) int {
	count := s.frameCount()
	if count <= 1 || s.FrameTime <= 0 {
		return 0
	}
	frame := int(elapsed / s.FrameTime)
	if s.Loop {
		return frame % count
	}
	if frame >= count {
		return count - 1
	}
	return frame
}

func (s *Sprite) frameRect(img *ebiten.Image, frame int) image.Rectangle {
	if len(s.Frames) == 0 || frame < 0 || frame >= len(s.Frames) {
		return img.Bounds()
	}
	r := s.Frames[frame]
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H)).Intersect(img.Bounds())
}

// drawSprite paints one frame of a sprite into the box at x, y of w by h,
// mirrored when the facing is not the sheet's. It reports false when there
// is nothing to draw — no sprite, or a sheet that could not be read — so the
// caller can fall back to a placeholder.
func drawSprite(dst *ebiten.Image, g *Game, s *Sprite, frame int, facing Facing, x, y, w, h float64) bool {
	if s == nil || s.Sheet == nil || g.loaded.isPlaceholder(s.Sheet) {
		if s != nil && s.Sheet != nil {
			_ = g.loaded.image(s.Sheet)
		}
		return false
	}
	img := g.loaded.image(s.Sheet)
	r := s.frameRect(img, frame)
	fw, fh := r.Dx(), r.Dy()
	if fw <= 0 || fh <= 0 {
		return false
	}
	sub := img.SubImage(r).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	if facing != s.Faces {
		op.GeoM.Scale(-w/float64(fw), h/float64(fh))
		op.GeoM.Translate(x+w, y)
	} else {
		op.GeoM.Scale(w/float64(fw), h/float64(fh))
		op.GeoM.Translate(x, y)
	}
	dst.DrawImage(sub, op)
	return true
}
