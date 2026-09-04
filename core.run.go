package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Run validates the game, then enters the ebiten main loop. The window is
// sized to Game.WindowScale times the internal resolution, 4× by default.
func Run(g *Game) error {
	// The domain may have swapped a manager in since NewGame, so the parts
	// that hold a registry directly are wired here, not at construction.
	g.resolveDerived()

	if err := g.loadSettings(); err != nil {
		logf("settings: %v", err)
	}
	if Lang == "" {
		Lang = Fallback
	}
	if err := g.Validate(); err != nil {
		return err
	}
	// If the domain didn't register any widgets, fall back to the SCUMM
	// preset so the game is still playable.
	if g.WidgetManager.Len() == 0 {
		RegisterDefaultUI(g)
	}

	// The opening screen goes in directly, and its OnEnter is queued with
	// whatever the world start adds, so the first script tick sees one
	// composed sequence.
	g.currentStage = g.StartStage
	if g.StartStage.World {
		// A game that opens straight into the world starts a session now;
		// one that opens on a menu starts it when the player says.
		g.startWorld()
	}
	g.arriveStage(g.StartStage)

	scale := g.WindowScale
	if scale <= 0 {
		scale = 4
	}
	ebiten.SetWindowSize(g.Width*scale, g.Height*scale)
	ebiten.SetWindowTitle(g.Title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(&engine{g: g})
}

// (Game.Run is provided here for convenience; some callers prefer it.)
func (g *Game) Run() error { return Run(g) }
