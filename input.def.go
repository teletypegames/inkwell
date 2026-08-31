package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input is read once per Update() and exposes a tiny "consume on use" API
// so that one click can be authoritative — e.g. clicking on a hotspot
// while a Say is active should only skip the Say, not also fire the
// hotspot underneath.
//
// Keys work the same way. A widget that binds the arrow keys claims them,
// so a menu walking a list does not also drive whatever else listens for
// the same key on the same frame.
type Input struct {
	mouseX, mouseY int
	leftPressed    bool
	leftConsumed   bool
	rightPressed   bool
	rightConsumed  bool
	wheelX, wheelY float64
	keysPressed    []ebiten.Key
	keysConsumed   []ebiten.Key
}

func newInput() *Input { return &Input{} }

func (i *Input) poll() {
	i.mouseX, i.mouseY = ebiten.CursorPosition()
	i.leftPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	i.leftConsumed = false
	i.rightPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	i.rightConsumed = false
	i.wheelX, i.wheelY = ebiten.Wheel()
	i.keysPressed = inpututil.AppendJustPressedKeys(i.keysPressed[:0])
	i.keysConsumed = i.keysConsumed[:0]
}

func (i *Input) Pos() (int, int)    { return i.mouseX, i.mouseY }
func (i *Input) Point() Point       { return Point{X: float64(i.mouseX), Y: float64(i.mouseY)} }
func (i *Input) LeftClicked() bool  { return i.leftPressed && !i.leftConsumed }
func (i *Input) RightClicked() bool { return i.rightPressed && !i.rightConsumed }
func (i *Input) ConsumeLeft()       { i.leftConsumed = true }
func (i *Input) ConsumeRight()      { i.rightConsumed = true }

// Wheel is the mouse wheel's travel on this frame, in ebiten's own units.
// Nothing consumes it — a widget that scrolls checks that the pointer is
// over it first.
func (i *Input) Wheel() (float64, float64) { return i.wheelX, i.wheelY }

// KeyPressed is the keyboard's LeftClicked: the key went down on this
// frame and nobody has claimed it yet. It is edge-triggered, so holding a
// key answers true once.
func (i *Input) KeyPressed(key ebiten.Key) bool {
	for _, k := range i.keysConsumed {
		if k == key {
			return false
		}
	}
	for _, k := range i.keysPressed {
		if k == key {
			return true
		}
	}
	return false
}

// ConsumeKey claims a key for this frame, so later KeyPressed calls read
// false.
func (i *Input) ConsumeKey(key ebiten.Key) {
	i.keysConsumed = append(i.keysConsumed, key)
}

// ClaimKey is KeyPressed followed by ConsumeKey — the shape a widget that
// binds half a dozen keys wants, so its own navigation does not leak to
// whatever else reads the same key.
func (i *Input) ClaimKey(key ebiten.Key) bool {
	if !i.KeyPressed(key) {
		return false
	}
	i.ConsumeKey(key)
	return true
}

// consumedClick is the Say-style "did the user click to skip?" probe.
func (i *Input) consumedClick() bool {
	if i.LeftClicked() {
		i.ConsumeLeft()
		return true
	}
	return false
}
