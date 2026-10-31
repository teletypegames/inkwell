package chatlog

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// SpeakerStyle says how a ChatLog names who is talking.
type SpeakerStyle int

const (
	// SpeakerInline puts "DEX: " in front of every reply, which is what a
	// log of a conversation between several people needs.
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

// Host is what a [ChatLog] needs from the game it is drawn for, the UICtx's
// Game. On a Game that does not implement it the log draws nothing.
type Host interface {
	widget.Screen
	// Messages is the log, oldest first.
	Messages() []inkbase.LogMessage
	// DeliveryOf is how the voice behind a speaker's name delivers its
	// lines, which is what Deliveries filters on.
	DeliveryOf(name string) inkbase.Delivery
	// SpeechColor is the colour of a speaker's lines, or fallback when the
	// speaker has none.
	SpeechColor(name string, fallback color.Color) color.Color
	// SpeakerLabel is what the log calls a speaker.
	SpeakerLabel(name string) inkbase.Text
	// ChatLogColors is the palette the log draws with.
	ChatLogColors() Colors
}

// ChatLog renders the in-game message log (player commands + in-world
// replies) as a scrolling box. The newest message is at the bottom; older
// messages scroll up out of view as the buffer fills.
//
// A log can be given over to part of the conversation rather than all of
// it: Deliveries keeps only the replies spoken in one of the voices it lists,
// which is how a game with a radio, a tape or a narrator hands that voice
// a panel of its own without the widget having to know what the thing is.
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
	// voices. Player commands and system lines are not a character's, so
	// they are kept either way. Nil keeps everything.
	Deliveries []inkbase.Delivery
	// Resting names the character the header carries before anybody this
	// log listens to has spoken. Only SpeakerHeader reads it.
	Resting inkbase.Talker
}

// GetName returns the name the log is registered under.
func (c *ChatLog) GetName() string { return c.Name }

// VisibleWhen returns the condition that switches the log on; nil means
// always.
func (c *ChatLog) VisibleWhen() inkbase.Condition { return c.When }

// Layer puts the log on [widget.LayerHUD].
func (c *ChatLog) Layer() widget.Layer { return widget.LayerHUD }

// Tick does nothing; the log reads the game's messages when it draws.
func (c *ChatLog) Tick(ctx *widget.UICtx) {}

// hears reports whether a reply from this speaker belongs in the log.
func (c *ChatLog) hears(h Host, speaker string) bool {
	if len(c.Deliveries) == 0 {
		return true
	}
	delivery := h.DeliveryOf(speaker)
	for _, d := range c.Deliveries {
		if d == delivery {
			return true
		}
	}
	return false
}

// latest is the last speaker the log heard, for the header.
func (c *ChatLog) latest(h Host) string {
	msgs := h.Messages()
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Kind == inkbase.LogResponse && msgs[i].Speaker != "" && c.hears(h, msgs[i].Speaker) {
			return msgs[i].Speaker
		}
	}
	if c.Resting == nil {
		return ""
	}
	return c.Resting.TalkerName()
}

// Draw paints the box in Bounds, the speaker header when Speaker is
// SpeakerHeader, and as many of the newest wrapped messages as fit, each in
// its kind's colour or its speaker's. A zero-sized Bounds draws nothing.
func (c *ChatLog) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	h, ok := ctx.Game.(Host)
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
				h.SpeechColor(who, colors.Response))
		}
		ruleY := float32(top+lh) + 2
		vector.StrokeLine(dst, float32(b.X)+float32(pad), ruleY,
			float32(b.X+b.W)-float32(pad), ruleY, 1, colors.Border, false)
		top += lh + 6
	}

	// Wrap each message; collect rendered lines (color + text).
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
			col = h.SpeechColor(m.Speaker, colors.Response)
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

// BlocksClickAt reports whether p falls inside Bounds. The log takes no
// clicks itself, but the screen under it is the log's, not the scene's.
func (c *ChatLog) BlocksClickAt(ctx *widget.UICtx, p inkbase.Point) bool { return c.Bounds.Contains(p) }
