package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// sceneTransition is what a scene change goes through; a Stage names its
// own in Stage.Transition.
var sceneTransition = Fade(0.25)

// transition is the one Transition in flight: which one, how far along, and
// the work to do at the darkest frame.
type transition struct {
	active bool
	spec   Transition
	t      float64
	out    bool // true: covering the old picture; false: uncovering the new
	onMid  func()
}

func (t *transition) start(spec Transition, onMid func()) {
	t.active = true
	t.spec = spec
	t.out = true
	t.t = 0
	t.onMid = onMid
}

func (t *transition) update(dt float64) {
	if !t.active {
		return
	}
	t.t += dt
	if t.t >= t.spec.Seconds() {
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
	cover := t.t / t.spec.Seconds()
	if !t.out {
		cover = 1 - cover
	}
	if cover < 0 {
		cover = 0
	} else if cover > 1 {
		cover = 1
	}
	t.spec.Draw(dst, w, h, cover)
}
