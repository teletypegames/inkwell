package pncdsl

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Run validates the game, then enters the ebiten main loop. The window is
// sized to 4× the internal resolution.
func Run(g *Game) error {
	if err := g.Validate(); err != nil {
		return err
	}
	// If the domain didn't register any widgets, fall back to the SCUMM
	// preset so the game is still playable.
	if g.UIManager.Len() == 0 {
		RegisterDefaultUI(g)
	}

	// Initial scene goes in directly (no transition) so OnEnter / OnStart
	// run cleanly as one composed sequence on the first script tick.
	s := g.SceneManager.MustGet(g.startID)
	g.currentScene = g.startID
	g.State.NoteVisit(g.startID)
	g.placeActors(g.startID)
	g.resetTriggers()
	if s.Music != "" {
		g.Audio.PlayMusic(s.Music)
	}
	var seq []Action
	if s.OnEnter != nil {
		seq = append(seq, s.OnEnter)
	}
	if g.onStart != nil {
		seq = append(seq, g.onStart)
	}
	if len(seq) > 0 {
		g.queueAction(Seq(seq...), "init")
	}

	ebiten.SetWindowSize(g.Width*4, g.Height*4)
	ebiten.SetWindowTitle(g.Title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(&engine{g: g})
}

// (Game.Run is provided here for convenience; some callers prefer it.)
func (g *Game) Run() error { return Run(g) }
