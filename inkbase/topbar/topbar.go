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

// Host is what a [TopBar] needs from the game it is drawn for, the UICtx's
// Game. On a Game that does not implement it the bar draws nothing.
type Host interface {
	widget.Screen
	// SceneTitle is the left section when LeftText is empty; an empty Text
	// leaves it blank.
	SceneTitle() inkbase.Text
	// StateVar reads a state variable, nil when it was never set.
	StateVar(name string) any
	// TopBarColors is the palette the bar draws with.
	TopBarColors() Colors
}

// TopBar is a thin strip across the top of the screen with three sections:
// scene title on the left, a score number in the center, and a day/time
// string on the right. Empty fields are skipped.
type TopBar struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Height is the strip's height in pixels; 0 means 12.
	Height int
	// LeftText overrides the auto-discovered scene Title/Name.
	LeftText inkbase.Text
	// ScoreVar is a State Var key whose value is fmt-printed as the
	// center score; empty = no score shown.
	ScoreVar string
	// ScoreMax, when above 0, prints the score as "Score: X/MAX".
	ScoreMax int
	// ScoreLabel overrides the engine's own word for it.
	ScoreLabel inkbase.Text
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
	Notes map[string]inkbase.Text
}

// GetName returns the name the bar is registered under.
func (t *TopBar) GetName() string { return t.Name }

// VisibleWhen returns the condition that switches the bar on; nil means
// always.
func (t *TopBar) VisibleWhen() inkbase.Condition { return t.When }

// Layer puts the bar on [widget.LayerHUD].
func (t *TopBar) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing; the bar reads the game when it draws.
func (t *TopBar) Tick(ctx *widget.UICtx) {}

// Draw paints the strip, Height pixels tall (12 when zero), with the left
// text and its note, the score in the centre and the time at the right
// edge, each read from State as the frame is drawn.
func (t *TopBar) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
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

// BlocksClickAt reports whether p lies above the bar's bottom edge.
func (t *TopBar) BlocksClickAt(ctx *widget.UICtx, p inkbase.Point) bool {
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
