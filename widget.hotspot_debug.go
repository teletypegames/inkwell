package inkwell

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DebugKind is what a DebugTarget is, and decides the colour of the plate
// HotspotDebug paints over it. The kinds are the things a player can point
// at: a hotspot painted into the scene, a way out of it, a character
// standing in it, an item in the inventory, and a target mounted on a
// HotspotPanel.
type DebugKind int

const (
	DebugHotspot DebugKind = iota
	DebugExit
	DebugCharacter
	DebugItem
	DebugPanel
)

// DebugTarget is one thing on screen HotspotDebug should mark: its kind,
// the name the game knows it by, and the screen rectangle it occupies.
type DebugTarget struct {
	Kind DebugKind
	Name string
	Area Rectangle
}

// Targeted is implemented by a widget that has clickable things of its own
// to show under the debug overlay — the inventory's slots, a panel's
// hotspots. The engine's own widgets implement it, and a game widget that
// embeds one of them inherits it, so the overlay finds a game's slots
// without being told about them.
type Targeted interface {
	DebugTargets(g *Game) []DebugTarget
}

// HotspotDebug paints a translucent plate over everything on screen a
// player could click or look at, one colour per kind, with the element's
// name written across it: the current scene's hotspots and exits, the
// characters standing in it, and whatever every visible Targeted widget
// offers — the inventory's items, a panel's hotspots. F1 toggles it at
// runtime; the default is off.
type HotspotDebug struct {
	Name    string
	When    Condition
	Enabled bool
	// ToggleKey, if non-zero, overrides the default F1 toggle.
	ToggleKey ebiten.Key
	// Colors overrides the plate colour of a kind; a kind left out keeps
	// its default. Alpha is honoured, so a plate can stay translucent.
	Colors map[DebugKind]color.Color
}

func (h *HotspotDebug) GetName() string        { return h.Name }
func (h *HotspotDebug) VisibleWhen() Condition { return h.When }
func (h *HotspotDebug) Layer() Layer           { return LayerDebug }

var debugDefaultColors = map[DebugKind]color.Color{
	DebugHotspot:   color.NRGBA{R: 0xE8, G: 0xC8, B: 0x30, A: 0x70},
	DebugExit:      color.NRGBA{R: 0x30, G: 0xC0, B: 0xE8, A: 0x70},
	DebugCharacter: color.NRGBA{R: 0xE8, G: 0x40, B: 0xC8, A: 0x70},
	DebugItem:      color.NRGBA{R: 0x40, G: 0xE8, B: 0x60, A: 0x70},
	DebugPanel:     color.NRGBA{R: 0xE8, G: 0x80, B: 0x30, A: 0x70},
}

var debugLabelColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

func (h *HotspotDebug) Tick(ctx *UICtx) {
	key := h.ToggleKey
	if key == 0 {
		key = ebiten.KeyF1
	}
	if ctx.Game.Input.ClaimKey(key) {
		h.Enabled = !h.Enabled
	}
}

func (h *HotspotDebug) Draw(dst *ebiten.Image, ctx *UICtx) {
	if !h.Enabled {
		return
	}
	for _, t := range h.targets(ctx.Game) {
		h.plate(dst, ctx.Game, t)
	}
}

// targets gathers everything the overlay marks this frame: the scene's own
// hotspots (an exit is one whose generated name says so), the characters
// placed in the scene, and the targets of every visible widget that offers
// any.
func (h *HotspotDebug) targets(g *Game) []DebugTarget {
	var out []DebugTarget
	if g.currentScene != "" {
		for _, hs := range g.SceneHotspots(g.currentScene) {
			if hs.Area == nil {
				continue
			}
			kind := DebugHotspot
			if strings.HasPrefix(hs.Name, "exit:") {
				kind = DebugExit
			}
			out = append(out, DebugTarget{Kind: kind, Name: hs.Name, Area: hs.Area.Bounds()})
		}
		for _, a := range g.sceneActors(g.SceneManager.MustGet(g.currentScene)) {
			c, ok := g.runtimeChar(a.CharacterName)
			if !ok {
				continue
			}
			w, ht := c.def.W, c.def.H
			if w == 0 {
				w = 24
			}
			if ht == 0 {
				ht = 56
			}
			out = append(out, DebugTarget{
				Kind: DebugCharacter,
				Name: c.def.Name,
				Area: Rect(c.pos.X-w/2, c.pos.Y-ht, w, ht),
			})
		}
	}
	ctx := g.makeCtx()
	for _, w := range g.WidgetManager.All() {
		src, ok := w.(Targeted)
		if !ok || !WidgetVisible(ctx, w) {
			continue
		}
		out = append(out, src.DebugTargets(g)...)
	}
	return out
}

func (h *HotspotDebug) plate(dst *ebiten.Image, g *Game, t DebugTarget) {
	col, ok := h.Colors[t.Kind]
	if !ok || col == nil {
		col = debugDefaultColors[t.Kind]
	}
	b := t.Area
	vector.DrawFilledRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), col, false)
	vector.StrokeRect(dst, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, opaque(col), false)
	if width := int(b.W) - 2; width >= GlyphW {
		g.DrawText(dst, ClipText(t.Name, width), int(b.X)+1, int(b.Y)+1, debugLabelColor)
	}
}

// opaque is the plate colour with its alpha removed, for the outline that
// keeps two overlapping plates apart.
func opaque(c color.Color) color.Color {
	r, gr, b, a := c.RGBA()
	if a == 0 {
		return c
	}
	return color.NRGBA{
		R: uint8(r * 0xFFFF / a >> 8),
		G: uint8(gr * 0xFFFF / a >> 8),
		B: uint8(b * 0xFFFF / a >> 8),
		A: 0xFF,
	}
}
