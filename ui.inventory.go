package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// InventoryBar shows the player's items as clickable slots. Left-click on
// a slot under the "look" verb prints the item's description; under any
// other verb, it toggles selection (the cursor "picks up" the item).
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

func (b *InventoryBar) GetName() string        { return b.Name }
func (b *InventoryBar) VisibleWhen() Condition { return b.When }
func (b *InventoryBar) Layer() Layer           { return LayerHUD }

func (b *InventoryBar) slotRect(idx int) Rectangle {
	cols := b.Cols
	if cols <= 0 {
		cols = b.Slots
	}
	step := b.SlotSize + b.Gap
	col := idx % cols
	row := idx / cols
	x := b.Origin.X + float64(col*step)
	y := b.Origin.Y + float64(row*step)
	return Rect(x, y, float64(b.SlotSize), float64(b.SlotSize))
}

func (b *InventoryBar) Tick(ctx *UICtx) {
	g := ctx.Game
	mp := g.Input.Point()
	items := g.Inventory.Items()

	// hover label: when over a slot with an item, show its description
	for i := 0; i < min(b.Slots, len(items)); i++ {
		if b.slotRect(i).Contains(mp) {
			if it, ok := g.ItemManager.Get(items[i]); ok && it.Description != "" {
				g.SetHoverLabel(it.Description)
			} else {
				g.SetHoverLabel(items[i])
			}
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
			if g.SelectedVerb() == "look" {
				if it, ok := g.ItemManager.Get(name); ok && it.Description != "" {
					g.queueAction(Say("player", it.Description), "inv-look")
					return
				}
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

func (b *InventoryBar) Draw(dst *ebiten.Image, ctx *UICtx) {
	g := ctx.Game
	th := g.Theme()
	items := g.Inventory.Items()

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
			img := g.loaded.image(g.AssetManager, it.Sprite)
			sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
			if sw == 0 || sh == 0 {
				continue
			}
			if g.loaded.isPlaceholder(it.Sprite) {
				// draw a small colored square as a stand-in
				inset := float32(2)
				vector.DrawFilledRect(dst,
					float32(r.X)+inset, float32(r.Y)+inset,
					float32(r.W)-2*inset, float32(r.H)-2*inset,
					th.VerbButtonSelectedBG, false)
				continue
			}
			op := &ebiten.DrawImageOptions{}
			inset := 1.0
			tgtW := r.W - 2*inset
			tgtH := r.H - 2*inset
			op.GeoM.Scale(tgtW/float64(sw), tgtH/float64(sh))
			op.GeoM.Translate(r.X+inset, r.Y+inset)
			dst.DrawImage(img, op)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BlocksClickAt — clickBlocker contract for the RadialVerbs widget.
func (b *InventoryBar) BlocksClickAt(p Point) bool {
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
