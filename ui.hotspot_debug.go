package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HotspotDebug overlays a colored outline on every hotspot of the current
// scene. F1 toggles it at runtime; the default is off.
type HotspotDebug struct {
	Name    string
	When    Condition
	Enabled bool
	// ToggleKey, if non-zero, overrides the default F1 toggle.
	ToggleKey ebiten.Key
}

func (h *HotspotDebug) GetName() string        { return h.Name }
func (h *HotspotDebug) VisibleWhen() Condition { return h.When }
func (h *HotspotDebug) Layer() Layer           { return LayerScene }

func (h *HotspotDebug) Tick(ctx *UICtx) {
	key := h.ToggleKey
	if key == 0 {
		key = ebiten.KeyF1
	}
	if inpututil.IsKeyJustPressed(key) {
		h.Enabled = !h.Enabled
	}
}

func (h *HotspotDebug) Draw(dst *ebiten.Image, ctx *UICtx) {
	if !h.Enabled {
		return
	}
	g := ctx.Game
	if g.currentScene == "" {
		return
	}
	col := g.Theme().HotspotOutline
	for _, hs := range g.SceneHotspots(g.currentScene) {
		b := hs.Area.Bounds()
		vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, col, false)
		if label := hs.Label.String(); label != "" {
			g.DrawText(dst, label, int(b.X)+1, int(b.Y)+1, col)
		}
	}
}
