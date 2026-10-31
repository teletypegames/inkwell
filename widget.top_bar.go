package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/topbar"

// TopBar is a thin strip across the top of the screen: the scene title on
// the left, a score in the centre and a time string on the right, in the
// active theme's top bar colours; see [topbar.TopBar].
type TopBar = topbar.TopBar

// TopBarColors is the palette a [TopBar] draws with, from the active theme's
// top bar tokens. It satisfies [topbar.Host].
func (g *Game) TopBarColors() topbar.Colors {
	th := g.Theme()
	return topbar.Colors{
		BG:     th.TopBarBG,
		Text:   th.TopBarText,
		Accent: th.TopBarAccent,
	}
}

// SceneTitle is what a bar calls the current scene: its Title, or its Name
// when it has none, and an empty Text before the first scene. It satisfies
// [topbar.Host].
func (g *Game) SceneTitle() Text {
	if g.currentScene == "" {
		return nil
	}
	s := g.SceneManager.MustGet(g.currentScene)
	if s.Title.String() != "" {
		return s.Title
	}
	return plain(s.Name)
}

// StateVar reads the State variable name, nil when it was never set. It
// satisfies [topbar.Host].
func (g *Game) StateVar(name string) any { return g.State.Var(name) }
