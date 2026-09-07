package inkwell

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// SpeakerStyle says how a ChatLog names who is talking.
type SpeakerStyle int

const (
	// SpeakerInline puts "DEX: " in front of every reply, which is what a
	// log of a conversation between several people needs.
	SpeakerInline SpeakerStyle = iota
	// SpeakerHeader names only the latest speaker, once, above a rule. A
	// panel given over to one voice at a time reads better this way — the
	// name is a heading rather than a word repeated down the margin.
	SpeakerHeader
	// SpeakerNone names nobody.
	SpeakerNone
)

// ChatLog renders the in-game message log (player commands + in-world
// replies) as a scrolling box. The newest message is at the bottom; older
// messages scroll up out of view as the buffer fills.
//
// A log can be given over to part of the conversation rather than all of
// it: Deliveries keeps only the replies spoken in one of the voices it lists,
// which is how a game with a radio, a tape or a narrator hands that voice
// a panel of its own without the widget having to know what the thing is.
type ChatLog struct {
	Name       string
	When       Condition
	Bounds     Rectangle
	LineHeight int
	Padding    int
	// ShowBorder draws a thin border around the box.
	ShowBorder bool

	// Speaker says how a reply names who said it.
	Speaker SpeakerStyle

	// Deliveries, when set, keeps only the replies spoken in one of these
	// voices. Player commands and system lines are not a character's, so
	// they are kept either way. Nil keeps everything.
	Deliveries []Delivery

	// Resting names the character the header carries before anybody this
	// log listens to has spoken. Only SpeakerHeader reads it.
	Resting Talker
}

func (c *ChatLog) GetName() string        { return c.Name }
func (c *ChatLog) VisibleWhen() Condition { return c.When }
func (c *ChatLog) Layer() Layer           { return LayerHUD }
func (c *ChatLog) Tick(ctx *UICtx)        {}

// hears reports whether a reply from this speaker belongs in the log.
func (c *ChatLog) hears(g *Game, speaker string) bool {
	if len(c.Deliveries) == 0 {
		return true
	}
	delivery := g.DeliveryOf(speaker)
	for _, d := range c.Deliveries {
		if d == delivery {
			return true
		}
	}
	return false
}

// speechColor is the speaker's own color where the speaker carries one,
// so a log with several voices in it stays readable as several voices.
func (c *ChatLog) speechColor(g *Game, speaker string, fallback color.Color) color.Color {
	return g.SpeechColor(speaker, fallback)
}

// latest is the last speaker the log heard, for the header.
func (c *ChatLog) latest(g *Game) string {
	msgs := g.Messages()
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Kind == LogResponse && msgs[i].Speaker != "" && c.hears(g, msgs[i].Speaker) {
			return msgs[i].Speaker
		}
	}
	return talkerName(c.Resting)
}

func (c *ChatLog) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	th := g.Theme()
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

	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), th.ChatLogBG, false)
	if c.ShowBorder {
		vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, th.CharacterPanelBorder, false)
	}

	top := int(b.Y) + pad
	if c.Speaker == SpeakerHeader {
		if who := c.latest(g); who != "" {
			g.DrawText(dst, g.SpeakerLabel(who).String(), int(b.X)+pad, top,
				c.speechColor(g, who, th.ChatLogResponse))
		}
		ruleY := float32(top+lh) + 2
		vector.StrokeLine(dst, float32(b.X)+float32(pad), ruleY,
			float32(b.X+b.W)-float32(pad), ruleY, 1, th.CharacterPanelBorder, false)
		top += lh + 6
	}

	// Wrap each message; collect rendered lines (color + text).
	type line struct {
		text string
		col  color.Color
	}
	wrapW := int(b.W) - pad*2
	var rendered []line
	for _, m := range g.Messages() {
		var col color.Color
		txt := m.Text
		switch m.Kind {
		case LogAction:
			col = th.ChatLogPrompt
		case LogResponse:
			if !c.hears(g, m.Speaker) {
				continue
			}
			col = c.speechColor(g, m.Speaker, th.ChatLogResponse)
			if c.Speaker == SpeakerInline && m.Speaker != "" {
				txt = g.SpeakerLabel(m.Speaker).String() + ": " + txt
			}
		default:
			col = th.ChatLogSystem
		}
		for _, w := range WrapText(txt, wrapW) {
			rendered = append(rendered, line{text: w, col: col})
		}
	}

	// Show the LAST maxLines wrapped lines.
	maxLines := (int(b.Y+b.H) - pad - top) / lh
	if maxLines <= 0 {
		return
	}
	if start := len(rendered) - maxLines; start > 0 {
		rendered = rendered[start:]
	}
	for i, ln := range rendered {
		g.DrawText(dst, ln.text, int(b.X)+pad, top+i*lh-1, ln.col)
	}
}

// BlocksClickAt — the log is non-interactive, but it sits over the bottom
// of the scene; treat it as opaque so RadialVerbs doesn't pop up over it.
func (c *ChatLog) BlocksClickAt(p Point) bool { return c.Bounds.Contains(p) }
