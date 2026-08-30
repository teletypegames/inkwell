package inkwell

// Status is the outcome of one tick of a Runner.
type Status int

const (
	StatusRunning Status = iota
	StatusDone
	StatusFailed
)

// Ctx is the per-tick context passed to running actions.
type Ctx struct {
	Game    *Game
	DT      float64
	Scene   *Scene
	Hotspot *Hotspot
	Item    *Item
}

// Action is an immutable spec of work; Start makes a fresh Runner with
// state. Storing an Action in a struct field (e.g. Hotspot.OnUse) is safe
// because the engine always calls Start before ticking, so two invocations
// can never share mutable state.
type Action interface {
	Start() Runner
}

// Runner is one in-flight execution of an Action. Tick advances it one
// frame and returns whether it's still running, done, or failed.
type Runner interface {
	Tick(ctx *Ctx) Status
}

// ----- immediate action helper -----------------------------------------

type immediateAction struct {
	fn func(*Ctx) Status
}

func (a *immediateAction) Start() Runner { return a }
func (a *immediateAction) Tick(ctx *Ctx) Status {
	if a.fn == nil {
		return StatusDone
	}
	return a.fn(ctx)
}

// Custom wraps a user function as an Action. Returns Done on first tick
// unless the function itself returns StatusRunning.
func Custom(fn func(*Ctx) Status) Action { return &immediateAction{fn: fn} }

// ----- Seq --------------------------------------------------------------

type seqAction struct{ children []Action }
type seqRunner struct {
	spec    *seqAction
	idx     int
	current Runner
}

func Seq(actions ...Action) Action {
	return &seqAction{children: flattenSeq(actions)}
}

func flattenSeq(in []Action) []Action {
	out := make([]Action, 0, len(in))
	for _, a := range in {
		if a == nil {
			continue
		}
		if s, ok := a.(*seqAction); ok {
			out = append(out, s.children...)
		} else {
			out = append(out, a)
		}
	}
	return out
}

func (a *seqAction) Start() Runner { return &seqRunner{spec: a} }

func (r *seqRunner) Tick(ctx *Ctx) Status {
	for {
		if r.idx >= len(r.spec.children) {
			return StatusDone
		}
		if r.current == nil {
			r.current = r.spec.children[r.idx].Start()
		}
		s := r.current.Tick(ctx)
		switch s {
		case StatusDone:
			r.idx++
			r.current = nil
			// keep going only if the just-finished action consumed no time
			// (zero-dt would loop forever otherwise — that's fine here since
			// immediate actions complete in one Tick call).
			continue
		case StatusFailed:
			return StatusFailed
		default:
			return StatusRunning
		}
	}
}

// ----- Par --------------------------------------------------------------

type parAction struct{ children []Action }
type parRunner struct {
	runners []Runner
	done    []bool
}

func Par(actions ...Action) Action { return &parAction{children: actions} }

func (a *parAction) Start() Runner {
	rs := make([]Runner, len(a.children))
	for i, c := range a.children {
		rs[i] = c.Start()
	}
	return &parRunner{runners: rs, done: make([]bool, len(rs))}
}

func (r *parRunner) Tick(ctx *Ctx) Status {
	allDone := true
	for i, rn := range r.runners {
		if r.done[i] {
			continue
		}
		s := rn.Tick(ctx)
		if s == StatusFailed {
			return StatusFailed
		}
		if s == StatusDone {
			r.done[i] = true
			continue
		}
		allDone = false
	}
	if allDone {
		return StatusDone
	}
	return StatusRunning
}

// ----- If ---------------------------------------------------------------

type ifAction struct {
	cond Condition
	then Action
	els  Action
}

func If(cond Condition, then Action, els ...Action) Action {
	var e Action
	if len(els) > 0 {
		e = Seq(els...)
	}
	return &ifAction{cond: cond, then: then, els: e}
}

func (a *ifAction) Start() Runner { return &ifRunner{spec: a} }

type ifRunner struct {
	spec    *ifAction
	started bool
	inner   Runner
}

func (r *ifRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		take := r.spec.then
		if r.spec.cond == nil || !r.spec.cond.Eval(ctx) {
			take = r.spec.els
		}
		if take == nil {
			return StatusDone
		}
		r.inner = take.Start()
	}
	if r.inner == nil {
		return StatusDone
	}
	return r.inner.Tick(ctx)
}

// ----- Wait -------------------------------------------------------------

type waitAction struct{ seconds float64 }
type waitRunner struct {
	spec    *waitAction
	elapsed float64
}

func Wait(seconds float64) Action { return &waitAction{seconds: seconds} }
func (a *waitAction) Start() Runner { return &waitRunner{spec: a} }
func (r *waitRunner) Tick(ctx *Ctx) Status {
	r.elapsed += ctx.DT
	if r.elapsed >= r.spec.seconds {
		return StatusDone
	}
	return StatusRunning
}

// ----- Say --------------------------------------------------------------

type sayAction struct{ speaker, text string }
type sayRunner struct {
	spec     *sayAction
	elapsed  float64
	duration float64
	started  bool
	bubble   bool
}

func Say(speaker, text string) Action { return &sayAction{speaker: speaker, text: text} }

func (a *sayAction) Start() Runner { return &sayRunner{spec: a} }

func (r *sayRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		// duration scales with text length, with a 1.2s floor
		r.duration = 1.2 + float64(len([]rune(r.spec.text)))*0.05
		r.bubble = ctx.Game.VoiceOf(r.spec.speaker) == VoiceBubble
		if r.bubble {
			ctx.Game.SetSpeech(r.spec.speaker, r.spec.text)
		}
		ctx.Game.LogResponse(r.spec.speaker, r.spec.text)
	}
	r.elapsed += ctx.DT
	// skip on click
	if ctx.Game.Input.consumedClick() {
		r.elapsed = r.duration
	}
	if r.elapsed >= r.duration {
		if r.bubble {
			ctx.Game.ClearSpeech()
		}
		return StatusDone
	}
	return StatusRunning
}

// ----- GoTo -------------------------------------------------------------

type gotoAction struct{ scene string }

func GoTo(scene string) Action { return &gotoAction{scene: scene} }
func (a *gotoAction) Start() Runner { return a }
func (a *gotoAction) Tick(ctx *Ctx) Status {
	ctx.Game.changeScene(a.scene)
	return StatusDone
}

// ----- Back -------------------------------------------------------------

type backAction struct{}

// Back returns to the scene the player came from. Most exits name their
// destination, because a door leads where it leads; a scene reached from
// several different rooms cannot, so its way out is a direction, not a place.
// A no-op when there is nowhere to go back to.
func Back() Action { return backAction{} }

func (a backAction) Start() Runner { return a }
func (a backAction) Tick(ctx *Ctx) Status {
	if prev := ctx.Game.PreviousScene(); prev != "" {
		ctx.Game.changeScene(prev)
	}
	return StatusDone
}

// ----- inventory --------------------------------------------------------

type giveAction struct{ item string }

func Give(item string) Action { return &giveAction{item: item} }
func (a *giveAction) Start() Runner { return a }
func (a *giveAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Add(a.item)
	return StatusDone
}

type takeAwayAction struct{ item string }

func TakeAway(item string) Action { return &takeAwayAction{item: item} }
func (a *takeAwayAction) Start() Runner { return a }
func (a *takeAwayAction) Tick(ctx *Ctx) Status {
	ctx.Game.Inventory.Remove(a.item)
	return StatusDone
}

// RequireItem fails silently with a generic line if the player doesn't have
// the named item selected or in inventory. Used at the top of Use handlers.
type requireItemAction struct{ item string }

func RequireItem(item string) Action { return &requireItemAction{item: item} }
func (a *requireItemAction) Start() Runner { return a }
func (a *requireItemAction) Tick(ctx *Ctx) Status {
	if ctx.Game.Inventory.Has(a.item) {
		return StatusDone
	}
	ctx.Game.FlashLine("Ehhez kell egy " + a.item + ".")
	return StatusFailed
}

// ----- flags / vars -----------------------------------------------------

type setFlagAction struct{ name string }

func SetFlag(name string) Action { return &setFlagAction{name: name} }
func (a *setFlagAction) Start() Runner { return a }
func (a *setFlagAction) Tick(ctx *Ctx) Status { ctx.Game.State.SetFlag(a.name); return StatusDone }

type clearFlagAction struct{ name string }

func ClearFlag(name string) Action { return &clearFlagAction{name: name} }
func (a *clearFlagAction) Start() Runner { return a }
func (a *clearFlagAction) Tick(ctx *Ctx) Status {
	ctx.Game.State.ClearFlag(a.name)
	return StatusDone
}

type setVarAction struct {
	name string
	v    any
}

func SetVar(name string, v any) Action { return &setVarAction{name: name, v: v} }
func (a *setVarAction) Start() Runner { return a }
func (a *setVarAction) Tick(ctx *Ctx) Status {
	ctx.Game.State.SetVar(a.name, a.v)
	return StatusDone
}

// ----- audio ------------------------------------------------------------

type playMusicAction struct{ name string }

func PlayMusic(name string) Action { return &playMusicAction{name: name} }
func (a *playMusicAction) Start() Runner { return a }
func (a *playMusicAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlayMusic(a.name)
	return StatusDone
}

type stopMusicAction struct{}

func StopMusic() Action { return &stopMusicAction{} }
func (a *stopMusicAction) Start() Runner { return a }
func (a *stopMusicAction) Tick(ctx *Ctx) Status { ctx.Game.Audio.StopMusic(); return StatusDone }

type playSoundAction struct{ name string }

func PlaySound(name string) Action { return &playSoundAction{name: name} }
func (a *playSoundAction) Start() Runner { return a }
func (a *playSoundAction) Tick(ctx *Ctx) Status {
	ctx.Game.Audio.PlaySound(a.name)
	return StatusDone
}

// ----- dialogue ---------------------------------------------------------

type runDialogueAction struct{ name string }

func RunDialogue(name string) Action { return &runDialogueAction{name: name} }
func (a *runDialogueAction) Start() Runner { return &runDialogueRunner{spec: a} }

type runDialogueRunner struct {
	spec    *runDialogueAction
	started bool
}

func (r *runDialogueRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.startDialogue(r.spec.name)
	}
	if ctx.Game.dialogueActive() {
		return StatusRunning
	}
	return StatusDone
}

type endDialogueAction struct{}

func EndDialogue() Action { return &endDialogueAction{} }
func (a *endDialogueAction) Start() Runner { return a }
func (a *endDialogueAction) Tick(ctx *Ctx) Status {
	ctx.Game.endDialogue()
	return StatusDone
}

type gotoNodeAction struct{ node string }

func GotoNode(node string) Action { return &gotoNodeAction{node: node} }
func (a *gotoNodeAction) Start() Runner { return a }
func (a *gotoNodeAction) Tick(ctx *Ctx) Status {
	ctx.Game.gotoDialogueNode(a.node)
	return StatusDone
}

// ----- scripts ----------------------------------------------------------

type runScriptAction struct{ name string }

func RunScript(name string) Action { return &runScriptAction{name: name} }
func (a *runScriptAction) Start() Runner { return &runScriptRunner{spec: a} }

type runScriptRunner struct {
	spec  *runScriptAction
	inner Runner
}

func (r *runScriptRunner) Tick(ctx *Ctx) Status {
	if r.inner == nil {
		s, ok := ctx.Game.ScriptManager.Get(r.spec.name)
		if !ok || s.Actions == nil {
			return StatusFailed
		}
		r.inner = s.Actions.Start()
	}
	return r.inner.Tick(ctx)
}

// ----- character movement ----------------------------------------------

type walkAction struct {
	character string
	to        Point
}

func Walk(character string, to Point) Action { return &walkAction{character: character, to: to} }

func (a *walkAction) Start() Runner { return &walkRunner{spec: a} }

type walkRunner struct {
	spec    *walkAction
	started bool
}

func (r *walkRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.walkCharacter(r.spec.character, r.spec.to)
	}
	if ctx.Game.characterMoving(r.spec.character) {
		return StatusRunning
	}
	return StatusDone
}

// ----- misc -------------------------------------------------------------

type showEndAction struct{ text string }

func ShowEnd(text string) Action { return &showEndAction{text: text} }
func (a *showEndAction) Start() Runner { return &showEndRunner{spec: a} }

type showEndRunner struct {
	spec    *showEndAction
	started bool
}

func (r *showEndRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		ctx.Game.showEndCard(r.spec.text)
	}
	return StatusRunning // never finishes; player closes the window
}
