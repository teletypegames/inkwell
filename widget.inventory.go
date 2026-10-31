package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// InventoryBar shows the player's items as clickable slots. Left-click on
// a slot under the "look" verb prints the item's description; under any
// other verb, it toggles selection (the cursor "picks up" the item).
// Clicking the item already held uses it on itself: its OnUseSelf runs
// when it has one, otherwise the click puts it back. A bar with no Slots
// shows and claims nothing.
type InventoryBar struct {
	Name     string
	When     Condition
	Origin   Point
	Slots    int
	Cols     int
	SlotSize int
	Gap      int

	PanelBG bool
}

// GetName returns the name the bar is registered under.
func (b *InventoryBar) GetName() string { return b.Name }

// VisibleWhen returns the condition that switches the bar on; nil means
// always.
func (b *InventoryBar) VisibleWhen() Condition { return b.When }

// Layer puts the bar on [LayerHUD].
func (b *InventoryBar) Layer() Layer { return LayerHUD }

func (b *InventoryBar) slotRect(idx int) Rectangle {
	cols := b.Cols
	if cols <= 0 {
		cols = max(b.Slots, 1)
	}
	step := b.SlotSize + b.Gap
	col := idx % cols
	row := idx / cols
	x := b.Origin.X + float64(col*step)
	y := b.Origin.Y + float64(row*step)
	return Rect(x, y, float64(b.SlotSize), float64(b.SlotSize))
}

// DebugTargets is one target per filled slot, named after the item in it,
// so the debug overlay can mark the inventory without knowing its layout.
func (b *InventoryBar) DebugTargets(g *Game) []DebugTarget {
	items := g.Inventory.Items()
	out := make([]DebugTarget, 0, len(items))
	for i := 0; i < min(b.Slots, len(items)); i++ {
		out = append(out, DebugTarget{Kind: DebugItem, Name: items[i], Area: b.slotRect(i)})
	}
	return out
}

// Tick sets the hover label to the item under the pointer and handles a
// left click on a filled slot, consuming it: under the look verb the player
// says the item's Description; with another item held, a pairing found in
// either item's OnUseWith is queued; with the same item held, its
// OnUseSelf is queued; otherwise the item is selected, or deselected when
// it already was.
func (b *InventoryBar) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	mp := g.Input.Point()
	items := g.Inventory.Items()

	// hover label: the label of the item under the pointer
	for i := 0; i < min(b.Slots, len(items)); i++ {
		if b.slotRect(i).Contains(mp) {
			g.SetHoverLabel(g.ItemLabel(items[i]).String())
			break
		}
	}

	if !g.Input.LeftClicked() {
		return
	}
	for i := 0; i < min(b.Slots, len(items)); i++ {
		if b.slotRect(i).Contains(mp) {
			g.Input.ConsumeLeft()
			name := items[i]
			if g.SelectedVerb() == "look" && g.Player != nil {
				if it, ok := g.ItemManager.Get(name); ok && !it.Description.Empty() {
					g.queueAction(Say(g.Player, it.Description), "inv-look")
					return
				}
			}
			if b.combine(g, name) || b.useSelf(g, name) {
				return
			}
			if g.Inventory.Selected() == name {
				g.Inventory.Select("")
			} else {
				g.Inventory.Select(name)
			}
			return
		}
	}
}

// combine is an item used on an item: with one in hand, clicking another
// runs whichever of the two pairs them in its OnUseWith. It reports whether
// a pairing was found, so a click that pairs nothing still selects.
func (b *InventoryBar) combine(g *Game, clicked string) bool {
	held, ok := g.ItemManager.Get(g.Inventory.Selected())
	if !ok || held.Name == clicked {
		return false
	}
	other, ok := g.ItemManager.Get(clicked)
	if !ok {
		return false
	}
	a := held.OnUseWith[other]
	if a == nil {
		a = other.OnUseWith[held]
	}
	if a == nil {
		return false
	}
	g.LogAction("> " + Format(HintVerbOn, verbLabel(g, "use"), g.ItemLabel(held.Name), g.ItemLabel(other.Name)).String())
	g.queueAction(a, "useWith items")
	g.Inventory.Select("")
	return true
}

// useSelf is the held item clicked again: an item with an OnUseSelf is
// used on its own and put back. It reports whether it ran one.
func (b *InventoryBar) useSelf(g *Game, clicked string) bool {
	if g.Inventory.Selected() != clicked {
		return false
	}
	it, ok := g.ItemManager.Get(clicked)
	if !ok || it.OnUseSelf == nil {
		return false
	}
	g.LogAction("> " + Format(HintVerb, verbLabel(g, "use"), g.ItemLabel(it.Name)).String())
	g.queueAction(it.OnUseSelf, "useSelf")
	g.Inventory.Select("")
	return true
}

// Draw paints the PanelBG backdrop when asked, then every slot, the
// selected item's slot highlighted and each filled slot showing its item's
// sprite, or a filled square when it has none.
func (b *InventoryBar) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := gameOfUI(ctx)
	th := g.Theme()
	items := g.Inventory.Items()
	if b.Slots <= 0 {
		return
	}

	if b.PanelBG {
		last := b.slotRect(b.Slots - 1)
		first := b.slotRect(0)
		x := first.X - 2
		y := first.Y - 2
		w := (last.X + last.W) - first.X + 4
		h := (last.Y + last.H) - first.Y + 4
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), th.PanelBG, false)
	}

	for i := 0; i < b.Slots; i++ {
		r := b.slotRect(i)
		bg := th.InventorySlotBG
		if i < len(items) && items[i] == g.Inventory.Selected() {
			bg = th.InventorySlotSelectedBG
		}
		vector.DrawFilledRect(dst, float32(r.X), float32(r.Y), float32(r.W), float32(r.H), bg, false)
		if i < len(items) {
			it, ok := g.ItemManager.Get(items[i])
			if !ok {
				continue
			}
			inset := 1.0
			if it.Sprite != nil && it.Sprite.Sheet != nil && g.loaded.isPlaceholder(it.Sprite.Sheet) {
				_ = g.loaded.image(it.Sprite.Sheet)
			}
			if !drawSprite(dst, g, it.Sprite, itemFrame(g, it.Sprite), FacingRight, r.X+inset, r.Y+inset, r.W-2*inset, r.H-2*inset) {
				vector.DrawFilledRect(dst,
					float32(r.X)+2, float32(r.Y)+2,
					float32(r.W)-4, float32(r.H)-4,
					th.VerbButtonSelectedBG, false)
			}
		}
	}
}

// BlocksClickAt reports whether p falls inside the rectangle the slots
// span, empty ones included.
func (b *InventoryBar) BlocksClickAt(ctx *UICtx, p Point) bool {
	if b.Slots <= 0 {
		return false
	}
	first := b.slotRect(0)
	last := b.slotRect(b.Slots - 1)
	x := first.X
	y := first.Y
	w := last.X + last.W - first.X
	h := last.Y + last.H - first.Y
	return Rect(x, y, w, h).Contains(p)
}
