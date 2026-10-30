package inkwell

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MenuItem is one row of a menu and — when it has children — the page those
// children make up. A menu is therefore a tree of a single type: the Menu's
// Root is a MenuItem whose Label titles the window, and every level under it
// is a MenuItem whose Label is that level's crumb in the title bar.
//
// A page is either written out (Items), worked out when the row is entered
// (Expand), or a paragraph (Body). A row with none of the three is a leaf,
// and Enter queues its Do.
type MenuItem struct {
	Label Text

	// Body makes the page this row opens a paragraph rather than rows: the
	// Text wrapped to the window, scrolled with the same keys. A codex
	// entry, a letter, a bio — anything the player reads rather than picks
	// from.
	Body Text

	// Value is the right-hand column: what the row is worth right now — a
	// count, a destination, "held". Nil leaves the column empty.
	Value func(*Ctx) Text

	// Items is the page this row opens, as somebody wrote it.
	Items []MenuItem

	// Expand is the page this row opens, built when the player enters it and
	// rebuilt while it is open. A list read out of a manager or out of the
	// scene the player is standing in is an Expand; anything an author typed
	// is Items.
	Expand func(*Ctx) []MenuItem

	// Do is what Enter queues on a leaf.
	Do Action

	// Enabled, when it evaluates false, dims the row and refuses Enter.
	Enabled Condition

	// KeepOpen leaves the menu up after Do ran. The default closes it, which
	// is what a row that sends the player somewhere wants; a row that toggles
	// something wants the opposite.
	KeepOpen bool
}

// page is the rows this item opens: the ones somebody wrote, or the ones
// Expand reads out of the game as it stands.
func (it MenuItem) page(ctx *Ctx) []MenuItem {
	if it.Expand != nil {
		return it.Expand(ctx)
	}
	return it.Items
}

// leaf reports whether entering this item runs it rather than opening a page.
func (it MenuItem) leaf() bool { return it.Expand == nil && len(it.Items) == 0 && len(it.Body) == 0 }

// reads reports whether the page this item opens is a Body to read rather
// than rows to pick from.
func (it MenuItem) reads() bool { return len(it.Body) > 0 }

// FixedValue is a row's value column when it does not depend on the game: a
// constant Text in the shape MenuItem.Value asks for.
func FixedValue(t Text) func(*Ctx) Text {
	return func(*Ctx) Text { return t }
}

// Menu is a nested list in a window: a page of rows, any of which can open a
// page of its own, as deep as the tree goes. It is the whole of the menu
// system — a debug menu, a pause menu and an in-fiction device screen are
// this one widget with a different Root.
//
// The keys are the ones a list has always answered to: up and down walk the
// rows, Enter (or right) opens one, Backspace (or left) goes back a level and
// closes at the root, Escape closes, PageUp/PageDown/Home/End jump. The mouse
// does the same: hovering moves the cursor, the wheel scrolls, the left
// button opens a row, the right button goes back, and the title bar carries a
// back arrow and a close box. Every key the menu uses it claims, so nothing
// underneath sees it.
type Menu struct {
	Name string
	When Condition

	// Root is the first page, and its Label is the window's title.
	Root MenuItem

	// Bounds is the window. The zero value takes the right half of the screen.
	Bounds Rectangle

	// ToggleKey opens and closes the menu. Zero binds no key at all, which
	// leaves Toggle to whatever the game wires it to.
	ToggleKey ebiten.Key

	// Sticky makes the menu the screen rather than a window over one: it is
	// open whenever it is visible, Escape and back-at-the-root do nothing,
	// and ToggleKey is ignored. A main menu at boot has nothing behind it to
	// close onto, and a menu that can be dismissed into an empty frame is a
	// way to lose a game rather than to leave one.
	Sticky bool

	// LineHeight is a row's height, Padding the window's inset. Zero gives
	// GlyphH + 2 and 6.
	LineHeight int
	Padding    int

	open    bool
	pages   []menuPage
	frame   Rectangle
	pointer Point
	wheel   float64
}

// menuPage is one open level: where it came from, the rows it is showing —
// or the lines, when the source is a Body — which row the cursor is on and
// which row or line it is scrolled to.
type menuPage struct {
	source MenuItem
	rows   []MenuItem
	lines  []string
	cursor int
	top    int
}

// length is how many rows or lines the page has to move over.
func (page *menuPage) length() int {
	if page.source.reads() {
		return len(page.lines)
	}
	return len(page.rows)
}

// GetName returns the name the menu is registered under.
func (m *Menu) GetName() string { return m.Name }

// VisibleWhen returns the condition that switches the menu on; nil means
// always.
func (m *Menu) VisibleWhen() Condition { return m.When }

// Layer puts the menu on [LayerMenu], above the game it opens over.
func (m *Menu) Layer() Layer { return LayerMenu }

// TicksDuringScript reports true: a menu has to open, answer its keys and
// close while a cutscene holds the frame.
func (m *Menu) TicksDuringScript() bool { return true }

// BlocksClickAt reports whether p falls inside the open window, so the verb
// coin does not open over it. A closed menu blocks nothing.
func (m *Menu) BlocksClickAt(p Point) bool { return m.open && m.frame.Contains(p) }

// IsOpen reports whether the window is up.
func (m *Menu) IsOpen() bool { return m.open }

// Open puts the window up on its root page. Close takes it down and forgets
// where the player was, so a menu always opens where it opened last time.
func (m *Menu) Open() {
	m.open = true
	m.pages = []menuPage{{source: m.Root}}
}

// Close takes the window down and drops every open page.
func (m *Menu) Close() {
	m.open = false
	m.pages = nil
}

// Toggle closes the window when it is open and opens it on its root page
// when it is not.
func (m *Menu) Toggle() {
	if m.open {
		m.Close()
		return
	}
	m.Open()
}

// Tick opens a Sticky menu, answers ToggleKey, and while the window is up
// rebuilds the open page and handles its keys and the mouse, claiming what
// it uses. A row entered by Enter or a click queues its Do on the game.
func (m *Menu) Tick(ctx *UICtx) {
	g := gameOfUI(ctx)
	if m.Sticky {
		if !m.open {
			m.Open()
		}
	} else if m.ToggleKey != 0 && g.Input.ClaimKey(m.ToggleKey) {
		m.Toggle()
	}
	if !m.open {
		return
	}
	m.frame = m.frameRect(g)
	cctx := g.makeCtx()
	m.refresh(cctx)
	m.keys(cctx)
	if m.open {
		m.mouse(cctx)
	}
}

// Draw paints the open window: the title bar with its breadcrumb, back
// arrow and close box, the visible rows or Body lines with the cursor row
// highlighted, and the footer listing the keys. A closed menu draws
// nothing.
func (m *Menu) Draw(dst *ebiten.Image, ctx *UICtx) {
	if !m.open || len(m.pages) == 0 {
		return
	}
	g := gameOfUI(ctx)
	th := g.Theme()
	m.frame = m.frameRect(g)
	frame := m.frame
	pad := m.padding()
	lh := m.lineHeight()
	baseline := (lh - GlyphH) / 2

	page := m.top()
	if page.rows == nil && page.lines == nil {
		m.refresh(g.makeCtx())
	}

	vector.DrawFilledRect(dst,
		float32(frame.X), float32(frame.Y), float32(frame.W), float32(frame.H),
		th.DialogBG, false)
	vector.StrokeRect(dst,
		float32(frame.X), float32(frame.Y), float32(frame.W), float32(frame.H),
		1, th.DialogBorder, false)

	titleY := int(frame.Y) + baseline + 1
	crumbX := int(frame.X) + pad
	if len(m.pages) > 1 {
		g.DrawText(dst, "←", crumbX, titleY, th.DialogSpeaker)
	}
	crumbX += GlyphW + 2
	closeX := int(frame.X+frame.W) - pad - GlyphW
	if !m.Sticky {
		g.DrawText(dst, "×", closeX, titleY, th.DialogSpeaker)
	}
	g.DrawText(dst, m.breadcrumb(closeX-crumbX-GlyphW), crumbX, titleY, th.DialogSpeaker)
	m.rule(dst, frame, float32(frame.Y)+float32(m.titleHeight())-0.5, th.DialogBorder)

	cctx := g.makeCtx()
	if page.source.reads() {
		for slot := 0; slot < m.visibleRows(frame); slot++ {
			index := page.top + slot
			if index >= len(page.lines) {
				break
			}
			rect := m.slotRect(frame, slot)
			g.DrawText(dst, page.lines[index], int(rect.X)+2, int(rect.Y)+baseline, th.DialogText)
		}
	} else if len(page.rows) == 0 {
		rect := m.slotRect(frame, 0)
		g.DrawText(dst, MenuEmpty.String(), int(rect.X)+2, int(rect.Y)+baseline, dimmed(th.DialogText))
	}
	for slot := 0; slot < m.visibleRows(frame); slot++ {
		index := page.top + slot
		if index >= len(page.rows) {
			break
		}
		row := page.rows[index]
		rect := m.slotRect(frame, slot)
		label := th.DialogText
		if index == page.cursor {
			vector.DrawFilledRect(dst,
				float32(rect.X), float32(rect.Y), float32(rect.W), float32(rect.H),
				th.DialogChoiceBG, false)
			label = th.DialogChoiceHover
		}
		if row.Enabled != nil && !row.Enabled.Eval(cctx) {
			label = dimmed(label)
		}
		right := int(rect.X + rect.W)
		if !row.leaf() {
			g.DrawText(dst, "→", right-GlyphW, int(rect.Y)+baseline, label)
		}
		right -= GlyphW + 2
		if row.Value != nil {
			if value := ClipText(row.Value(cctx).String(), int(rect.W)/2); value != "" {
				right -= TextWidth(value)
				g.DrawText(dst, value, right, int(rect.Y)+baseline, th.TopBarAccent)
				right -= GlyphW
			}
		}
		g.DrawText(dst, ClipText(row.Label.String(), right-int(rect.X)-2),
			int(rect.X)+2, int(rect.Y)+baseline, label)
	}

	footY := float32(frame.Y+frame.H) - float32(m.footerHeight())
	m.rule(dst, frame, footY+0.5, th.DialogBorder)
	keys := MenuKeys
	if m.Sticky {
		keys = MenuKeysSticky
	}
	g.DrawText(dst, ClipText(keys.String(), int(frame.W)-pad*2),
		int(frame.X)+pad, int(footY)+baseline, dimmed(th.StatusText))
}

// rule draws a one-pixel horizontal line. Its y is a pixel centre, so a
// rule meant to sit on the last row of the title band, just above the first
// row's highlight, is drawn at that row's centre rather than its edge.
func (m *Menu) rule(dst *ebiten.Image, frame Rectangle, y float32, c color.Color) {
	vector.StrokeLine(dst, float32(frame.X), y, float32(frame.X+frame.W), y, 1, c, false)
}

// ----- navigation -------------------------------------------------------

func (m *Menu) keys(ctx *Ctx) {
	game := gameOf(ctx)
	in := game.Input
	switch {
	case in.ClaimKey(ebiten.KeyEscape):
		if !m.Sticky {
			m.Close()
		}
	case in.ClaimKey(ebiten.KeyArrowUp):
		m.move(-1)
	case in.ClaimKey(ebiten.KeyArrowDown):
		m.move(1)
	case in.ClaimKey(ebiten.KeyPageUp):
		m.jump(-m.visibleRows(m.frame))
	case in.ClaimKey(ebiten.KeyPageDown):
		m.jump(m.visibleRows(m.frame))
	case in.ClaimKey(ebiten.KeyHome):
		m.moveTo(0)
	case in.ClaimKey(ebiten.KeyEnd):
		m.moveTo(len(m.top().rows) - 1)
	case in.ClaimKey(ebiten.KeyEnter), in.ClaimKey(ebiten.KeyNumpadEnter), in.ClaimKey(ebiten.KeyArrowRight):
		m.enter(ctx)
	case in.ClaimKey(ebiten.KeyBackspace), in.ClaimKey(ebiten.KeyArrowLeft):
		m.back(ctx)
	}
}

func (m *Menu) mouse(ctx *Ctx) {
	g := gameOf(ctx)
	point := g.Input.Point()
	inside := m.frame.Contains(point)

	if _, wheel := g.Input.Wheel(); wheel != 0 && inside {
		m.wheel += wheel
		for m.wheel >= 1 {
			m.wheel--
			m.move(-1)
		}
		for m.wheel <= -1 {
			m.wheel++
			m.move(1)
		}
	}

	// Hover only moves the cursor when the pointer itself moved, so a
	// pointer resting over the list does not fight the arrow keys.
	if point != m.pointer {
		m.pointer = point
		if slot, ok := m.slotAt(point); ok {
			m.moveTo(m.top().top + slot)
		}
	}

	if !inside {
		return
	}
	if g.Input.RightClicked() {
		g.Input.ConsumeRight()
		m.back(ctx)
		return
	}
	if !g.Input.LeftClicked() {
		return
	}
	g.Input.ConsumeLeft()
	switch {
	case !m.Sticky && m.closeRect().Contains(point):
		m.Close()
	case len(m.pages) > 1 && m.backRect().Contains(point):
		m.back(ctx)
	default:
		if slot, ok := m.slotAt(point); ok {
			m.moveTo(m.top().top + slot)
			m.enter(ctx)
		}
	}
}

func (m *Menu) enter(ctx *Ctx) {
	page := m.top()
	if page == nil || page.source.reads() || page.cursor >= len(page.rows) {
		return
	}
	row := page.rows[page.cursor]
	if row.Enabled != nil && !row.Enabled.Eval(ctx) {
		return
	}
	if !row.leaf() {
		m.pages = append(m.pages, menuPage{source: row})
		m.refresh(ctx)
		return
	}
	if row.Do == nil {
		return
	}
	gameOf(ctx).Do(row.Do)
	if !row.KeepOpen {
		m.Close()
	}
}

func (m *Menu) back(ctx *Ctx) {
	if len(m.pages) <= 1 {
		if !m.Sticky {
			m.Close()
		}
		return
	}
	m.pages = m.pages[:len(m.pages)-1]
	m.refresh(ctx)
}

// refresh rebuilds the open page's rows — an Expand page answers for the
// game as it is now, so teleporting out of a scene rewrites the list of what
// is clickable in it — and keeps the cursor on a row that exists. A Body
// page is wrapped to the window instead, in the player's language.
func (m *Menu) refresh(ctx *Ctx) {
	page := m.top()
	if page == nil {
		return
	}
	if page.source.reads() {
		page.rows = nil
		page.lines = WrapText(page.source.Body.String(), int(m.slotRect(m.frame, 0).W)-4)
	} else {
		page.rows = page.source.page(ctx)
	}
	m.moveTo(page.cursor)
}

// move steps the cursor, wrapping round the rows of a list; a Body has no
// cursor to wrap, so it scrolls and stops at either end.
func (m *Menu) move(step int) {
	page := m.top()
	if page == nil || page.length() == 0 {
		return
	}
	if page.source.reads() {
		m.jump(step)
		return
	}
	n := page.length()
	m.moveTo(((page.cursor+step)%n + n) % n)
}

func (m *Menu) jump(step int) {
	page := m.top()
	if page == nil {
		return
	}
	m.moveTo(page.cursor + step)
}

// moveTo puts the cursor on a row, clamped to the page, and scrolls the view
// until that row is in it.
func (m *Menu) moveTo(index int) {
	page := m.top()
	if page == nil {
		return
	}
	if index >= page.length() {
		index = page.length() - 1
	}
	if index < 0 {
		index = 0
	}
	page.cursor = index
	rows := m.visibleRows(m.frame)
	if page.source.reads() {
		page.top = page.cursor
	}
	if page.cursor < page.top {
		page.top = page.cursor
	}
	if rows > 0 && page.cursor >= page.top+rows {
		page.top = page.cursor - rows + 1
	}
	if page.top > page.length()-rows {
		page.top = page.length() - rows
	}
	if page.top < 0 {
		page.top = 0
	}
	if page.source.reads() {
		page.cursor = page.top
	}
}

func (m *Menu) top() *menuPage {
	if len(m.pages) == 0 {
		return nil
	}
	return &m.pages[len(m.pages)-1]
}

// ----- layout -----------------------------------------------------------

func (m *Menu) lineHeight() int {
	if m.LineHeight > 0 {
		return m.LineHeight
	}
	return GlyphH + 2
}

func (m *Menu) padding() int {
	if m.Padding > 0 {
		return m.Padding
	}
	return 6
}

func (m *Menu) titleHeight() int  { return m.lineHeight() + 2 }
func (m *Menu) footerHeight() int { return m.lineHeight() }

// frameRect is the window: what Bounds says, or the right half of the screen
// when nobody said.
func (m *Menu) frameRect(g *Game) Rectangle {
	if m.Bounds.W > 0 && m.Bounds.H > 0 {
		return m.Bounds
	}
	half := float64(g.Width) / 2
	return Rect(half, 0, half, float64(g.Height))
}

// visibleRows is how many rows fit between the title bar and the footer.
func (m *Menu) visibleRows(frame Rectangle) int {
	body := int(frame.H) - m.titleHeight() - m.footerHeight()
	if body < m.lineHeight() {
		return 0
	}
	return body / m.lineHeight()
}

// slotRect is the rectangle of the slot'th row on screen, which is the row
// page.top + slot of the open page.
func (m *Menu) slotRect(frame Rectangle, slot int) Rectangle {
	pad := float64(m.padding())
	return Rect(
		frame.X+pad,
		frame.Y+float64(m.titleHeight()+slot*m.lineHeight()),
		frame.W-pad*2,
		float64(m.lineHeight()),
	)
}

// slotAt is the row slot a point falls on, if it falls on one at all. A
// Body has no rows to fall on.
func (m *Menu) slotAt(point Point) (int, bool) {
	page := m.top()
	if page == nil || page.source.reads() {
		return 0, false
	}
	for slot := 0; slot < m.visibleRows(m.frame); slot++ {
		if page.top+slot >= len(page.rows) {
			break
		}
		if m.slotRect(m.frame, slot).Contains(point) {
			return slot, true
		}
	}
	return 0, false
}

func (m *Menu) backRect() Rectangle {
	return Rect(m.frame.X, m.frame.Y, float64(m.padding()+GlyphW*2), float64(m.titleHeight()))
}

func (m *Menu) closeRect() Rectangle {
	w := float64(m.padding() + GlyphW*2)
	return Rect(m.frame.X+m.frame.W-w, m.frame.Y, w, float64(m.titleHeight()))
}

// breadcrumb is the path across the open pages, cut from the front when the
// whole of it does not fit.
func (m *Menu) breadcrumb(maxPx int) string {
	labels := make([]string, len(m.pages))
	for i, page := range m.pages {
		labels[i] = page.source.Label.String()
	}
	line := strings.Join(labels, " / ")
	for len(labels) > 1 && TextWidth(line) > maxPx {
		labels = labels[1:]
		line = "… / " + strings.Join(labels, " / ")
	}
	return ClipText(line, maxPx)
}

// dimmed darkens a theme colour rather than fading it, so a disabled row
// reads as switched off in any theme without the theme having a token for
// it. RGBA() is alpha-premultiplied, and halving the channels keeps it so.
func dimmed(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA64{
		R: uint16(r / 2),
		G: uint16(g / 2),
		B: uint16(b / 2),
		A: uint16(a),
	}
}
