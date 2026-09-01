package inkwell

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Run validates the game, then enters the ebiten main loop. The window is
// sized to Game.WindowScale times the internal resolution, 4× by default.
func Run(g *Game) error {
	// The domain may have swapped a manager in since NewGame, so the parts
	// that hold a registry directly are wired here, not at construction.
	g.Audio.attach(g.AssetManager)

	if err := g.Validate(); err != nil {
		return err
	}
	// If the domain didn't register any widgets, fall back to the SCUMM
	// preset so the game is still playable.
	if g.WidgetManager.Len() == 0 {
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
	if g.onStart != "" {
		seq = append(seq, RunScript(g.onStart))
	}
	if g.onFinale != "" {
		seq = append(seq, RunScript(g.onFinale))
	}
	if len(seq) > 0 {
		g.queueAction(Seq(seq...), "init")
	}

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
