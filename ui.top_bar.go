package inkwell

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TopBar is a thin strip across the top of the screen with three sections:
// scene title on the left, a score number in the center, and a day/time
// string on the right. Empty fields are skipped.
type TopBar struct {
	Name   string
	Height int

	// LeftText overrides the auto-discovered scene Title/Name.
	LeftText string
	// ScoreVar is a State Var key whose value is fmt-printed as the
	// center score; empty = no score shown.
	ScoreVar string
	ScoreMax int // when >0, rendered as "Score: X/MAX"

	// TimeVar is a State Var key whose value is printed at the right
	// edge (e.g. "Day 2 - 14:32"); empty = nothing on the right.
	TimeVar string
}

func (t *TopBar) GetName() string { return t.Name }
func (t *TopBar) Tick(ctx *UICtx) {}

func (t *TopBar) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	th := g.Theme()
	h := t.Height
	if h <= 0 {
		h = 12
	}
	vector.DrawFilledRect(dst, 0, 0, float32(g.Width), float32(h), th.TopBarBG, false)

	left := t.LeftText
	if left == "" && g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		if s.Title != "" {
			left = s.Title
		} else {
			left = s.Name
		}
	}
	if left != "" {
		g.DrawText(dst, left, 4, -1, th.TopBarText)
	}

	if t.ScoreVar != "" {
		var sc string
		v := g.State.Var(t.ScoreVar)
		if t.ScoreMax > 0 {
			sc = fmt.Sprintf("Score: %v/%d", v, t.ScoreMax)
		} else {
			sc = fmt.Sprintf("Score: %v", v)
		}
		x := (g.Width - textWidth(sc)) / 2
		g.DrawText(dst, sc, x, -1, th.TopBarAccent)
	}

	if t.TimeVar != "" {
		tm := fmt.Sprint(g.State.Var(t.TimeVar))
		x := g.Width - textWidth(tm) - 4
		g.DrawText(dst, tm, x, -1, th.TopBarText)
	}
}

// BlocksClickAt — the TopBar swallows clicks so the radial menu doesn't
// pop up under it.
func (t *TopBar) BlocksClickAt(p Point) bool {
	h := t.Height
	if h <= 0 {
		h = 12
	}
	return p.Y < float64(h)
}
