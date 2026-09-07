package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"

	iwsprite "git.teletypegames.org/engines/inkwell/sprite"
)

type Sprite = iwsprite.Sprite
type Facing = iwsprite.Facing
type Pose = iwsprite.Pose

const (
	FacingRight = iwsprite.FacingRight
	FacingLeft  = iwsprite.FacingLeft

	PoseIdle = iwsprite.PoseIdle
	PoseWalk = iwsprite.PoseWalk
	PoseTalk = iwsprite.PoseTalk
	PoseWait = iwsprite.PoseWait
)

// waitAfter re-exports the constant so root files (character.pose.go) keep
// their lowercase spelling.
const waitAfter = iwsprite.WaitAfter

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
	r := s.FrameRect(img, frame)
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
