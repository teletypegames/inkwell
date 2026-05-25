package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// transition is a fade-to-black overlay used between scene swaps.
type transition struct {
	active   bool
	t        float64
	duration float64
	out      bool // true: fading to black; false: fading from black
	onMid    func()
}

func (t *transition) start(onMid func()) {
	t.active = true
	t.out = true
	t.t = 0
	t.duration = 0.25
	t.onMid = onMid
}

func (t *transition) update(dt float64) {
	if !t.active {
		return
	}
	t.t += dt
	if t.t >= t.duration {
		if t.out {
			if t.onMid != nil {
				t.onMid()
				t.onMid = nil
			}
			t.out = false
			t.t = 0
			return
		}
		t.active = false
	}
}

func (t *transition) draw(dst *ebiten.Image, w, h int) {
	if !t.active {
		return
	}
	alpha := t.t / t.duration
	if !t.out {
		alpha = 1 - alpha
	}
	if alpha < 0 {
		alpha = 0
	} else if alpha > 1 {
		alpha = 1
	}
	c := color.RGBA{0, 0, 0, uint8(alpha * 255)}
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(h), c, false)
}
