package chatlog

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// SpeakerStyle says how a [ChatLog] names who is talking.
type SpeakerStyle int

const (
	// SpeakerInline puts the speaker's label and a colon in front of every
	// reply.
	SpeakerInline SpeakerStyle = iota
	// SpeakerHeader names only the latest speaker, once, above a rule at the
	// top of the box.
	SpeakerHeader
	// SpeakerNone names nobody.
	SpeakerNone
)

// Colors is the palette a [ChatLog] draws with: BG fills the box, Border
// outlines it and draws the header rule, Prompt colors the player's commands,
// Response a reply whose speaker has no color of their own, and System every
// other line.
type Colors struct {
	BG       color.Color
	Border   color.Color
	Prompt   color.Color
	Response color.Color
	System   color.Color
}

type host interface {
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	Messages() []inkbase.LogMessage
	SpeakerDelivery(name string) inkbase.Delivery
	SpeakerColor(name string, fallback color.Color) color.Color
	SpeakerLabel(name string) inkbase.Text
	ChatLogColors() Colors
}

// ChatLog renders the game's message log (player commands, in-world replies
// and system lines) as a box with the newest line at the bottom; older lines
// scroll up out of view. It draws nothing on a host that does not provide
// the messages, speaker lookups and [Colors].
type ChatLog struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Bounds is the box; a zero width or height draws nothing.
	Bounds inkbase.Rectangle
	// LineHeight is a line's height in pixels; 0 means 10.
	LineHeight int
	// Padding is the inner margin in pixels; 0 means 3.
	Padding int
	// ShowBorder draws a thin border around the box.
	ShowBorder bool
	// Speaker says how a reply names who said it.
	Speaker SpeakerStyle
	// Deliveries, when set, keeps only the replies spoken in one of these
	// voices; commands and system lines are kept either way. Nil keeps
	// everything.
	Deliveries []inkbase.Delivery
	// Resting is the speaker name the header carries before anybody this log
	// listens to has spoken. Only SpeakerHeader reads it.
	Resting string
}

// GetName returns the chat log's Name, the key it is registered under.
func (c *ChatLog) GetName() string { return c.Name }

// VisibleWhen returns When, the condition that shows the chat log; nil shows
// it always.
func (c *ChatLog) VisibleWhen() inkbase.Condition { return c.When }

// Layer puts the chat log on [widget.LayerHUD].
func (c *ChatLog) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing: the log; the messages live in the host has no state of its own to advance.
func (c *ChatLog) Tick(ctx *widget.UICtx) {}

func (c *ChatLog) hears(h host, speaker string) bool {
	if len(c.Deliveries) == 0 {
		return true
	}
	delivery := h.SpeakerDelivery(speaker)
	for _, d := range c.Deliveries {
		if d == delivery {
			return true
		}
	}
	return false
}

func (c *ChatLog) latest(h host) string {
	msgs := h.Messages()
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Kind == inkbase.LogResponse && msgs[i].Speaker != "" && c.hears(h, msgs[i].Speaker) {
			return msgs[i].Speaker
		}
	}
	return c.Resting
}

// Draw paints the box, the speaker header when Speaker is SpeakerHeader, and
// as many of the latest wrapped lines as fit, each in its kind's color.
func (c *ChatLog) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	colors := h.ChatLogColors()
	b := c.Bounds
	if b.W <= 0 || b.H <= 0 {
		return
	}
	lh := c.LineHeight
	if lh <= 0 {
		lh = 10
	}
	pad := c.Padding
	if pad <= 0 {
		pad = 3
	}

	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), colors.BG, false)
	if c.ShowBorder {
		vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, colors.Border, false)
	}

	top := int(b.Y) + pad
	if c.Speaker == SpeakerHeader {
		if who := c.latest(h); who != "" {
			h.DrawText(dst, h.SpeakerLabel(who).String(), int(b.X)+pad, top,
				h.SpeakerColor(who, colors.Response))
		}
		ruleY := float32(top+lh) + 2
		vector.StrokeLine(dst, float32(b.X)+float32(pad), ruleY,
			float32(b.X+b.W)-float32(pad), ruleY, 1, colors.Border, false)
		top += lh + 6
	}

	type line struct {
		text string
		col  color.Color
	}
	wrapW := int(b.W) - pad*2
	var rendered []line
	for _, m := range h.Messages() {
		var col color.Color
		txt := m.Text
		switch m.Kind {
		case inkbase.LogAction:
			col = colors.Prompt
		case inkbase.LogResponse:
			if !c.hears(h, m.Speaker) {
				continue
			}
			col = h.SpeakerColor(m.Speaker, colors.Response)
			if c.Speaker == SpeakerInline && m.Speaker != "" {
				txt = h.SpeakerLabel(m.Speaker).String() + ": " + txt
			}
		default:
			col = colors.System
		}
		for _, w := range font.WrapText(txt, wrapW) {
			rendered = append(rendered, line{text: w, col: col})
		}
	}

	maxLines := (int(b.Y+b.H) - pad - top) / lh
	if maxLines <= 0 {
		return
	}
	if start := len(rendered) - maxLines; start > 0 {
		rendered = rendered[start:]
	}
	for i, ln := range rendered {
		h.DrawText(dst, ln.text, int(b.X)+pad, top+i*lh-1, ln.col)
	}
}

// BlocksClickAt reports whether p falls inside Bounds, so a click on the log
// does not reach the scene.
func (c *ChatLog) BlocksClickAt(p inkbase.Point) bool { return c.Bounds.Contains(p) }
