package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/textcard"

// TextCard is a block of lines centred on the screen — a credits roll that
// does not roll, a title card, an apology — in the active theme's end card
// colours; see [textcard.TextCard].
type TextCard = textcard.TextCard

// TextCardColors is the palette a [TextCard] draws with: the active theme's
// EndCardBG for the backdrop and EndCardText for the lines. It satisfies
// [textcard.Host].
func (g *Game) TextCardColors() textcard.Colors {
	th := g.Theme()
	return textcard.Colors{BG: th.EndCardBG, Text: th.EndCardText}
}
