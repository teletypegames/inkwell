package topbar

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

type Colors struct {
	BG     color.Color
	Text   color.Color
	Accent color.Color
}

// ScoreLabel is the default label for the score.
var ScoreLabel = inkbase.Text{
	{Lang: "en", Text: "Score"},
	{Lang: "hu", Text: "Pont"},
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	SceneTitle() inkbase.Text
	StateVar(string) any
	TopBarColors() Colors
}

type TopBar struct {
	Name       string
	When       inkbase.Condition
	Height     int
	LeftText   inkbase.Text
	ScoreVar   string
	ScoreMax   int
	ScoreLabel inkbase.Text
	TimeVar    string
	NoteVar    string
	Notes      map[string]inkbase.Text
}

func (t *TopBar) GetName() string                { return t.Name }
func (t *TopBar) VisibleWhen() inkbase.Condition { return t.When }
func (t *TopBar) Layer() widget.Layer            { return widget.LayerHUD }
func (t *TopBar) Tick(ctx *widget.UICtx)         {}

func (t *TopBar) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	colors := h.TopBarColors()
	w, _ := h.ScreenSize()
	barH := t.Height
	if barH <= 0 {
		barH = 12
	}
	vector.DrawFilledRect(dst, 0, 0, float32(w), float32(barH), colors.BG, false)

	left := t.LeftText.String()
	if left == "" {
		if title := h.SceneTitle(); !title.Empty() {
			left = title.String()
		}
	}
	if t.NoteVar != "" {
		if note := VarText(h.StateVar(t.NoteVar)); note != "" {
			if txt, ok := t.Notes[note]; ok {
				note = txt.String()
			}
			left += " — " + note
		}
	}
	if left != "" {
		h.DrawText(dst, left, 4, -1, colors.Text)
	}

	if t.ScoreVar != "" {
		var sc string
		v := h.StateVar(t.ScoreVar)
		label := t.ScoreLabel
		if label.Empty() {
			label = ScoreLabel
		}
		if t.ScoreMax > 0 {
			sc = fmt.Sprintf("%s: %v/%d", label, v, t.ScoreMax)
		} else {
			sc = fmt.Sprintf("%s: %v", label, v)
		}
		x := (w - font.TextWidth(sc)) / 2
		h.DrawText(dst, sc, x, -1, colors.Accent)
	}

	if t.TimeVar != "" {
		tm := VarText(h.StateVar(t.TimeVar))
		x := w - font.TextWidth(tm) - 4
		h.DrawText(dst, tm, x, -1, colors.Text)
	}
}

func (t *TopBar) BlocksClickAt(p inkbase.Point) bool {
	barH := t.Height
	if barH <= 0 {
		barH = 12
	}
	return p.Y < float64(barH)
}

// VarText renders a State var for the bar.
func VarText(v any) string {
	switch value := v.(type) {
	case nil:
		return ""
	case inkbase.Text:
		return value.String()
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}
