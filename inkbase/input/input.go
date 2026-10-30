package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

// Input is read once per Update() and exposes a "consume on use" API so that
// one click can be authoritative.
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

// New returns an empty Input; call [Input.Poll] at the start of each
// Update before reading it.
func New() *Input { return &Input{} }

// Poll reads the cursor position, this frame's fresh mouse clicks, the
// wheel and the newly pressed keys from Ebitengine, and clears every
// consumed mark from the previous frame. Call it once per Update.
func (i *Input) Poll() {
	i.mouseX, i.mouseY = ebiten.CursorPosition()
	i.leftPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	i.leftConsumed = false
	i.rightPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	i.rightConsumed = false
	i.wheelX, i.wheelY = ebiten.Wheel()
	i.keysPressed = inpututil.AppendJustPressedKeys(i.keysPressed[:0])
	i.keysConsumed = i.keysConsumed[:0]
}

// Pos returns the cursor position in screen pixels as of the last Poll.
func (i *Input) Pos() (int, int) { return i.mouseX, i.mouseY }

// Point returns the cursor position as of the last Poll as an
// [inkbase.Point].
func (i *Input) Point() inkbase.Point {
	return inkbase.Point{X: float64(i.mouseX), Y: float64(i.mouseY)}
}

// LeftClicked reports whether the left button went down this frame and
// nobody has consumed the click yet. It does not consume it.
func (i *Input) LeftClicked() bool { return i.leftPressed && !i.leftConsumed }

// RightClicked reports whether the right button went down this frame and
// nobody has consumed the click yet. It does not consume it.
func (i *Input) RightClicked() bool { return i.rightPressed && !i.rightConsumed }

// ConsumeLeft claims this frame's left click, so [Input.LeftClicked]
// reports false until the next Poll.
func (i *Input) ConsumeLeft() { i.leftConsumed = true }

// ConsumeRight claims this frame's right click, so [Input.RightClicked]
// reports false until the next Poll.
func (i *Input) ConsumeRight() { i.rightConsumed = true }

// Wheel is the mouse wheel's travel on this frame.
func (i *Input) Wheel() (float64, float64) { return i.wheelX, i.wheelY }

// KeyPressed is edge-triggered: the key went down this frame and nobody
// has claimed it yet.
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

// ConsumeKey claims a key for this frame.
func (i *Input) ConsumeKey(key ebiten.Key) {
	i.keysConsumed = append(i.keysConsumed, key)
}

// ClaimKey is KeyPressed followed by ConsumeKey.
func (i *Input) ClaimKey(key ebiten.Key) bool {
	if !i.KeyPressed(key) {
		return false
	}
	i.ConsumeKey(key)
	return true
}

// ConsumedClick probes whether the user clicked to skip, consuming the click.
func (i *Input) ConsumedClick() bool {
	if i.LeftClicked() {
		i.ConsumeLeft()
		return true
	}
	return false
}
