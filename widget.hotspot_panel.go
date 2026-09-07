package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// HotspotPanel makes a region of the HUD answer the verbs, the way a
// region of the picture already does. It owns nothing but a hit-test: the
// hotspots it hands back are dispatched through Game.InvokeVerb, so a
// device mounted on the panel gets look, use, talk and take — and the
// use-with flow, and the log lines, and the nothing-happens fallback —
// out of the same code a painted hotspot does.
//
// This is what keeps a game's own HUD from re-implementing the engine.
// Without it, a widget with two clickable slots on it has to read the
// selected verb, check for a held item, consume the click and answer
// every combination itself, in Go, in a switch — which is a hundred
// lines of engine, re-decided, with the game's own writing trapped
// inside it. With it, the slots are Hotspot literals like any other.
//
// The panel draws nothing. A game that wants the slots painted embeds it
// and declares a Draw of its own:
//
//	type tapeSlots struct {
//		inkwell.HotspotPanel
//	}
//
//	func (slots *tapeSlots) Draw(dst *ebiten.Image, ctx *inkwell.UICtx) { … }
type HotspotPanel struct {
	Name string
	When Condition
	// Bounds is the region the panel claims. A click inside it is the
	// panel's whether it lands on a hotspot or not, so the scene behind
	// never receives one.
	Bounds Rectangle
	// Hotspots is asked for the panel's targets, in screen coordinates,
	// on the frame it needs them. It is a function rather than a slice
	// because what is mounted on a panel changes as the game does — a
	// slot answers for whatever is in it — and a widget that asks each
	// time needs no invalidation.
	Hotspots func(g *Game) []Hotspot
}

func (p *HotspotPanel) GetName() string        { return p.Name }
func (p *HotspotPanel) VisibleWhen() Condition { return p.When }
func (p *HotspotPanel) Layer() Layer           { return LayerHUD }

// DebugTargets is the panel's hotspots as the debug overlay sees them: one
// plate per target, named as the hotspot is.
func (p *HotspotPanel) DebugTargets(g *Game) []DebugTarget {
	if p.Hotspots == nil {
		return nil
	}
	hs := p.Hotspots(g)
	out := make([]DebugTarget, 0, len(hs))
	for _, h := range hs {
		if h.Area == nil {
			continue
		}
		out = append(out, DebugTarget{Kind: DebugPanel, Name: h.Name, Area: h.Area.Bounds()})
	}
	return out
}

func (p *HotspotPanel) Draw(dst *ebiten.Image, ctx *UICtx) {}

func (p *HotspotPanel) BlocksClickAt(pt Point) bool { return p.Bounds.Contains(pt) }

// at is the topmost hotspot containing pt, resolved the same way the
// scene's are: first area to contain the point wins.
func (p *HotspotPanel) at(g *Game, pt Point) *Hotspot {
	if p.Hotspots == nil {
		return nil
	}
	spots := p.Hotspots(g)
	for i := range spots {
		if spots[i].Area != nil && spots[i].Area.Contains(pt) {
			return &spots[i]
		}
	}
	return nil
}

func (p *HotspotPanel) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	spot := p.at(g, g.Input.Point())
	if spot == nil {
		return
	}
	g.SetHoverLabel(hotspotLabel(spot).String())
	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	g.InvokeVerb(spot, g.SelectedVerb())
}
