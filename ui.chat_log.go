package pncdsl

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ChatLog renders the in-game message log (player commands + in-world
// replies) as a scrolling box. The newest message is at the bottom; older
// messages scroll up out of view as the buffer fills.
type ChatLog struct {
	Name       string
	Bounds     Rectangle
	LineHeight int
	Padding    int
	// ShowBorder draws a thin border around the box.
	ShowBorder bool
}

func (c *ChatLog) GetName() string { return c.Name }
func (c *ChatLog) Tick(ctx *UICtx) {}

func (c *ChatLog) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
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

	msgs := g.Messages()
	maxLines := int((b.H - float64(2*pad)) / float64(lh))
	if maxLines <= 0 {
		return
	}

	// Wrap each message; collect rendered lines (color + text).
	type line struct {
		text string
		col  color.Color
	}
	wrapW := int(b.W) - pad*2
	var rendered []line
	for _, m := range msgs {
		var col color.Color
		switch m.Kind {
		case LogAction:
			col = th.ChatLogPrompt
		case LogResponse:
			col = th.ChatLogResponse
		default:
			col = th.ChatLogSystem
		}
		txt := m.Text
		if m.Kind == LogResponse && m.Speaker != "" {
			txt = m.Speaker + ": " + m.Text
		}
		for _, w := range wrapText(txt, wrapW) {
			rendered = append(rendered, line{text: w, col: col})
		}
	}

	// Show the LAST maxLines wrapped lines.
	start := len(rendered) - maxLines
	if start < 0 {
		start = 0
	}
	for i, ln := range rendered[start:] {
		y := int(b.Y) + pad + i*lh
		g.DrawText(dst, ln.text, int(b.X)+pad, y-1, ln.col)
	}
}

// BlocksClickAt — the log is non-interactive, but it sits over the bottom
// of the scene; treat it as opaque so RadialVerbs doesn't pop up over it.
func (c *ChatLog) BlocksClickAt(p Point) bool { return c.Bounds.Contains(p) }
