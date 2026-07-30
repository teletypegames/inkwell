package inkwell

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// CharStat is one (Label, VarKey) row rendered inside a CharacterPanel.
// VarKey is read from g.State.Var; the value is fmt.Sprint'd verbatim.
type CharStat struct {
	Label  string
	VarKey string
}

// CharacterPanel is a small floating panel that shows portrait + stats for
// one registered character. The panel auto-hides when the character isn't
// listed as an actor in the current scene, so the same panel registration
// works across scenes.
type CharacterPanel struct {
	Name      string
	Bounds    Rectangle
	Character string // Character.Name to display
	Title     string // "PLAYER", "NPC", role badge
	Stats     []CharStat
}

func (c *CharacterPanel) GetName() string { return c.Name }
func (c *CharacterPanel) Tick(ctx *UICtx) {}

func (c *CharacterPanel) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	if c.Character == "" || !g.CharacterInScene(c.Character) {
		return
	}
	th := g.Theme()
	b := c.Bounds

	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), th.CharacterPanelBG, false)
	vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, th.CharacterPanelBorder, false)

	// portrait box on the left
	portraitW := 22.0
	px := b.X + 2
	py := b.Y + 2
	drawCharacterPortrait(dst, g, c.Character, px, py, portraitW, b.H-4)

	// text column
	tx := int(b.X + portraitW + 6)
	ty := int(b.Y) + 1

	if c.Title != "" {
		g.DrawText(dst, c.Title, tx, ty, th.CharacterPanelTitle)
		ty += 10
	}

	if ch, ok := g.CharacterManager.Get(c.Character); ok && ch.Name != "" {
		g.DrawText(dst, ch.Name, tx, ty, th.StatusText)
		ty += 10
	}

	for _, st := range c.Stats {
		v := g.State.Var(st.VarKey)
		line := fmt.Sprintf("%s: %v", st.Label, v)
		g.DrawText(dst, line, tx, ty, th.StatusText)
		ty += 10
	}
}

// BlocksClickAt — the panel sits on top of the scene; clicks inside it
// shouldn't trigger hotspots underneath.
func (c *CharacterPanel) BlocksClickAt(p Point) bool {
	return c.Bounds.Contains(p)
}

// drawCharacterPortrait renders a tiny version of the character's
// placeholder shape (or sprite if available) inside the panel.
func drawCharacterPortrait(dst *ebiten.Image, g *Game, name string, x, y, w, h float64) {
	c, ok := g.CharacterManager.Get(name)
	if !ok {
		return
	}
	if c.Sprite != "" && !g.loaded.isPlaceholder(c.Sprite) {
		img := g.loaded.image(g.AssetManager, c.Sprite)
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		if sw > 0 && sh > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(w/float64(sw), h/float64(sh))
			op.GeoM.Translate(x, y)
			dst.DrawImage(img, op)
			return
		}
	}
	bodyCol := color.RGBA{200, 200, 200, 255}
	if rgba, ok := c.SpeechColor.(color.RGBA); ok {
		bodyCol = rgba
	}
	if c.W > c.H && c.W > 0 {
		drawQuadrupedPlaceholder(dst, x, y, w, h, bodyCol)
	} else {
		drawHumanoidPlaceholder(dst, x, y, w, h, bodyCol)
	}
}
