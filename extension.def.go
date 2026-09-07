package inkwell

// Extension is a piece of machinery that lives beside the game for the whole
// run and is not a widget: a platform layer such as Steam, which has to be
// initialised before the window opens, pumped every frame whether a script
// holds the frame or not, and shut down when the loop returns. A widget
// stands still while a script runs and has no hook at either end, which is
// why this is its own contract.
//
// Start runs at the top of Run, before the settings document is read, so an
// extension may hand the Game a Store of its own and set SystemLang. An
// error from Start ends Run. Tick runs first thing in every Update, on every
// Stage. Stop runs once the ebiten loop has returned.
type Extension interface {
	Start(game *Game) error
	Tick(ctx *Ctx)
	Stop()
}

// Unlocker is an Extension that can record an Achievement earned. It is
// asked for by type assertion, the way an action asks its host, so an
// extension with no notion of achievements implements nothing extra.
type Unlocker interface {
	Unlock(achievement *Achievement)
}

func (g *Game) startExtensions() error {
	for _, extension := range g.Extensions {
		if err := extension.Start(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) tickExtensions(ctx *Ctx) {
	for _, extension := range g.Extensions {
		extension.Tick(ctx)
	}
}

func (g *Game) stopExtensions() {
	for i := len(g.Extensions) - 1; i >= 0; i-- {
		g.Extensions[i].Stop()
	}
}

func (g *Game) unlock(a *Achievement) {
	for _, extension := range g.Extensions {
		if unlocker, ok := extension.(Unlocker); ok {
			unlocker.Unlock(a)
		}
	}
}
