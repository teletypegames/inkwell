package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"

	iwsprite "git.teletypegames.org/engines/inkwell/sprite"
)

// Sprite is a picture or a sequence of them cut from one sheet; see
// [iwsprite.Sprite].
type Sprite = iwsprite.Sprite

// Facing is which way a character looks, left or right; see
// [iwsprite.Facing].
type Facing = iwsprite.Facing

// Pose is what a character is doing as far as its sprite is concerned; see
// [iwsprite.Pose].
type Pose = iwsprite.Pose

// The facings and the poses a character is drawn in.
const (
	FacingRight = iwsprite.FacingRight // looking right
	FacingLeft  = iwsprite.FacingLeft  // looking left

	PoseIdle = iwsprite.PoseIdle // standing; the fallback for a pose with no sprite
	PoseWalk = iwsprite.PoseWalk // moving
	PoseTalk = iwsprite.PoseTalk // while its line is up
	PoseWait = iwsprite.PoseWait // after standing still a while
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
