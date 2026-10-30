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

// Colors is the palette a [TopBar] draws with: BG fills the strip, Text is
// the left and right sections, Accent the score.
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

// TopBar is a thin strip across the top of the screen with three sections:
// the scene title (and a note) on the left, a score in the center, and a
// time string on the right. A section with nothing to show is skipped. It
// draws nothing on a host that does not provide the screen size, text
// drawing, scene title, state variables and [Colors].
type TopBar struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Height is the strip's height in pixels; 0 means 12.
	Height int
	// LeftText replaces the scene title on the left.
	LeftText inkbase.Text
	// ScoreVar is the state variable printed as the score; empty shows no
	// score.
	ScoreVar string
	// ScoreMax, when above 0, prints the score as "Score: X/MAX".
	ScoreMax int
	// ScoreLabel replaces [ScoreLabel] as the score's label.
	ScoreLabel inkbase.Text
	// TimeVar is the state variable printed at the right edge; empty shows
	// nothing there.
	TimeVar string
	// NoteVar is the state variable appended to the left text after an em
	// dash; empty or unset adds nothing.
	NoteVar string
	// Notes maps a value of NoteVar to the text printed for it. A value with
	// no entry prints as it stands.
	Notes map[string]inkbase.Text
}

// GetName returns the top bar's Name, the key it is registered under.
func (t *TopBar) GetName() string { return t.Name }

// VisibleWhen returns When, the condition that shows the top bar; nil shows
// it always.
func (t *TopBar) VisibleWhen() inkbase.Condition { return t.When }

// Layer puts the top bar on [widget.LayerHUD].
func (t *TopBar) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing: the bar; it reads the state as it draws has no state of its own to advance.
func (t *TopBar) Tick(ctx *widget.UICtx) {}

// Draw paints the strip and its three sections, reading the state
// variables afresh each frame.
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

// BlocksClickAt reports whether p lies above the bottom edge of the strip,
// so a click on the bar does not reach the scene.
func (t *TopBar) BlocksClickAt(p inkbase.Point) bool {
	barH := t.Height
	if barH <= 0 {
		barH = 12
	}
	return p.Y < float64(barH)
}

// VarText renders a state variable for the bar: a Text resolved into the
// player's language, a string as it is, nil as "", and anything else with
// fmt.Sprint.
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
