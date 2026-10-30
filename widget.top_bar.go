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
	When   Condition
	Height int

	// LeftText overrides the auto-discovered scene Title/Name.
	LeftText Text
	// ScoreVar is a State Var key whose value is fmt-printed as the
	// center score; empty = no score shown.
	ScoreVar string
	ScoreMax int // when >0, rendered as "Score: X/MAX"
	// ScoreLabel overrides the engine's own word for it.
	ScoreLabel Text

	// TimeVar is a State Var key whose value is printed at the right
	// edge (e.g. "Day 2 - 14:32"); empty = nothing on the right.
	TimeVar string

	// NoteVar is a State Var key whose value is appended to the left
	// text after an em dash — "ALLEY — paused". Empty or unset = no
	// note, and the title stands alone.
	NoteVar string

	// Notes is what to print for each value NoteVar takes. The value
	// itself is game state a Condition compares against, so it has to
	// stay the same in every language; what the player reads does not.
	// A value with no entry prints as it stands.
	Notes map[string]Text
}

// GetName returns the name the bar is registered under.
func (t *TopBar) GetName() string { return t.Name }

// VisibleWhen returns the condition that switches the bar on; nil means
// always.
func (t *TopBar) VisibleWhen() Condition { return t.When }

// Layer puts the bar on [LayerHUD].
func (t *TopBar) Layer() Layer { return LayerHUD }

// Tick does nothing; the bar reads the game when it draws.
func (t *TopBar) Tick(ctx *UICtx) {}

// Draw paints the strip, Height pixels tall (12 when zero), with the left
// text and its note, the score in the centre and the time at the right
// edge, each read from State as the frame is drawn.
func (t *TopBar) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	th := g.Theme()
	h := t.Height
	if h <= 0 {
		h = 12
	}
	vector.DrawFilledRect(dst, 0, 0, float32(g.Width), float32(h), th.TopBarBG, false)

	left := t.LeftText.String()
	if left == "" && g.currentScene != "" {
		s := g.SceneManager.MustGet(g.currentScene)
		if left = s.Title.String(); left == "" {
			left = s.Name
		}
	}
	if t.NoteVar != "" {
		if note := varText(g.State.Var(t.NoteVar)); note != "" {
			if txt, ok := t.Notes[note]; ok {
				note = txt.String()
			}
			left += " — " + note
		}
	}
	if left != "" {
		g.DrawText(dst, left, 4, -1, th.TopBarText)
	}

	if t.ScoreVar != "" {
		var sc string
		v := g.State.Var(t.ScoreVar)
		label := t.ScoreLabel
		if label.Empty() {
			label = ScoreLabel
		}
		if t.ScoreMax > 0 {
			sc = fmt.Sprintf("%s: %v/%d", label, v, t.ScoreMax)
		} else {
			sc = fmt.Sprintf("%s: %v", label, v)
		}
		x := (g.Width - TextWidth(sc)) / 2
		g.DrawText(dst, sc, x, -1, th.TopBarAccent)
	}

	if t.TimeVar != "" {
		tm := varText(g.State.Var(t.TimeVar))
		x := g.Width - TextWidth(tm) - 4
		g.DrawText(dst, tm, x, -1, th.TopBarText)
	}
}

// BlocksClickAt reports whether p lies above the bar's bottom edge, so the
// verb coin does not open over it.
func (t *TopBar) BlocksClickAt(p Point) bool {
	h := t.Height
	if h <= 0 {
		h = 12
	}
	return p.Y < float64(h)
}

// varText renders a State var for the bar. A var holding a Text resolves
// into the player's language, which is what a game whose strip says
// something in words rather than in numbers needs; anything else prints
// as it stands, and an unset var says nothing at all.
func varText(v any) string {
	switch value := v.(type) {
	case nil:
		return ""
	case Text:
		return value.String()
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}
