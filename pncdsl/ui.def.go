package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Layout constants in the internal 320×200 coordinate space.
const (
	uiSceneTop    = 0
	uiSceneBottom = 140
	uiStatusY     = 142
	uiPanelY      = 152
	uiPanelH      = 48
	uiVerbsX      = 4
	uiVerbsW      = 120
	uiInvX        = 132
	uiInvW        = 184
	uiInvSlotW    = 22
)

// UI is the per-game UI state. It's a simple bag of fields the engine
// writes into and the renderer reads. No event/signal layer — the engine
// drives everything top-down each tick.
type UI struct {
	g *Game

	hoverLabel string
	flash      string
	flashTimer float64

	speech         *speechBubble
	dialog         *dialogBox
	endCard        string

	verbButtons []verbButton
}

func newUI(g *Game) *UI {
	return &UI{
		g:      g,
		speech: &speechBubble{},
	}
}

func (u *UI) rebuildVerbButtons() {
	verbs := u.g.VerbManager.Names()
	u.verbButtons = u.verbButtons[:0]
	cols := 2
	bw := uiVerbsW / cols
	bh := 14
	for i, name := range verbs {
		col := i % cols
		row := i / cols
		x := uiVerbsX + col*bw
		y := uiPanelY + row*bh
		u.verbButtons = append(u.verbButtons, verbButton{
			Name:   name,
			Label:  u.g.VerbManager.MustGet(name).Label,
			Bounds: Rect(float64(x), float64(y), float64(bw-2), float64(bh-2)),
		})
	}
}

type verbButton struct {
	Name   string
	Label  string
	Bounds Rectangle
}

// SetSpeech / ClearSpeech are called from the Say action.
func (u *UI) SetSpeech(speaker, text string) {
	u.speech.set(u.g, speaker, text)
}

func (u *UI) ClearSpeech() { u.speech.clear() }

// FlashLine shows a short status string (e.g. "Ehhez kell a kulcs.") for ~2s.
func (u *UI) FlashLine(text string) {
	u.flash = text
	u.flashTimer = 2.0
}

func (u *UI) tick(dt float64) {
	if u.flashTimer > 0 {
		u.flashTimer -= dt
		if u.flashTimer <= 0 {
			u.flash = ""
		}
	}
	u.speech.tick(dt)
}

// drawHUD paints the bottom UI panel (status line + verbs + inventory).
// Dialog box, if active, is drawn separately by the engine on top.
func (u *UI) drawHUD(dst *ebiten.Image) {
	// status / hover / flash line
	status := u.flash
	if status == "" {
		status = u.composedHoverLabel()
	}
	if status != "" {
		w := textWidth(status)
		drawText(dst, status, (320-w)/2, uiStatusY, color.White)
	}

	// panel background
	vector.DrawFilledRect(dst, 0, float32(uiPanelY-2), 320, float32(uiPanelH+2), color.RGBA{20, 20, 30, 255}, false)

	// verb buttons
	for _, b := range u.verbButtons {
		bg := color.RGBA{50, 50, 70, 255}
		if b.Name == u.g.selectedVerb {
			bg = color.RGBA{120, 90, 40, 255}
		}
		vector.DrawFilledRect(dst, float32(b.Bounds.X), float32(b.Bounds.Y), float32(b.Bounds.W), float32(b.Bounds.H), bg, false)
		drawText(dst, b.Label, int(b.Bounds.X)+2, int(b.Bounds.Y)-1, color.White)
	}

	// inventory slots
	items := u.g.Inventory.Items()
	for k := 0; k < 8; k++ {
		x := uiInvX + k*uiInvSlotW
		y := uiPanelY
		bg := color.RGBA{40, 40, 50, 255}
		if k < len(items) && items[k] == u.g.Inventory.Selected() {
			bg = color.RGBA{120, 90, 40, 255}
		}
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(uiInvSlotW-2), float32(uiInvSlotW-2), bg, false)
		if k < len(items) {
			it, ok := u.g.ItemManager.Get(items[k])
			if ok {
				img := u.g.loaded.image(u.g.AssetManager, it.Sprite)
				op := &ebiten.DrawImageOptions{}
				sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
				if sw == 0 || sh == 0 {
					continue
				}
				slotPx := float64(uiInvSlotW - 4)
				sx := slotPx / float64(sw)
				sy := slotPx / float64(sh)
				op.GeoM.Scale(sx, sy)
				op.GeoM.Translate(float64(x+1), float64(y+1))
				dst.DrawImage(img, op)
			}
		}
	}
}

func (u *UI) composedHoverLabel() string {
	verb := ""
	if v, ok := u.g.VerbManager.Get(u.g.selectedVerb); ok {
		verb = v.Label
	}
	target := u.hoverLabel
	if sel := u.g.Inventory.Selected(); sel != "" {
		if target != "" {
			return verb + " " + sel + " ezen: " + target
		}
		return verb + " " + sel
	}
	if target == "" {
		return ""
	}
	return verb + " " + target
}

// hitTestUI returns the click target (verbButton, inventory slot index) or "".
// "verb:<name>", "inv:<i>", "" for none.
func (u *UI) hitTest(p Point) string {
	for _, b := range u.verbButtons {
		if b.Bounds.Contains(p) {
			return "verb:" + b.Name
		}
	}
	items := u.g.Inventory.Items()
	for k := 0; k < len(items); k++ {
		x := uiInvX + k*uiInvSlotW
		r := Rect(float64(x), float64(uiPanelY), float64(uiInvSlotW-2), float64(uiInvSlotW-2))
		if r.Contains(p) {
			return "inv:" + items[k]
		}
	}
	return ""
}
