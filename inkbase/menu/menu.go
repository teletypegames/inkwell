package menu

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/font"
	"git.teletypegames.org/engines/inkwell/inkbase/input"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Texts the [Menu] prints in its own chrome. A game may replace them before
// the first frame to change the wording or add a language.
var (
	// MenuEmpty is printed on a page that has no rows.
	MenuEmpty = inkbase.Text{
		{Lang: "en", Text: "(nothing here)"},
		{Lang: "hu", Text: "(nincs itt semmi)"},
	}
	// MenuKeys is the key help in the footer of a menu that can be closed.
	MenuKeys = inkbase.Text{
		{Lang: "en", Text: "↑↓ move  ENTER open  BACKSPACE back  ESC close"},
		{Lang: "hu", Text: "↑↓ mozgás  ENTER nyit  BACKSPACE vissza  ESC zár"},
	}
	// MenuKeysSticky is the key help in the footer of a Sticky menu, which
	// has no close key.
	MenuKeysSticky = inkbase.Text{
		{Lang: "en", Text: "↑↓ move  ENTER choose  BACKSPACE back"},
		{Lang: "hu", Text: "↑↓ mozgás  ENTER választ  BACKSPACE vissza"},
	}
)

// Colors is the palette a [Menu] draws with, read from the host every frame:
// BG fills the window, Border outlines it and draws its rules, Text is a
// row's label and a page's paragraph, Hover fills the row under the cursor,
// Speaker is the title bar (arrows and breadcrumb), Accent the value column,
// and Status the key help in the footer.
type Colors struct {
	BG, Border, Text, Hover, Speaker, Accent, Status color.Color
}

type host interface {
	ScreenSize() (int, int)
	DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
	WidgetInput() *input.Input
	DoAction(inkbase.Action)
	MakeCtx() *inkbase.Ctx
	MenuColors() Colors
}

// MenuItem is one row of a menu and, when it has children, the page those
// children make up. The Menu's Root is a MenuItem whose Label titles the
// window, and every level under it is a MenuItem whose Label is that level's
// crumb in the title bar.
//
// A page is either written out (Items), worked out when the row is entered
// (Expand), or a paragraph (Body). A row with none of the three is a leaf,
// and Enter runs its Do.
type MenuItem struct {
	// Label is the row's text, and the page's crumb once it is entered.
	Label inkbase.Text
	// Body makes the page this row opens a paragraph to read, wrapped to
	// the window and scrolled with the same keys, rather than rows.
	Body inkbase.Text
	// Value is the right-hand column: what the row is worth right now. Nil
	// leaves the column empty.
	Value func(*inkbase.Ctx) inkbase.Text
	// Items is the page this row opens, as written.
	Items []MenuItem
	// Expand builds the page this row opens when the player enters it, and
	// again on every frame it is open. It wins over Items.
	Expand func(*inkbase.Ctx) []MenuItem
	// Do is the action Enter runs on a leaf.
	Do inkbase.Action
	// Enabled, when it evaluates false, dims the row and refuses Enter.
	Enabled inkbase.Condition
	// KeepOpen leaves the menu up after Do ran; by default it closes.
	KeepOpen bool
}

func (it MenuItem) page(ctx *inkbase.Ctx) []MenuItem {
	if it.Expand != nil {
		return it.Expand(ctx)
	}
	return it.Items
}

func (it MenuItem) leaf() bool  { return it.Expand == nil && len(it.Items) == 0 && len(it.Body) == 0 }
func (it MenuItem) reads() bool { return len(it.Body) > 0 }

// FixedValue returns a [MenuItem.Value] function that ignores the game and
// always returns t.
func FixedValue(t inkbase.Text) func(*inkbase.Ctx) inkbase.Text {
	return func(*inkbase.Ctx) inkbase.Text { return t }
}

// Menu is a nested list in a window: a page of rows, any of which can open a
// page of its own, as deep as the [MenuItem] tree goes. Up and down walk the
// rows, Enter or right opens one, Backspace or left goes back a level (and
// closes at the root), Escape closes, PageUp, PageDown, Home and End jump.
// The mouse does the same: hovering moves the cursor, the wheel scrolls, the
// left button opens a row, the right button goes back, and the title bar has
// a back arrow and a close box. Every key the menu uses it claims, so
// nothing underneath sees it.
//
// The host the widget is ticked and drawn for must provide the screen size,
// text drawing, input, action dispatch, a Ctx and the [Colors]; on any other
// host it does nothing.
type Menu struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the widget; nil shows it always.
	When inkbase.Condition
	// Root is the first page, and its Label is the window's title.
	Root MenuItem
	// Bounds is the window. A zero width or height takes the right half of
	// the screen.
	Bounds inkbase.Rectangle
	// ToggleKey opens and closes the menu. Zero binds no key.
	ToggleKey ebiten.Key
	// Sticky keeps the menu open: it opens itself on the first tick and
	// neither Escape, Backspace at the root nor a close box closes it.
	Sticky bool
	// LineHeight is a row's height in pixels; 0 means the glyph height + 2.
	LineHeight int
	// Padding is the inner margin in pixels; 0 means 6.
	Padding int

	open    bool
	pages   []menuPage
	frame   inkbase.Rectangle
	pointer inkbase.Point
	wheel   float64
}

type menuPage struct {
	source MenuItem
	rows   []MenuItem
	lines  []string
	cursor int
	top    int
}

func (page *menuPage) length() int {
	if page.source.reads() {
		return len(page.lines)
	}
	return len(page.rows)
}

// GetName returns the menu's Name, the key it is registered under.
func (m *Menu) GetName() string { return m.Name }

// VisibleWhen returns When, the condition that shows the menu; nil shows
// it always.
func (m *Menu) VisibleWhen() inkbase.Condition { return m.When }

// Layer puts the menu on [widget.LayerMenu], over the HUD and dialogs.
func (m *Menu) Layer() widget.Layer { return widget.LayerMenu }

// TicksDuringScript reports true: the menu answers its keys even while a
// script holds the game, so a pause menu works in a cutscene.
func (m *Menu) TicksDuringScript() bool { return true }

// BlocksClickAt reports whether p falls on the open menu's window, so a
// click there does not reach the scene. A closed menu blocks nothing.
func (m *Menu) BlocksClickAt(p inkbase.Point) bool { return m.open && m.frame.Contains(p) }

// IsOpen reports whether the menu is showing.
func (m *Menu) IsOpen() bool { return m.open }

// Open shows the menu on its Root page, dropping any pages it had open.
func (m *Menu) Open() {
	m.open = true
	m.pages = []menuPage{{source: m.Root}}
}

// Close hides the menu and forgets its open pages, so the next Open starts
// at the Root again.
func (m *Menu) Close() {
	m.open = false
	m.pages = nil
}

// Toggle closes an open menu and opens a closed one.
func (m *Menu) Toggle() {
	if m.open {
		m.Close()
		return
	}
	m.Open()
}

// Tick opens a Sticky menu or toggles on ToggleKey, then, while the menu is
// open, rebuilds the current page and handles its keys and the mouse.
func (m *Menu) Tick(ctx *widget.UICtx) {
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	in := h.WidgetInput()
	if m.Sticky {
		if !m.open {
			m.Open()
		}
	} else if m.ToggleKey != 0 && in.ClaimKey(m.ToggleKey) {
		m.Toggle()
	}
	if !m.open {
		return
	}
	w, sh := h.ScreenSize()
	m.frame = m.frameRect(w, sh)
	cctx := h.MakeCtx()
	m.refresh(cctx)
	m.keys(h, cctx)
	if m.open {
		m.mouse(h, cctx)
	}
}

// Draw paints the open menu onto dst: the window, the title bar with its
// breadcrumb, the visible rows or paragraph lines, and the key help.
func (m *Menu) Draw(dst *ebiten.Image, ctx *widget.UICtx) {
	if !m.open || len(m.pages) == 0 {
		return
	}
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	colors := h.MenuColors()
	w, sh := h.ScreenSize()
	m.frame = m.frameRect(w, sh)
	frame := m.frame
	pad := m.padding()
	lh := m.lineHeight()
	baseline := (lh - font.GlyphH) / 2

	page := m.top()
	if page.rows == nil && page.lines == nil {
		m.refresh(h.MakeCtx())
	}

	vector.DrawFilledRect(dst,
		float32(frame.X), float32(frame.Y), float32(frame.W), float32(frame.H),
		colors.BG, false)
	vector.StrokeRect(dst,
		float32(frame.X), float32(frame.Y), float32(frame.W), float32(frame.H),
		1, colors.Border, false)

	titleY := int(frame.Y) + baseline + 1
	crumbX := int(frame.X) + pad
	if len(m.pages) > 1 {
		h.DrawText(dst, "←", crumbX, titleY, colors.Speaker)
	}
	crumbX += font.GlyphW + 2
	closeX := int(frame.X+frame.W) - pad - font.GlyphW
	if !m.Sticky {
		h.DrawText(dst, "×", closeX, titleY, colors.Speaker)
	}
	h.DrawText(dst, m.breadcrumb(closeX-crumbX-font.GlyphW), crumbX, titleY, colors.Speaker)
	m.rule(dst, frame, float32(frame.Y)+float32(m.titleHeight())-0.5, colors.Border)

	cctx := h.MakeCtx()
	if page.source.reads() {
		for slot := 0; slot < m.visibleRows(frame); slot++ {
			index := page.top + slot
			if index >= len(page.lines) {
				break
			}
			rect := m.slotRect(frame, slot)
			h.DrawText(dst, page.lines[index], int(rect.X)+2, int(rect.Y)+baseline, colors.Text)
		}
	} else if len(page.rows) == 0 {
		rect := m.slotRect(frame, 0)
		h.DrawText(dst, MenuEmpty.String(), int(rect.X)+2, int(rect.Y)+baseline, dimmed(colors.Text))
	}
	for slot := 0; slot < m.visibleRows(frame); slot++ {
		index := page.top + slot
		if index >= len(page.rows) {
			break
		}
		row := page.rows[index]
		rect := m.slotRect(frame, slot)
		label := colors.Text
		if index == page.cursor {
			vector.DrawFilledRect(dst,
				float32(rect.X), float32(rect.Y), float32(rect.W), float32(rect.H),
				colors.Hover, false)
			label = colors.Hover
		}
		if row.Enabled != nil && !row.Enabled.Eval(cctx) {
			label = dimmed(label)
		}
		right := int(rect.X + rect.W)
		if !row.leaf() {
			h.DrawText(dst, "→", right-font.GlyphW, int(rect.Y)+baseline, label)
		}
		right -= font.GlyphW + 2
		if row.Value != nil {
			if value := font.ClipText(row.Value(cctx).String(), int(rect.W)/2); value != "" {
				right -= font.TextWidth(value)
				h.DrawText(dst, value, right, int(rect.Y)+baseline, colors.Accent)
				right -= font.GlyphW
			}
		}
		h.DrawText(dst, font.ClipText(row.Label.String(), right-int(rect.X)-2),
			int(rect.X)+2, int(rect.Y)+baseline, label)
	}

	footY := float32(frame.Y+frame.H) - float32(m.footerHeight())
	m.rule(dst, frame, footY+0.5, colors.Border)
	keys := MenuKeys
	if m.Sticky {
		keys = MenuKeysSticky
	}
	h.DrawText(dst, font.ClipText(keys.String(), int(frame.W)-pad*2),
		int(frame.X)+pad, int(footY)+baseline, dimmed(colors.Status))
}

func (m *Menu) rule(dst *ebiten.Image, frame inkbase.Rectangle, y float32, c color.Color) {
	vector.StrokeLine(dst, float32(frame.X), y, float32(frame.X+frame.W), y, 1, c, false)
}

func (m *Menu) keys(h host, ctx *inkbase.Ctx) {
	in := h.WidgetInput()
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
		m.enter(h, ctx)
	case in.ClaimKey(ebiten.KeyBackspace), in.ClaimKey(ebiten.KeyArrowLeft):
		m.back(ctx)
	}
}

func (m *Menu) mouse(h host, ctx *inkbase.Ctx) {
	in := h.WidgetInput()
	point := in.Point()
	inside := m.frame.Contains(point)

	if _, wheel := in.Wheel(); wheel != 0 && inside {
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

	if point != m.pointer {
		m.pointer = point
		if slot, ok := m.slotAt(point); ok {
			m.moveTo(m.top().top + slot)
		}
	}

	if !inside {
		return
	}
	if in.RightClicked() {
		in.ConsumeRight()
		m.back(ctx)
		return
	}
	if !in.LeftClicked() {
		return
	}
	in.ConsumeLeft()
	switch {
	case !m.Sticky && m.closeRect().Contains(point):
		m.Close()
	case len(m.pages) > 1 && m.backRect().Contains(point):
		m.back(ctx)
	default:
		if slot, ok := m.slotAt(point); ok {
			m.moveTo(m.top().top + slot)
			m.enter(h, ctx)
		}
	}
}

func (m *Menu) enter(h host, ctx *inkbase.Ctx) {
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
	h.DoAction(row.Do)
	if !row.KeepOpen {
		m.Close()
	}
}

func (m *Menu) back(ctx *inkbase.Ctx) {
	if len(m.pages) <= 1 {
		if !m.Sticky {
			m.Close()
		}
		return
	}
	m.pages = m.pages[:len(m.pages)-1]
	m.refresh(ctx)
}

func (m *Menu) refresh(ctx *inkbase.Ctx) {
	page := m.top()
	if page == nil {
		return
	}
	if page.source.reads() {
		page.rows = nil
		page.lines = font.WrapText(page.source.Body.String(), int(m.slotRect(m.frame, 0).W)-4)
	} else {
		page.rows = page.source.page(ctx)
	}
	m.moveTo(page.cursor)
}

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

func (m *Menu) lineHeight() int {
	if m.LineHeight > 0 {
		return m.LineHeight
	}
	return font.GlyphH + 2
}

func (m *Menu) padding() int {
	if m.Padding > 0 {
		return m.Padding
	}
	return 6
}

func (m *Menu) titleHeight() int  { return m.lineHeight() + 2 }
func (m *Menu) footerHeight() int { return m.lineHeight() }

func (m *Menu) frameRect(screenW, screenH int) inkbase.Rectangle {
	if m.Bounds.W > 0 && m.Bounds.H > 0 {
		return m.Bounds
	}
	half := float64(screenW) / 2
	return inkbase.Rect(half, 0, half, float64(screenH))
}

func (m *Menu) visibleRows(frame inkbase.Rectangle) int {
	body := int(frame.H) - m.titleHeight() - m.footerHeight()
	if body < m.lineHeight() {
		return 0
	}
	return body / m.lineHeight()
}

func (m *Menu) slotRect(frame inkbase.Rectangle, slot int) inkbase.Rectangle {
	pad := float64(m.padding())
	return inkbase.Rect(
		frame.X+pad,
		frame.Y+float64(m.titleHeight()+slot*m.lineHeight()),
		frame.W-pad*2,
		float64(m.lineHeight()),
	)
}

func (m *Menu) slotAt(point inkbase.Point) (int, bool) {
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

func (m *Menu) backRect() inkbase.Rectangle {
	return inkbase.Rect(m.frame.X, m.frame.Y, float64(m.padding()+font.GlyphW*2), float64(m.titleHeight()))
}

func (m *Menu) closeRect() inkbase.Rectangle {
	w := float64(m.padding() + font.GlyphW*2)
	return inkbase.Rect(m.frame.X+m.frame.W-w, m.frame.Y, w, float64(m.titleHeight()))
}

func (m *Menu) breadcrumb(maxPx int) string {
	labels := make([]string, len(m.pages))
	for i, page := range m.pages {
		labels[i] = page.source.Label.String()
	}
	line := strings.Join(labels, " / ")
	for len(labels) > 1 && font.TextWidth(line) > maxPx {
		labels = labels[1:]
		line = "… / " + strings.Join(labels, " / ")
	}
	return font.ClipText(line, maxPx)
}

func dimmed(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA64{
		R: uint16(r / 2),
		G: uint16(g / 2),
		B: uint16(b / 2),
		A: uint16(a),
	}
}
