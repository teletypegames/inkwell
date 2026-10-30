package inkwell

// action.game.go holds the adventure actions that access *Game directly and
// have not yet moved to a sub-package — either because they reference root
// types the sub-package cannot import (Stage, Scene, Script, AudioPlayer)
// or because the AudioPlayer type split prevents delegation to inkbase.

// ----- audio --------------------------------------------------------------

type playMusicAction struct{ music *Music }

// PlayMusic is an Action that starts music in place of the current track;
// see [AudioPlayer.PlayMusic]. It is done at once.
func PlayMusic(music *Music) Action      { return &playMusicAction{music: music} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).Audio.PlayMusic(a.music)
	return StatusDone
}

type playSoundAction struct{ sound *Sound }

// PlaySound is an Action that plays sound once; see
// [AudioPlayer.PlaySound]. It is done at once, without waiting for the
// sound to end.
func PlaySound(sound *Sound) Action      { return &playSoundAction{sound: sound} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).Audio.PlaySound(a.sound)
	return StatusDone
}

type stopMusicAction struct{}

// StopMusic is an Action that stops the current track; see
// [AudioPlayer.StopMusic].
func StopMusic() Action                         { return &stopMusicAction{} }
func (a *stopMusicAction) Start() Runner        { return a }
func (a *stopMusicAction) Tick(ctx *Ctx) Status { gameOf(ctx).Audio.StopMusic(); return StatusDone }

// ----- stages -------------------------------------------------------------

type gotoStageAction struct{ stage *Stage }

// GoToStage is an Action that shows stage, at the top of the next frame;
// see [Game.GoToStage]. It is done at once.
func GoToStage(stage *Stage) Action      { return &gotoStageAction{stage: stage} }
func (a *gotoStageAction) Start() Runner { return a }
func (a *gotoStageAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).GoToStage(a.stage)
	return StatusDone
}

type goBackAction struct{}

// GoBack is an Action that returns to the previous Stage, or does nothing
// when there is none; see [Game.PreviousStage].
func GoBack() Action                  { return &goBackAction{} }
func (a *goBackAction) Start() Runner { return a }
func (a *goBackAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	if previous := game.PreviousStage(); previous != nil {
		game.GoToStage(previous)
	}
	return StatusDone
}

// ----- lifecycle ----------------------------------------------------------

type quitAction struct{}

// Quit is an Action that ends the game: the main loop exits at the top of
// the next frame.
func Quit() Action                         { return &quitAction{} }
func (a *quitAction) Start() Runner        { return a }
func (a *quitAction) Tick(ctx *Ctx) Status { gameOf(ctx).quitting = true; return StatusDone }

type restartAction struct{}

// Restart is an Action that begins a new session, at the top of the next
// frame; see [Game.Restart].
func Restart() Action                         { return &restartAction{} }
func (a *restartAction) Start() Runner        { return a }
func (a *restartAction) Tick(ctx *Ctx) Status { gameOf(ctx).Restart(); return StatusDone }

// ----- slots --------------------------------------------------------------

type loadSlotAction struct{ slot int }

// LoadSlot is an Action that restores the session saved in slot, at the
// top of the next frame; see [Game.LoadLater].
func LoadSlot(slot int) Action          { return &loadSlotAction{slot: slot} }
func (a *loadSlotAction) Start() Runner { return a }
func (a *loadSlotAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).LoadLater(a.slot)
	return StatusDone
}

type saveSlotAction struct{ slot int }

// SaveSlot is an Action that saves the session into slot at once; it fails
// when [Game.Save] returns an error, which goes to the debug log.
func SaveSlot(slot int) Action          { return &saveSlotAction{slot: slot} }
func (a *saveSlotAction) Start() Runner { return a }
func (a *saveSlotAction) Tick(ctx *Ctx) Status {
	if err := gameOf(ctx).Save(a.slot); err != nil {
		logf("saveSlot %d: %v", a.slot, err)
		return StatusFailed
	}
	return StatusDone
}

// ----- settings -----------------------------------------------------------

type useLangAction struct{ lang string }

// UseLang is an Action that sets [Lang] to lang and saves the settings;
// a failed save is ignored.
func UseLang(lang string) Action       { return &useLangAction{lang: lang} }
func (a *useLangAction) Start() Runner { return a }
func (a *useLangAction) Tick(ctx *Ctx) Status {
	Lang = a.lang
	_ = gameOf(ctx).SaveSettings()
	return StatusDone
}

type useThemeAction struct{ theme *Theme }

// UseTheme is an Action that makes theme the active one and saves the
// settings; see [Game.UseTheme]. A failed save is ignored.
func UseTheme(theme *Theme) Action      { return &useThemeAction{theme: theme} }
func (a *useThemeAction) Start() Runner { return a }
func (a *useThemeAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.UseTheme(a.theme)
	_ = game.SaveSettings()
	return StatusDone
}

// ----- script runner ------------------------------------------------------

type runScriptAction struct{ script *Script }

// RunScript is an Action that runs the script's Actions to the end and
// takes their Status; a nil script, or one with no Actions, fails.
func RunScript(script *Script) Action    { return &runScriptAction{script: script} }
func (a *runScriptAction) Start() Runner { return &runScriptRunner{spec: a} }

type runScriptRunner struct {
	spec  *runScriptAction
	inner Runner
}

func (r *runScriptRunner) Tick(ctx *Ctx) Status {
	if r.inner == nil {
		if r.spec.script == nil || r.spec.script.Actions == nil {
			return StatusFailed
		}
		r.inner = r.spec.script.Actions.Start()
	}
	return r.inner.Tick(ctx)
}
