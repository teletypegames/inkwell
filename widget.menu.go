package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/menu"

// MenuItem is one row of a [Menu] and, when it has children, the page those
// children make up; see [menu.MenuItem].
type MenuItem = menu.MenuItem

// Menu is a nested list in a window — the whole of the menu system: a debug
// menu, a pause menu and an in-fiction device screen are this one widget
// with a different Root. It draws with the active theme's dialog colours;
// see [menu.Menu] and [Game.MenuColors].
type Menu = menu.Menu

// FixedValue is a row's value column when it does not depend on the game: a
// constant Text in the shape MenuItem.Value asks for; see [menu.FixedValue].
func FixedValue(t Text) func(*Ctx) Text { return menu.FixedValue(t) }

// MenuColors is the palette a [Menu] draws with, from the active theme: the
// dialog background, border and text, the dialog choice colours for the
// row under the cursor, the dialog speaker colour for the title bar, the top
// bar's accent for the value column and the status text for the footer. It
// satisfies [menu.Host].
func (g *Game) MenuColors() menu.Colors {
	th := g.Theme()
	return menu.Colors{
		BG:      th.DialogBG,
		Border:  th.DialogBorder,
		Text:    th.DialogText,
		HoverBG: th.DialogChoiceBG,
		Hover:   th.DialogChoiceHover,
		Speaker: th.DialogSpeaker,
		Accent:  th.TopBarAccent,
		Status:  th.StatusText,
	}
}
