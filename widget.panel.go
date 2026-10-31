package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/panel"

// Panel is a generic background widget — a coloured rectangle, optionally
// outlined and opaque to clicks; see [panel.Panel]. Its colours default to
// the active theme's, through [Game.PanelColors].
type Panel = panel.Panel

// PanelColors is the fallback palette of a [Panel] that sets no colour of
// its own: the active theme's PanelBG and DialogBorder. It satisfies
// [panel.Host].
func (g *Game) PanelColors() panel.Colors {
	th := g.Theme()
	return panel.Colors{BG: th.PanelBG, Border: th.DialogBorder}
}
