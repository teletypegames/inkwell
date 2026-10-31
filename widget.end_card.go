package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/endcard"

// EndCard is a fullscreen overlay shown when the ShowEnd action fires.
// While active it consumes every click so nothing underneath reacts. It
// draws in the active theme's end card colours; see [endcard.EndCard].
type EndCard = endcard.EndCard

// EndCardText is the text ShowEnd put up, or "" while the game is not over.
// It satisfies [endcard.Host].
func (g *Game) EndCardText() string { return g.endCard }

// EndCardColors is the palette an [EndCard] draws with: the active theme's
// EndCardBG and EndCardText. It satisfies [endcard.Host].
func (g *Game) EndCardColors() endcard.Colors {
	th := g.Theme()
	return endcard.Colors{BG: th.EndCardBG, Text: th.EndCardText}
}
