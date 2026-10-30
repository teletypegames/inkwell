// Package action holds the generic actions — the ones any Ebitengine game
// built on inkbase can use — and the combinators every script is built from.
//
// # Action and Runner
//
// An [inkbase.Action] is an immutable spec of work, safe to keep in a field
// and run many times. Its Start makes an [inkbase.Runner], one in-flight
// execution carrying whatever state it needs (elapsed time, sub-runners, a
// started flag). The host never ticks an Action directly; it asks for a
// fresh Runner every time the action runs.
//
// A Runner's Tick is called once per frame with an [inkbase.Ctx] — the
// frame's DT in seconds, the [inkbase.State], and the Game as an opaque any
// — and answers an [inkbase.Status]: StatusRunning (tick me again),
// StatusDone (a sequence moves on) or StatusFailed (a sequence stops here,
// which is how RequireItem aborts the rest of a script).
//
// An inkwell game runs one action at a time, from a queue; parallel work
// goes inside one action, with [Par].
//
// # The actions
//
//   - [Seq] runs its actions in order; nested sequences are flattened and
//     nil actions skipped. [Par] runs them together and is done when all
//     are. [If] runs then when the condition holds, else the rest in
//     sequence. [Wait] holds for a number of seconds.
//   - [SetFlag], [ClearFlag] and [SetVar] write the State. [SetVarIfEmpty]
//     writes a var only while nothing has written it yet — what a scene's
//     OnEnter wants for a first-arrival value. [Noting] writes a var for
//     exactly as long as an inner action runs, and clears it after.
//     [Learn] marks an entity known.
//   - [Flash] puts a line in the status strip ([Flasher]).
//   - [GoToStage] shows another Stage, and [GoBack] the one before it
//     ([StageHost]); [Restart] begins a new session ([LifecycleHost]);
//     [LoadSlot] and [SaveSlot] restore and write a save slot ([SlotHost]);
//     [Quit] closes the window at the end of the frame ([Quitter]). The
//     screen changes land at the top of the next frame.
//   - [UseTheme] ([ThemeHost]) and [UseLang] ([SettingsSaver]) switch the
//     theme or the language and save the settings.
//   - [PlayMusic], [StopMusic] and [PlaySound] drive the host's
//     [asset.AudioPlayer]
//     ([AudioHost]).
//   - [RunScript] runs a Script's actions.
//   - [Custom] wraps a function as an action that finishes in one tick.
//
// An action that needs more than the State type-asserts ctx.Game to the
// host interface named above; a host that does not implement it gets an
// action that does nothing.
//
// # Custom actions
//
// For a one-shot change, Custom is shortest:
//
//	bump := action.Custom(func(ctx *inkbase.Ctx) inkbase.Status {
//		score, _ := ctx.State.Var("score").(int)
//		ctx.State.SetVar("score", score+10)
//		return inkbase.StatusDone
//	})
//
// Work that runs over several frames — an animation, a tween, a poll —
// implements Action and Runner directly:
//
//	type fade struct{ seconds float64 }
//
//	func (f fade) Start() inkbase.Runner { return &fadeRunner{spec: f} }
//
//	type fadeRunner struct {
//		spec    fade
//		elapsed float64
//	}
//
//	func (r *fadeRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
//		r.elapsed += ctx.DT
//		if r.elapsed >= r.spec.seconds {
//			return inkbase.StatusDone
//		}
//		// change state, draw a fade, …
//		return inkbase.StatusRunning
//	}
package action
