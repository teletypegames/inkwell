package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input is read once per Update() and exposes a tiny "consume on use" API
// so that one click can be authoritative — e.g. clicking on a hotspot
// while a Say is active should only skip the Say, not also fire the
// hotspot underneath.
type Input struct {
	mouseX, mouseY int
	leftPressed    bool
	leftConsumed   bool
	rightPressed   bool
	rightConsumed  bool
}

func newInput() *Input { return &Input{} }

func (i *Input) poll() {
	i.mouseX, i.mouseY = ebiten.CursorPosition()
	i.leftPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	i.leftConsumed = false
	i.rightPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	i.rightConsumed = false
}

func (i *Input) Pos() (int, int)     { return i.mouseX, i.mouseY }
func (i *Input) Point() Point        { return Point{X: float64(i.mouseX), Y: float64(i.mouseY)} }
func (i *Input) LeftClicked() bool   { return i.leftPressed && !i.leftConsumed }
func (i *Input) RightClicked() bool  { return i.rightPressed && !i.rightConsumed }
func (i *Input) ConsumeLeft()        { i.leftConsumed = true }
func (i *Input) ConsumeRight()       { i.rightConsumed = true }

// consumedClick is the Say-style "did the user click to skip?" probe.
func (i *Input) consumedClick() bool {
	if i.LeftClicked() {
		i.ConsumeLeft()
		return true
	}
	return false
}
