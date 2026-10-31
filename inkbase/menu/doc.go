// Package menu is the menu system: [Menu], a nested list in a window — a
// page of rows, any of which can open a page of its own, as deep as the tree
// goes. It is a widget like any other, on the menu layer, gated by When, and
// it is the whole of the menu system: a debug menu, a pause menu, a main menu
// and an in-fiction device screen are this one widget with a different Root.
// The root inkwell package re-exports it — its *Game is the menu's [Host] —
// and builds the debug menu and the explorer on it.
//
// # Rows and pages
//
// A row and a page are the same type, [MenuItem]. Its Label names the row on
// its parent's page and titles the page it opens, which is what makes the
// breadcrumb in the title bar fall out for free. A row opens one of three
// kinds of page:
//
//   - Items: a page as an author typed it. It costs nothing per frame.
//   - Expand: a page read out of the game. It is rebuilt while it is open,
//     so a list of what is clickable in the current scene rewrites itself
//     when the player leaves the scene.
//   - Body: a paragraph rather than rows — the Text wrapped to the window in
//     the player's language and scrolled with the same keys, with nothing to
//     pick. A codex entry, a letter, a bio.
//
// A row with none of the three is a leaf, and Enter hands its Do to the host.
// Value is the right-hand column — a count, a target, "held" — and
// [FixedValue] makes one that does not depend on the game. Enabled false
// dims the row (darkened, so it needs no colour of its own) and refuses
// Enter. KeepOpen leaves the menu up after Do. Value, Expand and Enabled take
// the same [inkbase.Ctx] an action or condition takes.
//
//	&menu.Menu{
//		Name:      "pause",
//		ToggleKey: ebiten.KeyEscape,
//		Root: menu.MenuItem{
//			Label: inkbase.Plain("PAUSED"),
//			Items: []menu.MenuItem{
//				{Label: resume, Do: closeMenu},
//				{Label: settings, Items: settingsRows},
//				{Label: quit, Do: action.Quit()},
//			},
//		},
//	}
//
// # Keys and the mouse
//
// The menu answers what a list has always answered to, and claims every key
// it uses, so nothing underneath sees it:
//
//   - up and down, or the wheel, walk the rows, wrapping;
//   - PageUp, PageDown, Home and End jump;
//   - Enter or right opens the row, or runs it;
//   - Backspace, left or the right mouse button go back a level, and close
//     at the root;
//   - Escape, or the close box in the title bar, closes;
//   - hovering moves the cursor, only when the pointer itself moved, so it
//     does not fight the keys;
//   - the left mouse button opens the row under the pointer, and the back
//     arrow in the title bar goes back a level.
//
// Only clicks inside the window are consumed, so a menu down one side of the
// screen leaves the other side playable. ToggleKey opens and closes it; zero
// binds no key and leaves that to Open, Close and Toggle, which is how a HUD
// button drives one.
//
// Sticky makes the menu the screen rather than a window over one: it is open
// whenever it is visible, Escape and back-at-the-root do nothing, ToggleKey
// is ignored and the close box is not drawn. A main menu at boot has nothing
// behind it to close onto, and a menu that can be dismissed into an empty
// frame is a way to lose a game. A sticky menu is what a Stage is usually
// made of, and its key legend ([MenuKeysSticky] rather than [MenuKeys]) offers
// no Escape. [MenuEmpty] is what a page with no rows says.
//
// A Menu is Awake: it ticks while a script holds the frame, because a menu
// that goes deaf when a cutscene starts cannot be closed. It blocks clicks
// only inside its open window.
//
// # The host
//
// The menu reaches its game through [Host], which the UICtx's Game must
// implement: ScreenSize and DrawText (widget.Screen), WidgetInput, Do,
// MakeCtx and MenuColors() [Colors]. On a Game that does not, the menu does
// nothing. inkwell's *Game is a Host, and there the colours come from the
// active theme's dialog tokens — the background and border, the speaker
// colour for the title bar, the choice colours for the cursor row, the top
// bar's accent for the value column — so a menu needs no theme work to match
// the game.
//
// [MenuEmpty], [MenuKeys] and [MenuKeysSticky] are package variables a game
// may replace before the first frame. In an inkwell game the root package's
// variables of the same names are the ones to change: the engine copies them
// here every frame.
package menu
