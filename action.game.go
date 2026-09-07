package inkwell

// action.game.go holds the adventure actions that access *Game directly and
// have not yet moved to a sub-package — either because they reference root
// types the sub-package cannot import (Stage, Scene, Script, AudioPlayer)
// or because the AudioPlayer type split prevents delegation to inkbase.

// ----- audio --------------------------------------------------------------

type playMusicAction struct{ music *Music }

func PlayMusic(music *Music) Action      { return &playMusicAction{music: music} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).Audio.PlayMusic(a.music)
	return StatusDone
}

type playSoundAction struct{ sound *Sound }

func PlaySound(sound *Sound) Action      { return &playSoundAction{sound: sound} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).Audio.PlaySound(a.sound)
	return StatusDone
}

type stopMusicAction struct{}

func StopMusic() Action                         { return &stopMusicAction{} }
func (a *stopMusicAction) Start() Runner        { return a }
func (a *stopMusicAction) Tick(ctx *Ctx) Status { gameOf(ctx).Audio.StopMusic(); return StatusDone }

// ----- stages -------------------------------------------------------------

type gotoStageAction struct{ stage *Stage }

func GoToStage(stage *Stage) Action      { return &gotoStageAction{stage: stage} }
func (a *gotoStageAction) Start() Runner { return a }
func (a *gotoStageAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).GoToStage(a.stage)
	return StatusDone
}

type goBackAction struct{}

func GoBack() Action                   { return &goBackAction{} }
func (a *goBackAction) Start() Runner  { return a }
func (a *goBackAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	if previous := game.PreviousStage(); previous != nil {
		game.GoToStage(previous)
	}
	return StatusDone
}

// ----- lifecycle ----------------------------------------------------------

type quitAction struct{}

func Quit() Action                     { return &quitAction{} }
func (a *quitAction) Start() Runner    { return a }
func (a *quitAction) Tick(ctx *Ctx) Status { gameOf(ctx).quitting = true; return StatusDone }

type restartAction struct{}

func Restart() Action                     { return &restartAction{} }
func (a *restartAction) Start() Runner    { return a }
func (a *restartAction) Tick(ctx *Ctx) Status { gameOf(ctx).Restart(); return StatusDone }

// ----- slots --------------------------------------------------------------

type loadSlotAction struct{ slot int }

func LoadSlot(slot int) Action           { return &loadSlotAction{slot: slot} }
func (a *loadSlotAction) Start() Runner  { return a }
func (a *loadSlotAction) Tick(ctx *Ctx) Status {
	gameOf(ctx).LoadLater(a.slot)
	return StatusDone
}

type saveSlotAction struct{ slot int }

func SaveSlot(slot int) Action           { return &saveSlotAction{slot: slot} }
func (a *saveSlotAction) Start() Runner  { return a }
func (a *saveSlotAction) Tick(ctx *Ctx) Status {
	if err := gameOf(ctx).Save(a.slot); err != nil {
		logf("saveSlot %d: %v", a.slot, err)
		return StatusFailed
	}
	return StatusDone
}

// ----- settings -----------------------------------------------------------

type useLangAction struct{ lang string }

func UseLang(lang string) Action           { return &useLangAction{lang: lang} }
func (a *useLangAction) Start() Runner     { return a }
func (a *useLangAction) Tick(ctx *Ctx) Status {
	Lang = a.lang
	_ = gameOf(ctx).SaveSettings()
	return StatusDone
}

type useThemeAction struct{ theme *Theme }

func UseTheme(theme *Theme) Action       { return &useThemeAction{theme: theme} }
func (a *useThemeAction) Start() Runner  { return a }
func (a *useThemeAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.UseTheme(a.theme)
	_ = game.SaveSettings()
	return StatusDone
}

// ----- script runner ------------------------------------------------------

type runScriptAction struct{ script *Script }

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
