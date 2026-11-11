// Package inkwell is a point-and-click adventure game framework for Go,
// built on top of Ebitengine (https://ebitengine.org/). A game is composed by
// registering plain struct literals into per-entity [Manager] registries
// hanging off a single [Game] root. The library handles input, rendering,
// dialog trees, cutscenes, the HUD, themes and lazy asset loading; the game
// only declares content.
//
// inkwell is not a generic engine. It is a thin Go-level layer on top of
// Ebitengine that implements the fixed skeleton of the point-and-click genre
// as it existed in the LucasArts SCUMM era (Maniac Mansion, Monkey Island,
// Day of the Tentacle):
//
//   - scenes with hotspots and click zones,
//   - a verb-based interaction grammar (Look, Use, Talk, Take, …),
//   - an inventory, item-on-item and item-on-hotspot combinations,
//   - dialog trees with conditional choices,
//   - a scripting and cutscene system built from composable [Action] values,
//   - a widget-tree HUD that can be extended with anything Ebitengine can draw,
//   - swappable colour themes,
//   - placeholder graphics, so the game runs without any art on disk.
//
// The generic infrastructure — registries, text in several languages, world
// state, the action and condition framework and the widget contract — lives
// in the inkbase package and its sub-packages, which a game that is not a
// point-and-click adventure can import alone. inkwell builds on inkbase and
// re-exports what a game needs, so an adventure imports this one package.
//
// The design collapses to one organising pattern: every named entity is
// registered with Manager.Register(&Entity{Name: "..."}). No builders, no
// fluent chains, no scattered registries. The same mental model applies to
// items, scenes, characters, dialogs, scripts, assets, verbs, widgets and
// themes alike.
//
// A complete example game ("Morning Coffee") lives in the separate
// inkwell-demo repository, ssh://git@git.teletypegames.org:2222/demos/inkwell-demo,
// and consumes this library as a plain Go module.
//
// # Quick start
//
// The module needs the Go version named in its go.mod and a working OpenGL
// context; Ebitengine handles the window and the render loop.
//
//	go mod init my-game
//	go get git.teletypegames.org/engines/inkwell
//
// If the host cannot be reached over HTTPS, point Go at the SSH endpoint:
//
//	git config --global \
//	    url."ssh://git@git.teletypegames.org:2222/".insteadOf "https://git.teletypegames.org/"
//	export GOPRIVATE=git.teletypegames.org/*
//
// The minimum game:
//
//	package main
//
//	import (
//		"log"
//
//		"git.teletypegames.org/engines/inkwell"
//	)
//
//	func main() {
//		player := &inkwell.Character{Name: "player", W: 28, H: 62}
//		room := &inkwell.Scene{
//			Name:   "room",
//			Layers: []inkwell.SceneLayer{{Asset: &inkwell.Asset{Name: "bg/room", Path: "assets/bg/room.png"}, Z: -1}},
//			Actors: []inkwell.SceneActor{{Character: player, At: inkwell.Point{X: 160, Y: 140}}},
//			Hotspots: []inkwell.Hotspot{{
//				Name:  "door",
//				Area:  inkwell.Rect(260, 50, 40, 90),
//				Label: inkwell.Text{{Lang: "en", Text: "door"}},
//				OnLook: inkwell.Say(player, inkwell.Text{
//					{Lang: "en", Text: "A wooden door."},
//				}),
//			}},
//		}
//		world := &inkwell.Stage{Name: "game", World: true}
//
//		g := inkwell.New(inkwell.Game{
//			Title:      "Sample",
//			Width:      320,
//			Height:     200,
//			StartStage: world,
//			StartScene: room,
//			Player:     player,
//		})
//		g.CharacterManager.Register(player)
//		g.SceneManager.Register(room)
//		g.StageManager.Register(world)
//
//		inkwell.RegisterDefaultUI(g) // optional: Run installs it when no widget is registered
//		if err := inkwell.Run(g); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # The Manager pattern
//
// Every name-addressable entity goes through one generic registry,
// [Manager], over anything that is [Named] (it has GetName). Each entity kind
// has its own alias: [ItemManager], [SceneManager], [CharacterManager],
// [SpeakerManager], [DialogManager], [ScriptManager], [AssetManager],
// [VerbManager], [WidgetManager], [ThemeManager], [StageManager] and
// [AchievementManager].
//
// A Manager offers:
//
//   - Register(v) adds v and hands it back; it panics on an empty or
//     duplicate name.
//   - Set(v) replaces an entry in place, keeping its position, and registers
//     it when the name is new.
//   - Get(name) looks an entry up, with a boolean that is false when absent;
//     MustGet panics instead.
//   - Has, Len, and Names, All and Each, which all walk the entries in
//     insertion order (the order widgets are drawn in within one layer).
//
// Register returns what it registered so that a registration can be an
// expression rather than a statement, which lets an entity file be a
// declaration instead of a procedure:
//
//	// the whole of a scene file
//	package scene
//
//	var Alley = SceneManager.Register(&inkwell.Scene{
//		Name: "alley",
//		…
//	})
//
// A package-level variable is initialised in dependency order, so the
// manager is built before any file that writes into it, whichever file the
// compiler reads first. Registration order is the order the files are
// presented in, as it would be with init functions.
//
// The result is the handle. Every manager holds pointers — Manager[*Scene],
// not Manager[Scene] — so the variable a file binds and the entry the engine
// looks up are one instance, and every field that refers to another entity
// is a pointer to it: Stage.Background is an *Asset, Exit.To a *Scene,
// SceneActor.Character a *Character. A reference is therefore an identifier,
// checked by the compiler and followed by an editor, and Name is left with
// the one job only a string can do: it is the key a save file, a [State]
// var and a command-line flag use to find the entity again.
//
// The result is also safe to ignore, so the plain statement form still
// compiles, and a file that has to build something before registering it (a
// loop, or a value wired to itself) keeps its init function.
//
// Most registerable types also implement TypeLabel, which the manager uses
// in its panic messages ("duplicate item …"). A duplicate registration is a
// construction-time bug, code that ran once on startup; crashing loudly
// surfaces it in development, where quietly accepting the second one would
// mask a shadowed entity at runtime. When an entry genuinely has to be
// rewritten — a derived field filled in after the fact, a hot-reloaded
// entity — Set is the deliberate overwrite, and it keeps the entry's place
// in the insertion order.
//
// # Building a game
//
// [Game] is the root aggregate: the title and resolution, one manager per
// entity kind, what the game boots as, the domain-wide defaults and the
// runtime services ([State], [Inventory], [AudioPlayer], [Input]). [New]
// builds one from a declared spec. Everything a game is rather than does is
// a field, so the whole of it is one literal at the call site:
//
//	game := inkwell.New(inkwell.Game{
//		Title:       "Real World",
//		Width:       640,
//		Height:      380,
//		WindowScale: 2,
//
//		StartStage:  stage.Game,
//		StartScene:  scene.Alley,
//		StartScript: script.Opening,
//		StartTheme:  theme.RealWorld93,
//
//		Player:      character.Paul,
//		Walkboxes:   scene.Floor,
//		UseWithFail: script.UseWithFail,
//
//		SceneManager:  scene.SceneManager,
//		ScriptManager: script.ScriptManager,
//		WidgetManager: widget.WidgetManager,
//		ThemeManager:  theme.ThemeManager,
//	})
//
// What the spec leaves out, New fills in: an empty manager for every kind the
// caller did not hand over, the runtime halves (State, Inventory, Audio,
// Input), the SCUMM verb set (look, use, talk, take) when no verbs were
// registered, and the four preset themes. Widgets are not registered
// automatically: call [RegisterDefaultUI] (or one of its siblings), or let
// [Run] install the default set when the WidgetManager is empty.
//
// A manager handed over is used in place, not copied, so the engine and the
// game share one registry per kind. New adds the preset themes to whichever
// ThemeManager it ends up with, leaving alone any name the game has already
// taken, and a VerbManager handed over with verbs in it does not get the
// SCUMM four installed over it.
//
// [NewGame] is New for a game built in steps: it declares nothing but the
// title and the window, and the caller assigns the rest before Run, which is
// where the engine wires up the parts that hold a registry directly.
//
//	g := inkwell.NewGame("Real World", 640, 380)
//	g.SceneManager = mySceneManager
//	g.AssetManager = myAssetManager
//
// The fields that say what the game boots as:
//
//   - StartStage is the screen the game opens on — its main menu, or the
//     world itself.
//   - StartScene is the scene a session opens in.
//   - StartScript is queued after that scene's OnEnter, once; FinaleScript
//     is queued straight after it, which is what a "boot into the ending"
//     debug flag wants.
//   - StartTheme is selected at boot; nil takes "classic-scumm".
//   - Flags are set in State before the first frame.
//   - SystemLang is the language the machine or platform reports (see
//     Text and languages).
//
// [Game.StartAt], [Game.OnStart] and [Game.OnFinale] are setters for
// StartScene, StartScript and FinaleScript for a game built in steps; they
// return the *Game so they chain. The opening and the finale are registered
// [Script] values rather than bare actions, so the opening a game plays is
// content like any other.
//
// The domain-wide defaults:
//
//   - Player is the character a scene with no actors of its own gets, and
//     the one WalkOnClick moves.
//   - WalkOnClick makes a left click on empty floor walk Player there,
//     routed through the walkboxes. It is off by default.
//   - Walkboxes is the floor of a scene that declares none.
//   - SceneRect is the part of the window the picture occupies; a game with
//     a HUD along the foot sets it so that exit strips land in the painting.
//   - ExitLook and ExitTake answer look and take on every exit, and
//     UseWithFail answers an item used on a hotspot nobody paired it with.
//   - WindowScale multiplies the resolution for the window; 0 means 4.
//   - MaxLogLines bounds the chat log; 0 keeps every line.
//   - App names the game to the machine (its settings directory, its
//     localStorage prefix) and Store overrides where settings and saves go.
//
// # Running a game
//
// [Game.Theme] is the active theme and [Game.UseTheme] switches it; the next
// frame uses the new colours, and the theme is remembered by name, which is
// what the settings store.
//
// The UI's transient state lives on the Game, written by actions and the
// engine and read by widgets: the selected verb ([Game.SelectedVerb]), the
// hover label and the composed "verb + target" hint ([Game.HoverHint]), the
// line being spoken ([Game.SetSpeech]) and the two-second status banner
// ([Game.FlashLine]), which the game loop expires whether or not a
// [StatusLine] shows it.
//
// The chat log is a buffer of [LogMessage] values, each a [LogKind]: player
// commands ([Game.LogAction]), in-world replies ([Game.LogResponse]) and
// system lines ([Game.LogSystem]). The engine and Say push into it,
// [Game.Messages] reads it, and the [ChatLog] widget draws it.
//
// The scene helpers:
//
//   - [Game.CurrentScene] is where the player is, empty before the first
//     scene is entered; [Game.PreviousScene] is the one before it, where
//     [Back] leads. Both survive a save.
//   - [Game.SceneArea] is SceneRect, or the whole window when it was never
//     set; exit strips are measured against it.
//   - [Game.SceneHotspots] is everything clickable in a scene: its own
//     hotspots first, then one per exit, so a painted thing beats the edge
//     strip an exit sits on wherever the two overlap. [Game.HotspotAt] hit
//     tests the current scene.
//   - [Game.InvokeVerb] is the one path from "the player picked a verb and a
//     target" to an action on the queue (see Hotspot). It is exported
//     because a target does not have to be in the picture: a HUD widget
//     that owns something clickable hands its own [Hotspot] here and gets
//     the same logging and resolution as a painted one.
//   - [Game.CharacterInScene] is what lets a [CharacterPanel] hide when its
//     character is not standing in the current scene.
//
// The screen helpers: [Game.CurrentStage] and [Game.PreviousStage],
// [Game.OnWorldStage], [Game.Started] (does a session exist?), and
// [Game.GoToStage], [Game.Restart] and [Game.LoadLater], which all land at
// the top of the next frame, and [Game.Quit].
//
// # The action queue
//
// [Game.Do] queues an action. Queued actions run one at a time and in order:
// the engine starts the next one on the first frame the previous one is done,
// so a widget, a trigger and a hotspot click can all hand work in without
// any of them being dropped or racing the others. The engine's own hotspot
// actions go through the same queue. The action running holds
// the frame — the script slot — so no ordinary widget ticks and the scene
// takes no clicks until it is done; parallel work goes inside one action with
// [Par]. See [git.teletypegames.org/engines/inkwell/act] and [git.teletypegames.org/engines/inkwell/inkbase/action] for the
// action model.
//
// # Extensions and Steam
//
// An [Extension] is machinery that lives beside the game for the whole run
// and is not a widget: a platform layer, which has to be up before the window
// opens, pumped every frame whether a script holds the frame or not, and shut
// down when the loop returns. Game.Extensions lists them.
//
// Start runs at the top of Run, before the settings are read, so an
// extension may hand the game a [Store] of its own and may set
// Game.SystemLang. An error from Start ends Run. Tick runs first thing in
// every frame, on every Stage. Stop runs in reverse order once the loop has
// returned, whichever way it returned. An extension that also implements
// [Unlocker] is handed every [Achievement] the world earns; it is asked for
// by type assertion, so an extension with no notion of achievements
// implements nothing extra.
//
// The one extension the library ships is the Steamworks layer in
// [git.teletypegames.org/engines/inkwell/steam], which this package does not
// import — a game that never names it links none of it. It is one literal in
// the Game:
//
//	Extensions: []inkwell.Extension{
//		&steam.Steam{
//			AppID:     480,
//			Languages: map[string]string{"english": "en", "hungarian": "hu"},
//			Cloud:     []string{inkwell.SlotsKey},
//		},
//	},
//
// See that package for what it does at each step and what it needs at
// runtime.
//
// # Text and languages
//
// Every string the player reads is a [Text]: the same string in as many
// languages as somebody has written it in, one [Translation] per language.
// The wording stays where the thing it names is — a hotspot's label in the
// hotspot, a line of dialog in the scene that speaks it — and nothing is
// lifted into a catalogue keyed by identifier, so a translator reading a
// scene file can see what the line is answering. One translation per line,
// stacked, reads as a table:
//
//	inkwell.Say(paul, inkwell.Text{
//		{Lang: "en", Text: "Locked. And I'm not guessing four digits."},
//		{Lang: "hu", Text: "Zárva. És nem fogok négy számjegyet találgatni."},
//	})
//
// Adding a language is adding a line. The zero value says nothing in any
// language, and a Text with one language behaves like a plain string. Inside
// a []Text or a map[string]Text the element type is implied, which is what
// gofmt -s leaves behind.
//
// String resolves a Text against the package variable [Lang], falls back to
// [Fallback], and falls back again to whichever translation was written
// first, so a line nobody has translated yet still reaches the player. In
// resolves against a named language. Set both variables once, before Run:
//
//	inkwell.Fallback = "en"
//	inkwell.Lang = *langFlag
//
// They are package variables rather than Game fields because a Text resolves
// itself wherever it is read — in a widget's Draw, in the Say runner, in a
// hover label — and threading a *Game through all of those would add a
// parameter for a value there is only ever one of. [AssetFS] is a package
// variable for the same reason. Resolution happens at display time, so
// content built during initialisation is unaffected by when the language is
// set. Run picks the language in this order: whatever was set before Run (a
// -lang flag), the stored setting, Game.SystemLang, Fallback. [UseLang]
// switches it at runtime and saves the setting.
//
// A sentence built around another Text has to come out in one language
// throughout, which is what [Format] is for: the frame and every Text
// argument resolve in the same language, and anything else is printed as it
// stands. It is not a join — English and Hungarian put the preposition in
// different places, and a frame per language is the only thing that survives
// contact with a second one:
//
//	var exitLook = inkwell.Text{
//		{Lang: "en", Text: "That way: %s."},
//		{Lang: "hu", Text: "Arra van: %s."},
//	}
//
//	g.ExitLook = func(e inkwell.Exit) inkwell.Action {
//		return inkwell.Say(paul, inkwell.Format(exitLook, e.Label))
//	}
//
// Map is the other composer: it rewrites every translation and keeps the
// languages, which is what a derived label wants.
//
// Two things are keyed by something other than wording, because wording is
// not a stable identity. A Once dialog choice is recorded against its
// position in its node, so a used-up choice stays used up in another language
// and a translator can reword it without breaking a save. And a [TopBar]
// note is a State var whose value is identical in every language, with
// TopBar.Notes mapping it to what the player reads.
//
// # What the engine says itself
//
// The handful of sentences the engine puts on screen without the game having
// written them are package variables in English and Hungarian, so a game in
// another language can teach them one by appending a [Translation]:
//
//   - [MsgNothing]: a verb nothing is bound to.
//   - [MsgNotThat]: an unpaired use-with, when UseWithFail is nil.
//   - [MsgNeedItem]: [RequireItem] on an item the player has not got.
//   - [HintVerb] and [HintVerbOn]: the status line over a plain target, and
//     while an item is held.
//   - [ScoreLabel]: the [TopBar] score prefix.
//   - [MenuEmpty], [MenuKeys] and [MenuKeysSticky]: a page with no rows, and
//     the key legend along the foot of a [Menu].
//
// For example:
//
//	inkwell.MsgNothing = append(inkwell.MsgNothing, inkwell.Translation{
//		Lang: "de",
//		Text: "Nichts von Bedeutung.",
//	})
//
// The four preset verbs carry both languages too; a game that wants a third
// replaces them through VerbManager.Set.
//
// # The font
//
// A game in a second language needs a font that can draw it. The engine
// ships its own bitmap atlas, baked from bitmapfont by the internal/fontgen
// generator, covering ASCII, the Latin-1 supplement, Latin Extended-A (ő and
// ű, and most of what Polish, Czech, Croatian and Turkish need), dashes,
// curly quotes, the bullet and the ellipsis, the euro sign and the four
// arrows — 348 glyphs in under 8 KB. A rune outside the coverage draws a
// hollow box, so a gap is visible rather than silent.
//
// Every glyph sits in a fixed [GlyphW] × [GlyphH] (6×16) cell, a white glyph
// over a black drop shadow, so [Game.DrawText] tints the glyph with the
// colour it is given and leaves the shadow black. [TextWidth], [WrapText]
// (word wrap; a newline breaks where it stands) and [ClipText] (cut to width
// with an ellipsis) are what the built-in widgets lay out with, and every
// vertical measurement in a HUD should derive from GlyphH. They take a
// resolved string, not a Text: a Text reaches them through String.
//
// To widen the coverage, edit the generator's ranges and run go generate in
// the module root (the generator is its own module, so the engine's
// dependency list stays short); it refuses any rune the 6-pixel cell cannot
// hold.
//
// # A ladder of lines
//
// A [Ladder] is what somebody says the n-th time the player does the same
// fruitless thing. Answering a repeated failure with one line is the surest
// way to sound like a machine; a Ladder answers the first attempt with the
// first rung and keeps a character's patience visibly running out:
//
//	var refuse = inkwell.Ladder{
//		Rungs: []inkwell.Text{
//			{{Lang: "en", Text: "No. That's not going to work."}},
//			{{Lang: "en", Text: "Tried that. It didn't improve."}},
//		},
//		Last: inkwell.Text{{Lang: "en", Text: "Leave it. It is not that kind of puzzle."}},
//	}
//
//	inkwell.Say(paul, refuse.At(state.Talked(key)))
//
// [Ladder.At] counts the first attempt as 1. Past the last rung it picks one
// at random, unless it names a Last, which stands from then on. The counter
// is the game's: State.NoteTalked and State.Talked already count how many
// times a thing has been worried at, keyed however the game wants.
// Game.UseWithFail is the usual caller.
//
// # Geometry
//
// [Point] is a position in game pixels. [Shape] is anything with Contains
// and Bounds; [Rectangle] (made with [Rect]) and [Polygon] (made with [Poly],
// hit-tested by ray casting) both are, so either can be a Hotspot.Area.
// [PolyRect] is a rectangle as a four-corner polygon: a walkbox is a Polygon
// so that a floor can be any shape, but most floors start as a box.
//
// # Assets and audio
//
// An [Asset] is a picture: a logical Name ("bg/kitchen") and a Path. Assets
// are lazy — registering one only stores the spec, and the file is opened on
// first use and cached. A missing or unreadable file falls back to a
// deterministic placeholder colour derived from the name, so the game runs
// without any art, which is useful for prototyping and CI; characters
// without art are drawn as a stylised figure instead. PNG and JPEG are
// decoded. Files are read through [AssetFS] — an embed.FS, or nil for the
// operating system's filesystem. The font is not an asset; it is baked in.
//
// Sound is not an Asset. Three things are heard, and they are three types
// because they are played three ways:
//
//   - A [Voice] is a recorded line. It sits on a [Translation], because the
//     English and the Hungarian recording are two files, and Say plays the
//     one the resolved language carries, for exactly as long as it lasts, in
//     place of the text-length timer. A skipped line stops talking. The
//     player switches voices off in the settings ([SetVoice]), and Say falls
//     back to the timer.
//   - A [Music] is a track a [Scene] or a [Stage] plays while it shows,
//     looping when Loop says so. Entering one whose Music is already playing
//     changes nothing, and a nil Music keeps whatever plays. [PlayMusic] and
//     [StopMusic] are the same as actions.
//   - A [Sound] is a noise — a door, a latch — played once by [PlaySound],
//     over the music and over other sounds.
//
// The [AudioPlayer] (Game.Audio) creates the shared Ebitengine audio context
// lazily at 48 kHz, or reuses one a host process already made. Music streams
// from disk; voices and sounds are decoded once and cached. WAV, Ogg Vorbis
// and MP3 are decoded, picked by file extension. Every failure — a missing
// file, an unsupported codec, no audio device — degrades to a debug-log line
// and silence.
//
// # Scenes, exits and hotspots
//
// A [Scene] is a place in the fiction: a picture in planes (Layers), the
// things in it a verb can reach (Hotspots), the ways out of it (Exits), the
// floor (Walkboxes), Triggers, the people standing in it (Actors), a Music
// track, and OnEnter and OnLeave actions.
//
// Layers are full-screen images, each with transparency where the planes
// behind show through, at a depth Z; the cast stands at 0. The background is
// the layer at -1, which everybody walks in front of; a foreground — a fence
// along the bottom edge, a counter — is the layer at 1, which everybody walks
// behind. A layer's Asset need not be registered, a missing lowest layer is
// the usual placeholder, a missing layer above it draws nothing, and an
// empty Layers leaves the theme's backdrop showing.
//
// A nil Actors means the usual cast, Game.Player at that character's Start;
// an empty non-nil slice is an empty scene. Only the current scene's cast is
// drawn, less anyone taken out with [Hide] — which lasts across scenes until
// [Show], so a cutscene can pass through a scene without the player:
//
//	inkwell.Seq(
//		inkwell.Hide(character.Paul),
//		inkwell.GoTo(scene.Rooftop),
//		inkwell.Say(speaker.Dex, inkwell.Text{{Lang: "en", Text: "He never saw it coming."}}),
//		inkwell.GoTo(scene.Alley),
//		inkwell.Show(character.Paul),
//	)
//
// A Show after the GoTo brings the player in once the new picture is
// uncovered; one before it has the player standing there as it uncovers,
// but also in the old scene while it fades, when that scene's cast includes
// the player. Until the Show the player stays hidden, and a save keeps it so. Walkboxes work the same way: nil
// falls back to Game.Walkboxes, an empty slice turns routing off. With
// walkboxes, a walk is routed by a breadth-first search over the polygons
// that share an edge, through the midpoint of each shared edge, and a
// destination outside every walkbox is clipped to the nearest boundary;
// without them, characters walk in a straight line.
//
// Exit.To and GoTo hold the *Scene itself, and two package-level variables
// may not refer to each other in their initialisers, which two scenes that
// lead to each other would. So a scene file declares the scene first and its
// exits second, with [Scene.Leads], which appends to Exits and hands the
// scene back:
//
//	var Alley = SceneManager.Register(&inkwell.Scene{Name: "alley", …})
//
//	var _ = Alley.Leads(
//		inkwell.Exit{To: Street, Side: inkwell.ExitLeft},
//		inkwell.Exit{
//			Label:   inkwell.Text{{Lang: "en", Text: "the fire escape"}},
//			Side:    inkwell.ExitBack,
//			Needs:   "has_ladder",
//			Blocked: inkwell.Say(paul, inkwell.Text{{Lang: "en", Text: "Can't reach it."}}),
//		},
//	)
//
// Scene.Derived hands back exits worked out from the rest of the game rather
// than written into the scene — the pins on a map screen, every location that
// says it leads there. Run asks each scene once, before validation, and
// appends the result to Exits, so everything downstream sees one kind of
// exit.
//
// The engine turns each [Exit] into a Hotspot with the exit cursor, the
// exit's Label, and a use that travels: [GoTo] the target, or [Back] when To
// is nil. When Needs names a flag, the journey is wrapped in
// If(Flag(Needs), travel, Blocked) — the exit is visible either way, since a
// locked door still says it is a door. With no Area an exit is a strip along
// one edge of [Game.SceneArea], chosen by its [ExitSide]. Look and take come
// from Game.ExitLook and Game.ExitTake unless the exit says otherwise, so
// that phrasing lives in one place. The generated hotspot is named
// [ExitName] of the target ("exit:noodle_house", or "exit:back"), so the
// location graph can be read back out of the registered scenes.
//
// A [Hotspot] is a named, clickable [Shape] with a hover Label, a
// [CursorKind] and one action per verb: OnLook, OnUse, OnTalk, OnTake, OnGive,
// and OnVerb for custom verbs. Hotspots belong to their scene and have no
// manager. A click goes through [Game.InvokeVerb]:
//
//   - With an item in hand: the hotspot's OnUseWith entry for that item,
//     then its OnUseWithAny — one answer for anything put on a slot, a lock
//     or a person, without enumerating the item catalogue — then
//     Game.UseWithFail, then the [MsgNotThat] flash. The item is deselected.
//   - Otherwise: the hotspot's handler for the selected verb, then the
//     [Verb]'s Default, then the [MsgNothing] flash.
//
// Every resolved click logs a player-command line for the chat log. A right
// click drops the held item, or resets the verb to look.
//
// A [Trigger] arms when its scene is entered. Every idle frame the engine
// samples its When condition, and on a rising edge (false to true) queues
// its Do; Once makes it fire at most once per arming. Triggers are not
// sampled while a script runs, but the edge is kept across the busy window,
// and loading a save re-arms them as entering the scene would.
//
// # Items and the inventory
//
// An [Item] is a thing the player can carry, registered in the ItemManager.
// [Give] puts it in the [Inventory] (which remembers it by name) and marks it
// known; [TakeAway] removes it. Item.OnUseWith is an item used on another
// item: with one in hand, clicking another in the [InventoryBar] runs
// whichever of the two pairs them, and OnUseSelf answers an item used on its
// own — the held item clicked again. UseWithFail is the same shape as
// ExitLook and ExitTake: the pair the author never wrote is still content,
// and only the game knows what its characters say about it. The Inventory is
// pure runtime state owned by the Game, changed by actions and by the
// InventoryBar.
//
// # Characters and speakers
//
// A [Character] is someone who stands in scenes: a Label for the UI, a Lore
// paragraph for a dossier screen, Poses, a walking Speed (60 pixels a second
// by default), a SpeechColor, a [Delivery], a Start point and a size used by
// the placeholder figure (humanoid when narrower than tall, quadruped
// otherwise).
//
// Delivery is where a character reaches the player, and Say obeys it: a
// DeliveryBubble character speaks in a bubble and the log, a DeliveryLog one
// in the log alone. A voice with no body — a radio, a tape, a narrator — is a
// [Speaker], registered in the SpeakerManager: a name, a Label, a Delivery and
// a Color, and nothing about walking or being drawn. Both are a [Talker], so
// [Say], [Reply] and a [DialogLine] take either; [Game.Speaker] looks a name
// up in the characters first and the speakers second, and a name found in
// neither speaks in a bubble under its own name.
//
// A [Sprite] is a picture or a sequence of them: a sheet and the frames cut
// from it, in order, at FrameTime seconds each, looping or holding the last.
// No frames means the whole sheet is the one frame, so a still item and a
// walking character are the same type. Character.Poses gives one Sprite per
// [Pose], and the engine picks the pose: PoseWalk while moving, PoseTalk while
// the character's line is up, PoseWait after standing still for a while,
// PoseIdle otherwise — or whatever [SetPose] holds until the character next
// walks or speaks. A pose the art does not cover falls back to PoseIdle, and
// no PoseIdle draws the placeholder. Direction is not a second sheet: a
// Sprite says which way it Faces ([Facing]), the character remembers the way
// it last moved, and the frame is mirrored when the two differ; the facing
// survives a save.
//
// # Dialogs and scripts
//
// A [Dialog] is a conversation: named [DialogNode] values, each a few
// [DialogLine] values spoken one at a time followed by the [DialogChoice]
// values the player picks from. [RunDialog] opens one and holds the script
// until it closes; the [DialogBox] widget shows the lines and the choices;
// a choice runs its Actions in sequence, and [EndDialog] and [GotoNode] close
// the conversation or move it. See [git.teletypegames.org/engines/inkwell/dialog] for the
// details.
//
// A [Script] is a named composite action — an intro, a victory, a transition
// reused from several places. [RunScript] runs one.
//
// # Verbs
//
// A [Verb] is a canonical name ("look"), a Label in every language, and a
// Default action run when a hotspot has no handler for it. New registers the
// four SCUMM verbs (look, use, talk, take) unless the game brought its own,
// and a game adds more:
//
//	g.VerbManager.Register(&inkwell.Verb{
//		Name:    "push",
//		Label:   inkwell.Text{{Lang: "en", Text: "Push"}, {Lang: "hu", Text: "Lökd"}},
//		Default: inkwell.Say(player, inkwell.Text{{Lang: "en", Text: "It won't move."}}),
//	})
//
// A hotspot answers a custom verb through OnVerb. The verb bar and the verb
// wheel read the VerbManager every frame, so a verb added at runtime appears
// at once.
//
// # Stages
//
// A [Stage] is one full-screen thing the game can be showing, and exactly one
// is showing at a time. A main menu is a Stage, a settings screen is a Stage,
// a credits card is a Stage — and so is the game itself (World: true), which
// is the point: the world is not the thing every other screen floats over,
// it is one of the screens. A Scene is a place in the fiction and the engine
// always has one; a Stage is what owns the frame. Off the world Stage the
// scene is not drawn, characters do not walk, triggers do not fire and a
// click reaches no hotspot.
//
// Widgets follow without a new field, because [OnStage] is a [Condition]
// like any other:
//
//	game := StageManager.Register(&inkwell.Stage{Name: "game", World: true})
//	mainMenu := StageManager.Register(&inkwell.Stage{Name: "main_menu"})
//
//	WidgetManager.Register(&inkwell.Menu{
//		Name:   "main_menu",
//		When:   inkwell.OnStage(mainMenu),
//		Sticky: true,
//		Root:   inkwell.MenuItem{ … },
//	})
//
// [GoToStage], [Restart] and [LoadSlot] all land at the top of the next frame
// rather than where they are called, so a menu row can say what it means in
// one sequence:
//
//	Do: inkwell.Seq(inkwell.Restart(), inkwell.GoToStage(game)),
//
// A Stage may name how it gives way: a change out of it goes through its
// [Transition], and nil is a cut. The whole pending change lands on the fully
// covered frame, a Restart or LoadSlot asked for at the same time included.
//
// Hold shows a Stage for a number of seconds and then goes to Next, through
// the Stage's Transition, which is all a splash card or a title sequence is.
// The clock starts once OnEnter is done. Skippable lets Space, Enter or a
// click move on at once; the key or click is claimed, so nothing underneath
// reads it.
//
//	StageManager.Register(&inkwell.Stage{
//		Name:       "splash_studio",
//		Transition: inkwell.Fade(0.5),
//		Hold:       2,
//		Skippable:  true,
//		Next:       splashTitle,
//	})
//
// A game whose StartStage is the world starts a session at boot: the start
// scene is entered and StartScript is queued. A game that opens on a menu
// starts one when the player says so, which is all Restart is.
// [Game.Started] (and the [Started] condition) report whether a session
// exists, which is what a "Resume" or "Save" row is enabled by.
//
// # Achievements
//
// An [Achievement] is something the player is recognised for, declared as
// content: a Name a platform knows it by (the Steam API name), a Title and a
// Text in every language the game speaks, and the [Condition] that earns it.
// It is registered in Game.AchievementManager. The engine evaluates every
// achievement on each frame the world is showing, a running script included,
// and hands the ones that have just come true to every [Unlocker] among the
// extensions. It keeps no record of its own beyond a per-process set; see
// [git.teletypegames.org/engines/inkwell/achievement].
//
// # Actions
//
// An [Action] is an immutable spec of work, safe to keep in a Hotspot field
// and run on every click; its Start makes a [Runner], one in-flight
// execution that is ticked each frame with a [Ctx] and answers a [Status]:
// StatusRunning, StatusDone (a sequence moves on) or StatusFailed (a sequence
// stops here). The built-in constructors:
//
//   - Flow: [Seq], [Par], [If], [Wait], [Noting].
//   - Speech: [Say], [Reply], [Flash], [ShowEnd].
//   - Scenes and characters: [GoTo], [Back], [Walk], [SetPose].
//   - Inventory and knowledge: [Give], [TakeAway], [RequireItem], [Learn].
//   - World state: [SetFlag], [ClearFlag], [SetVar], [SetVarIfEmpty].
//   - Dialogs and scripts: [RunDialog], [EndDialog], [GotoNode], [RunScript].
//   - Screens and sessions: [GoToStage], [GoBack], [Restart], [LoadSlot],
//     [SaveSlot], [Quit].
//   - Settings: [UseTheme], [UseLang], [SetVoice].
//   - Audio: [PlayMusic], [StopMusic], [PlaySound].
//   - Escape hatch: [Custom], which wraps a function as an action.
//
// The Ctx a runner gets carries the frame's DT, the [State], and the Game as
// an any; an adventure action type-asserts it to *Game:
//
//	bumpScore := inkwell.Custom(func(ctx *inkwell.Ctx) inkwell.Status {
//		score, _ := ctx.State.Var("score").(int)
//		ctx.State.SetVar("score", score+10)
//		return inkwell.StatusDone
//	})
//
// [git.teletypegames.org/engines/inkwell/act] documents the adventure actions and the host
// interfaces they reach the Game through;
// [git.teletypegames.org/engines/inkwell/inkbase/action] documents the generic ones and how
// to write an action of your own.
//
// # Conditions
//
// A [Condition] is a stateless predicate evaluated lazily against a Ctx — in
// [If], in DialogChoice.Show, in a widget's When, in a Trigger and in an
// Achievement. [Not], [And] and [Or] combine; [Flag], [VarEq] and [Knows]
// read the State; [HasItem] and [SelectedItem] the inventory; [InScene],
// [OnStage], [OnWorld] and [Started] where the game is. [Given] is a
// condition written as a literal — every field set on it has to hold — for
// the place a condition is read most, a dialog choice:
//
//	Show: inkwell.Given{
//		Item:    mysteryTape,
//		NotFlag: "dex_told",
//	},
//
// See [git.teletypegames.org/engines/inkwell/cond] and [git.teletypegames.org/engines/inkwell/inkbase/condition].
//
// # World state
//
// [State] is the run's mutable record: flags (the workhorse of puzzle state,
// "cupboard_open"), vars of any type (scores, the clock, a note for the top
// bar), how many times each scene was visited, how many times each dialog
// choice or other counter was bumped (Talked and NoteTalked), and which
// [Knowable] entities the player has met. Known is kept for a dossier or
// codex screen: a character or speaker is known once it has said a line, an
// item once it has been in the inventory, and a script says [Learn] for the
// rest. Its value is the order of discovery, so a screen sorted by it reads
// as the player's own history; [Knows] is the same fact as a condition. An
// item and a character may share a name, because the kind is part of the
// key.
//
// # Widgets
//
// The HUD is a set of [Widget] values registered in the WidgetManager. A
// widget handles its own input, hit testing, state and drawing; the library
// imposes no layout. Each declares a [Layer] ([Layered]) — the draw order
// runs from LayerScene up to LayerCursor, ticks run top layer down so the
// widget drawn on top gets the click first, and registration order decides
// inside one layer. A widget claims a click or a key through [Input], and
// widgets below see it as consumed. Every built-in widget has a When
// condition ([Gated]): a widget whose condition is false neither ticks nor
// draws nor blocks a click, so a HUD hides for a cutscene by saying so on its
// widgets. A running script has the frame to itself except for the widgets
// that implement [Awake] — [Menu] and [DialogBox] — which tick before the
// script's runner. Whatever click survives the widgets reaches the scene,
// unless a visible widget claims the point as its own (see below).
//
// Besides the inkbase layers, this package adds LayerDebug, between the
// dialog and the menus, for the [HotspotDebug] overlay. A wrapper widget that
// forwards to an inner one should return [LayerOf] the inner one from its
// own Layer, so that wrapping does not move it.
//
// A widget that occupies an area of the screen implements [ClickBlocker]:
// BlocksClickAt(ctx *UICtx, p Point) bool, reporting whether p is its own.
// A left click there that the widget's Tick did not take never reaches the
// scene (no hotspot, no walk), and the [RadialVerbs] coin does not open over
// it. The built-ins that do are the verb bar, the inventory, the dialog box
// while a dialog is open, the top bar, the character panel while its
// character is in the scene, the chat log, a blocking panel, the hotspot
// panel, an open menu and the radial menu itself while it is visible.
// ClickBlocker is the inkbase widget contract of the same name, so a widget
// from the inkbase widget packages blocks the same way.
//
// [git.teletypegames.org/engines/inkwell/inkbase/widget] documents the contract in full. A
// custom widget reaches the Game by type assertion:
//
//	type Minimap struct {
//		Name   string
//		Bounds inkwell.Rectangle
//	}
//
//	func (m *Minimap) GetName() string        { return m.Name }
//	func (m *Minimap) Tick(ctx *inkwell.UICtx) {}
//	func (m *Minimap) Draw(dst *ebiten.Image, ctx *inkwell.UICtx) {
//		g := ctx.Game.(*inkwell.Game)
//		… // render a scene thumbnail, mark the cast
//	}
//
// # Built-in widgets
//
// Fields left at their zero value fall back to a sensible default. The
// generic ones — [Menu], [Panel], [TextCard], [KeyBinding], [Letterbox],
// [EndCard], [StatusLine], [TopBar] and [ChatLog] — are the inkbase widgets,
// re-exported: the Game is their host, and draws them with the active theme
// (MenuColors, ChatLogColors and the like are the theme's tokens handed to
// them). The rest need the adventure and live here.
//
//   - [Panel]: a coloured rectangle, optionally bordered, as a backdrop; with
//     Blocks it is opaque to clicks.
//   - [StatusLine]: the flash message, or the composed hover hint.
//   - [VerbBar]: the permanent SCUMM verb grid.
//   - [RadialVerbs]: a verb coin that opens where the trigger button is
//     pressed, or a permanent verb wheel.
//   - [InventoryBar]: the item slots; hovering names an item, clicking
//     selects it, and under look it says the description.
//   - [SpeechBubble]: the line being spoken, over the speaker's head, in
//     the speaker's colour.
//   - [DialogBox]: the open conversation's lines and choices; while a dialog
//     is up it takes every click.
//   - [EndCard]: the full-screen card [ShowEnd] puts up.
//   - [Cursor]: the crosshair, or the held item's picture.
//   - [HotspotDebug]: a translucent plate, one colour per [DebugKind], over
//     everything a player could point at — hotspots, exits, the cast,
//     inventory slots and panel targets — toggled by F1. A widget with
//     clickable things of its own joins in by implementing [Targeted].
//   - [TopBar]: the scene title (with a note from a State var), a score and
//     a clock.
//   - [CharacterPanel]: a card with a character's portrait, a role badge and
//     [CharStat] rows, hidden while the character is not in the scene.
//   - [ChatLog]: the message log; a [SpeakerStyle] says how replies are
//     attributed, and Deliveries keeps only the replies delivered certain
//     ways, which is how a radio or a narrator gets a panel of its own.
//   - [KeyBinding]: a key and what it does, so "Escape opens the main menu"
//     is a registration like everything else.
//   - [TextCard]: a block of lines centred on the screen — a title card, a
//     credits page.
//   - [Letterbox]: the cutscene bars, on LayerScene, under the HUD.
//   - [HotspotPanel]: makes a region of the HUD answer the verbs the way a
//     region of the picture does — its hotspots go through
//     [Game.InvokeVerb] — so a device mounted on the panel is Hotspot
//     literals rather than a hand-written switch. It draws nothing; a game
//     embeds it and adds a Draw.
//   - [Menu]: the menu system (below).
//
// [RegisterDefaultUI] installs the SCUMM layout (verb bar and inventory),
// [RegisterRadialVerbUI] a verb coin with a wider inventory, and
// [RegisterRichUI] a top bar, character panels (nil skips one), a chat log
// and a verb wheel.
//
// # Menus, the debug menu and the explorer
//
// [Menu] is a nested list in a window: a page of rows, any of which can open
// a page of its own, as deep as the tree goes. A row and a page are the same
// type, [MenuItem]: Items is a page as authored, Expand a page read out of
// the game while it is open, Body a paragraph to read, and a row with none of
// them is a leaf whose Do is queued on Enter. [FixedValue] makes a Value
// column that does not depend on the game. The menu claims every key it
// uses; Sticky makes it the screen rather than a window over one, which is
// what a Stage is usually made of. Colours come from the active theme's
// dialog tokens. [git.teletypegames.org/engines/inkwell/inkbase/menu] describes the model,
// the keys and the mouse in full; this Menu is the same widget, wired to the
// Game.
//
// [DebugMenuRoot] is the one use of the menu system the library ships: what
// the game is made of, and what is clickable in the current scene. It is
// data, not a widget, and nothing registers it — a game asks for it and gates
// it on its own debug switch:
//
//	g.WidgetManager.Register(&inkwell.Menu{
//		Name:      "debug_menu",
//		Root:      inkwell.DebugMenuRoot,
//		When:      whenDebug,
//		ToggleKey: ebiten.KeyD,
//	})
//
// Its pages:
//
//	DEBUG
//	├── Scenes         every registered scene · Enter teleports
//	├── Items          give · take away · select · use on self · use with → target
//	├── Action points  the current scene's hotspots, exits included, in hit-test order
//	│   └── <hotspot>  the verbs it answers · use with → item
//	├── Exits          the current scene's exits · go · set or clear the flag it needs
//	└── EXPLORER       the Game itself, field by field
//
// The rows name registry keys, not the prose a player reads, since a debug
// menu is read by whoever writes the content. Every verb goes through
// [Game.InvokeVerb], the same path a click takes, so the chat log reads as if
// the player had done it; a use-with row hands the item over first. A game
// that wants sections of its own builds a root around it:
//
//	Root: inkwell.MenuItem{
//		Label: inkwell.DebugMenuRoot.Label,
//		Items: append(inkwell.DebugMenuRoot.Items, mySection),
//	},
//
// [ExplorerRoot], the last page of the debug menu, is the Game itself, field
// by field, worked out by reflection when a page opens: a manager opens onto
// its entities, an entity onto its literal, an action onto the actions it is
// made of, a Text onto its translations. The value column is what a row is
// worth now — a name, the current line of a Text, the truth of a Condition,
// the kind of an Action, a colour, a count. A pointer to a registered entity
// is a link to it. Unexported fields are read too, so the explorer can see
// inside a Say; it stops at functions, where the writing turns into code.
// Every action row opens with "run", so any authored sequence can be played
// from where it is written. [Explore] runs a game with nothing on screen but
// this, validated and wired exactly as Run wires it:
//
//	if explore {
//		return inkwell.Explore(game)
//	}
//	return inkwell.Run(game)
//
// # Themes
//
// A [Theme] is a named set of colour tokens every built-in widget draws with
// — panel, status, verbs, inventory, speech, dialog, overlays, top bar, chat
// log and character panel. New registers the four presets
// ([RegisterPresetThemes]: classic-scumm, sierra-coin, paper-notebook and
// terminal-green) and selects StartTheme, or classic-scumm. A custom theme is
// registered like anything else and selected with [Game.UseTheme], or with
// the [UseTheme] action, which also saves the setting; widgets read the
// active theme every frame, so the change shows on the next one. See
// [git.teletypegames.org/engines/inkwell/theme].
//
// # Input
//
// [Input] (Game.Input) is polled once per frame: the pointer position, the
// two mouse buttons' presses, the wheel's travel and every key that went
// down. Clicks and keys are consume-on-use — whoever consumes first wins,
// and later checks see nothing — which is how one click is authoritative
// across the widget stack and the scene. KeyPressed is edge-triggered, and
// ClaimKey is KeyPressed plus ConsumeKey. [git.teletypegames.org/engines/inkwell/inkbase/input]
// documents the same model for inkbase games.
//
// # The game loop
//
// [Run] (or [Game.Run]) resolves every Scene.Derived, starts the extensions,
// reads the settings and settles the language, validates the game, installs
// the default UI when no widget is registered, shows the StartStage —
// starting a session at once when it is the world — sizes the window to the
// resolution times WindowScale, and enters the Ebitengine loop.
//
// Each frame's update, in order:
//
//	tick the extensions
//	apply the pending screen change          GoToStage, Restart, LoadSlot
//	if quitting: stop
//	count the flash down, poll input, advance the transition,
//	start the next queued action
//	if a transition is covering the picture: stop here
//	on the world Stage: evaluate the achievements
//	if an action is running:
//	    tick the Awake widgets, top layer down
//	    tick the action
//	    on the world Stage: move the characters
//	    stop here
//	on the world Stage: sample the scene's triggers
//	tick every visible widget, top layer down
//	on the world Stage: resolve the click on the scene, move the characters
//
// Each frame's draw: the theme's backdrop, the scene layers at or below 0 by
// depth, the characters sorted by their feet, the layers above 0, every
// visible widget by layer and registration order, and the transition's cover
// on top.
//
// A moving character steps along its waypoints at Speed pixels per second;
// [Walk] holds its script until the character arrives. A click on empty floor
// with WalkOnClick does not queue a Walk: a queued walk would be a script,
// and a script owns the frame, so the next click could not redirect it.
//
// # Scene transitions
//
// A [Transition] is how one picture gives way to another, in two halves: the
// old picture is covered, the change lands on the frame it is fully hidden,
// and the new picture is uncovered. [Fade] and [FadeTo] are the ones the
// library ships; a Stage names its own, and a scene change uses a quarter
// second fade. At the midpoint of a scene change the current scene changes,
// the visit is counted, the actors are placed and the new scene's music
// plays; input is frozen while the old picture is going. A change made by
// [GoTo] or [Back] runs the old scene's OnLeave and the new scene's OnEnter
// inside that action, right after the swap, and the action is done when the
// new picture is fully uncovered, so a script can go on in the new scene:
//
//	inkwell.Seq(
//		inkwell.GoTo(scene.Harbour),
//		inkwell.Say(paul, inkwell.Text{{Lang: "en", Text: "The boat is gone."}}),
//		inkwell.GoTo(scene.Alley),
//		inkwell.Say(paul, inkwell.Text{{Lang: "en", Text: "Back where I started."}}),
//	)
//
// [Game.ChangeScene] queues the two hooks instead. The opening scene is entered
// without a transition. See [git.teletypegames.org/engines/inkwell/inkbase/transition].
//
// # Validation
//
// [Game.Validate] cross-checks what the types cannot. Run calls it after
// every Scene.Derived has been resolved, and a test can call it on its own.
// It checks that StartScene and StartStage are set and registered; that a
// Stage that holds has a registered Next; that every scene actor's character,
// and every exit's target, is registered (a nil target means "back the way
// you came"); that Player, when set, is registered; that every achievement
// has a When; and that the active theme is registered. Most of what such a
// validator would check is a type here: a scene, an item or a dialog is
// referred to by pointer, so an unknown one is a compile error. Duplicate
// names are caught earlier still, by Register's panic.
//
// # Save and load
//
// A Game keeps what has to outlive the process through its [Store]: two
// documents, named by [SettingsKey] and [SlotsKey]. Game.App names the game
// to the machine ([Game.AppName] is App or "inkwell"), and the platform picks
// the place: natively a [FileStore], one JSON file per document under the
// user's configuration directory (~/Library/Application Support/<app> on
// macOS, ~/.config/<app> on Linux, %AppData%\<app> on Windows), written to a
// temporary file and renamed so a crash leaves the previous document; in a
// browser, one localStorage entry per document. A game that wants neither
// hands in its own Store; a test hands in a [MemoryStore]. Read returns
// [ErrNotStored] for a key nothing has written.
//
// [Settings] are read before the first frame: the stored language is used
// when the run has not chosen one, the stored theme when it is registered,
// the stored window scale when the game named none — so a -lang flag wins for
// one run without being written back. UseLang, UseTheme and SetVoice save
// the settings at once. Extra is the game's own: read and write it through
// [Game.Settings], call [Game.SaveSettings], and it round-trips beside the
// engine's fields.
//
// [Game.Save] and [Game.Load] write and read one slot of the slots document;
// [SaveSlot] and [LoadSlot] are the actions. [Game.Slots] reads the header of
// slots 1 to n as [SlotInfo] values for a load or save screen; a missing or
// unreadable slot comes back as one that does not exist. Each slot keeps
// when and where it was saved in the clear, beside the session itself as a
// base64 JSON string, so listing slots decodes no session. A session holds
// the current and previous scene, the selected verb and the active theme,
// every character's position, target, movement and facing, the whole State,
// and the inventory with its selection. It does not hold the registrations
// — the game rebuilds those on every launch — nor the running action or
// open dialog (Load cancels both), the poses (re-derived on the next frame)
// or walk paths (a character still moving finishes in a straight line). A
// save is checked against the current registries before anything is
// touched, so a corrupt or outdated save returns an error without changing
// the running game, and a save from a newer format version is refused. Vars
// are stored as JSON, so numbers come back as float64.
//
// # Errors
//
// The sentinel errors — [ErrNoStartScene], [ErrNoStartStage],
// [ErrUnknownScene], [ErrUnknownStage], [ErrUnknownCharacter],
// [ErrUnknownAsset], [ErrUnknownScript] and [ErrNotStored] — are wrapped with
// context and checked with errors.Is. Manager.Register and MustGet panic
// instead of returning an error, because they signal construction-time bugs
// that should crash loudly in development.
//
// # Testing and debugging
//
// A game's cheapest CI gate is a test that builds the game and validates it;
// it opens no window and finishes in milliseconds:
//
//	func TestBuildValidates(t *testing.T) {
//		g := domain.Build()
//		if err := g.Validate(); err != nil {
//			t.Fatalf("validate: %v", err)
//		}
//	}
//
// Runtime debug switches:
//
//	inkwell.DebugLog = true                                            // engine log to stderr
//	&inkwell.HotspotDebug{Enabled: true}                               // start with the F1 overlay on
//	&inkwell.Menu{Root: inkwell.DebugMenuRoot, ToggleKey: ebiten.KeyD} // the debug menu
//
// # Project layout
//
// The adventure layer is this package plus a few sub-packages it re-exports,
// so a game imports this package alone:
//
//   - [git.teletypegames.org/engines/inkwell/act]: the adventure actions
//     (Say, Walk, Give, RunDialog, …).
//   - [git.teletypegames.org/engines/inkwell/cond]: the adventure conditions
//     (HasItem, InScene, OnStage, Given, …).
//   - [git.teletypegames.org/engines/inkwell/dialog]: dialog trees.
//   - [git.teletypegames.org/engines/inkwell/item]: items.
//   - [git.teletypegames.org/engines/inkwell/sprite]: sprites, poses and
//     facing.
//   - [git.teletypegames.org/engines/inkwell/verb]: verbs and the SCUMM set.
//   - [git.teletypegames.org/engines/inkwell/theme]: themes and the presets.
//   - [git.teletypegames.org/engines/inkwell/achievement]: achievements.
//   - [git.teletypegames.org/engines/inkwell/steam]: the Steamworks extension,
//     which a game imports itself.
//
// The generic foundation is [git.teletypegames.org/engines/inkwell/inkbase]
// and its sub-packages: font, input, asset and transition (Ebitengine
// services), stage (the Stage), action and condition (the generic actions
// and conditions), widget (the widget contract) and the generic widgets —
// menu, panel, textcard, keybinding, letterbox, endcard, statusline, topbar
// and chatlog. Those widgets are the ones this package re-exports (Menu,
// ChatLog, TopBar, …): *Game implements each widget package's Host, drawing
// them with the active theme, so the same widget serves an adventure and a
// game built on inkbase alone.
//
// The library is MIT licensed; see LICENSE.md.
package inkwell
