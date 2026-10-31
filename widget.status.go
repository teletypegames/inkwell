package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/statusline"

// StatusLine renders a single line: either the active flash message (set
// by [Game.FlashLine]) or the hover hint (verb + target under the cursor),
// in the active theme's FlashText and StatusText colours; see
// [statusline.StatusLine].
type StatusLine = statusline.StatusLine

// StatusLineColors is the palette a [StatusLine] draws with: the active
// theme's FlashText for a flash and StatusText for the hover hint. It
// satisfies [statusline.Host].
func (g *Game) StatusLineColors() statusline.Colors {
	th := g.Theme()
	return statusline.Colors{Flash: th.FlashText, Status: th.StatusText}
}

// FlashText is the flash message up now, or "" when there is none; see
// [Game.FlashLine]. It satisfies [statusline.Host].
func (g *Game) FlashText() string { return g.flash }

// tickFlash counts the flash message's timer down by dt and clears the
// message when it runs out. The game loop calls it every frame, so a flash
// expires whether or not a StatusLine is showing it.
func (g *Game) tickFlash(dt float64) {
	if g.flashTimer <= 0 {
		return
	}
	g.flashTimer -= dt
	if g.flashTimer <= 0 {
		g.flash = ""
		g.flashTimer = 0
	}
}
