# inkwell — reference manual

A point-and-click adventure game framework for Go, built on top of
[Ebitengine](https://ebitengine.org/). Games are composed by registering
plain struct literals into per-entity `*Manager` registries hanging off a
single `*Game` root. The library handles input, rendering, dialog trees,
cutscenes, HUD, themes and lazy asset loading; the domain only declares
content.

This document is the full library reference. The companion `inkwell-demo`
repo at <ssh://git@git.teletypegames.org:2222/games/inkwell-demo> is a small
end-to-end example ("Morning Coffee") that consumes this library as a
plain Go module.

**Module path:** `git.teletypegames.org/games/inkwell`

```go
import "git.teletypegames.org/games/inkwell"
```

---

## Table of contents

1. [Introduction](#1-introduction)
2. [Quick start](#2-quick-start)
3. [The Manager pattern](#3-the-manager-pattern)
4. [The Game aggregate](#4-the-game-aggregate)
5. [Text and languages](#5-text-and-languages)
   - 5.1 [Writing one](#51-writing-one)
   - 5.2 [Resolving one](#52-resolving-one)
   - 5.3 [Composing](#53-composing)
   - 5.4 [What the engine says itself](#54-what-the-engine-says-itself)
   - 5.5 [Language-stable state](#55-language-stable-state)
   - 5.6 [The font](#56-the-font)
   - 5.7 [A ladder of lines](#57-a-ladder-of-lines)
6. [Geometry primitives](#6-geometry-primitives)
7. [Entity reference](#7-entity-reference)
   - 7.1 [Asset](#71-asset)
   - 7.2 [Scene](#72-scene)
   - 7.3 [Exit](#73-exit)
   - 7.4 [Hotspot](#74-hotspot)
   - 7.5 [Trigger](#75-trigger)
   - 7.6 [Item](#76-item)
   - 7.7 [Inventory](#77-inventory)
   - 7.8 [Character](#78-character)
   - 7.9 [Dialog](#79-dialog)
   - 7.10 [Script](#710-script)
   - 7.11 [Verb](#711-verb)
   - 7.12 [Stage](#712-stage)
8. [The action system](#8-the-action-system)
   - 8.1 [Action and Runner](#81-action-and-runner)
   - 8.2 [Status and Ctx](#82-status-and-ctx)
   - 8.3 [Built-in actions](#83-built-in-actions)
   - 8.4 [Custom actions](#84-custom-actions)
9. [Conditions](#9-conditions)
10. [World state](#10-world-state)
11. [The widget system](#11-the-widget-system)
    - 11.1 [Widget interface](#111-widget-interface)
    - 11.2 [Z-order and input consumption](#112-z-order-and-input-consumption)
    - 11.3 [Built-in widgets](#113-built-in-widgets)
    - 10.4 [The clickBlocker contract](#114-the-clickblocker-contract)
    - 11.5 [Registration helpers](#115-registration-helpers)
    - 11.6 [Custom widgets](#116-custom-widgets)
    - 11.7 [The menu system](#117-the-menu-system)
    - 11.8 [The debug menu](#118-the-debug-menu)
12. [Themes](#12-themes)
13. [Asset pipeline](#13-asset-pipeline)
14. [Audio](#14-audio)
15. [Input](#15-input)
16. [The game loop](#16-the-game-loop)
17. [Scene transitions](#17-scene-transitions)
18. [Validation](#18-validation)
19. [Save and load](#19-save-and-load)
20. [Errors](#20-errors)
21. [Testing](#21-testing)
22. [Project layout](#22-project-layout)
23. [License](#23-license)

---

## 1. Introduction

`inkwell` is **not** a generic engine — it is a thin Go-level DSL layer on
top of Ebitengine that implements the fixed skeleton of the point-and-click
adventure genre as it existed in the LucasArts SCUMM era (Maniac Mansion,
Monkey Island, Day of the Tentacle):

- scenes with hotspots and click zones,
- a verb-based interaction grammar (Look, Use, Talk, Take, …),
- inventory, item-on-item and item-on-hotspot combinations,
- dialog trees with conditional choices,
- a scripting / cutscene system built from composable `Action` values,
- a widget-tree HUD that you can extend with anything Ebitengine can draw,
- swappable color themes,
- placeholder graphics so the game runs without any art on disk.

The design collapses to **one organising pattern**: every named entity is
registered via `Manager.Register(Entity{Name: "..."})`. No builders, no
fluent chains, no scattered registries. The same mental model applies to
items, scenes, characters, dialogs, scripts, assets, verbs, widgets and
themes alike.

A complete example game ("Morning Coffee") lives in the separate
[`inkwell-demo`](ssh://git@git.teletypegames.org:2222/games/inkwell-demo) repo
and consumes this library as a regular Go module.

---

## 2. Quick start

**Prerequisites.** Go 1.24+ (for generic type aliases) and a working OpenGL
context. Ebitengine handles windowing and the render loop.

**Use the library in your own project:**

```bash
go mod init my-game
go get git.teletypegames.org/games/inkwell
```

If the private host can't be reached over HTTPS, point Go at the SSH
endpoint:

```bash
git config --global \
    url."ssh://git@git.teletypegames.org:2222/".insteadOf "https://git.teletypegames.org/"
export GOPRIVATE=git.teletypegames.org/*
```

**Minimum game:**

```go
package main

import (
    "log"
    "git.teletypegames.org/games/inkwell"
)

func main() {
    player := &inkwell.Character{Name: "player", W: 28, H: 62}
    room := &inkwell.Scene{
        Name:       "room",
        Background: &inkwell.Asset{Name: "bg/room", Path: "assets/bg/room.png"},
        Actors:     []inkwell.SceneActor{{Character: player, At: inkwell.Point{X: 160, Y: 140}}},
        Hotspots: []inkwell.Hotspot{{
            Name:   "door",
            Area:   inkwell.Rect(260, 50, 40, 90),
            Label:  inkwell.Text{{Lang: "en", Text: "door"}},
            OnLook: inkwell.Say(player, inkwell.Text{
                {Lang: "en", Text: "A wooden door."},
            }),
        }},
    }
    world := &inkwell.Stage{Name: "game", World: true}

    g := inkwell.New(inkwell.Game{
        Title:      "Sample",
        Width:      320,
        Height:     200,
        StartStage: world,
        StartScene: room,
        Player:     player,
    })
    g.CharacterManager.Register(player)
    g.SceneManager.Register(room)
    g.StageManager.Register(world)

    inkwell.RegisterDefaultUI(g) // optional; auto-called by Run if no widgets registered
    if err := inkwell.Run(g); err != nil {
        log.Fatal(err)
    }
}
```

---

## 3. The Manager pattern

Every name-addressable entity goes through one shared generic type:

```go
// inkwell/core.manager.go

type Named interface {
    GetName() string
}

type Manager[T Named] struct { /* ... */ }
```

Each entity kind has its own alias (generic type alias, Go 1.24+):

```go
type ItemManager       = Manager[*Item]
type SceneManager      = Manager[*Scene]
type CharacterManager  = Manager[*Character]
type SpeakerManager    = Manager[*Speaker]
type DialogManager     = Manager[*Dialog]
type ScriptManager     = Manager[*Script]
type AssetManager      = Manager[*Asset]
type VerbManager       = Manager[*Verb]
type WidgetManager     = Manager[Widget]
type ThemeManager      = Manager[*Theme]
```

### 3.1 Public API

| Method                 | Behaviour                                                  |
|------------------------|------------------------------------------------------------|
| `Register(v T) T`      | Adds `v` to the registry and hands it back. Panics on empty or duplicate `Name`. |
| `Set(v T)`             | Replaces an entry in place, keeping its position; registers it when the name is new. |
| `Get(name) (T, bool)`  | Looks up by name. The boolean is `false` if absent.        |
| `MustGet(name) T`      | Same as `Get`, but panics on missing names.                |
| `Has(name) bool`       | True if the name is registered.                            |
| `Len() int`            | Number of registered entries.                              |
| `Names() []string`     | Returns names in **insertion order** (used for widget Z-order). |
| `All() []T`            | Returns every entry in **insertion order**.                |
| `Each(fn func(T))`     | Iterates in insertion order.                               |

`Register` returns what it registered so that a registration can be an
**expression** rather than a statement, which is what lets an entity file be a
declaration instead of a procedure:

```go
// scene.alley.go — the whole file
package scene

var Alley = SceneManager.Register(&Scene{
	Name: "alley",
	…
})
```

That is not merely shorter than wrapping the call in an `init()`. A
package-level variable is initialised in **dependency order**, so the manager
is built before any file that writes into it, whichever one the compiler reads
first; the `init()` form rests instead on the separate rule that every package
variable is ready before any `init()` runs. Registration order is unchanged
either way — the order the files are presented in.

The result is the handle. Every manager holds **pointers** — `Manager[*Scene]`,
not `Manager[Scene]` — so the variable a file binds and the entry the engine
looks up are one instance, and every field that refers to another entity is a
pointer to it: `Stage.Background *Asset`, `Exit.To *Scene`,
`SceneActor.Character *Character`, `Say(speaker Talker, …)`. A reference is
therefore an identifier, checked by the compiler and followed by an editor, and
`Name` is left with the one job only a string can do: it is the key a save
file, a `State` var and a command-line flag use to find the entity again.

The result is also safe to ignore, so the plain statement form still compiles,
and a file that has to *build* something before registering it (a loop, or a
value wired to itself) keeps its `init()` and should.

### 3.2 The `Named` contract

Every registerable struct implements `GetName() string`. Most also
implement an optional `TypeLabel() string` that the manager uses in panic
messages for clearer errors:

```go
func (i Item) GetName() string   { return i.Name }
func (i Item) TypeLabel() string { return "item" }
```

### 3.3 Why panic on duplicates?

A duplicate registration is a construction-time bug — code that ran once
on startup. Crashing loudly during `Build()` surfaces the problem in
development; quietly accepting the second registration would silently mask
shadowed entities at runtime.

When a registration genuinely has to be rewritten — a derived field filled in
after the fact, a hot-reloaded entity — `Set` is the deliberate overwrite. It
keeps the entry's place in the insertion order, so widget Z-order and any other
order-sensitive iteration survive the rewrite.

---

## 4. The Game aggregate

```go
// inkwell/core.game.go

type Game struct {
    Title         string
    Width, Height int

    // entity managers
    ItemManager      *ItemManager
    SceneManager     *SceneManager
    CharacterManager *CharacterManager
    SpeakerManager   *SpeakerManager
    DialogManager    *DialogManager
    ScriptManager    *ScriptManager
    AssetManager     *AssetManager
    VerbManager      *VerbManager
    WidgetManager    *WidgetManager
    ThemeManager     *ThemeManager
    StageManager     *StageManager

    // exits
    SceneRect Rectangle          // the part of the window the picture occupies
    ExitLook  func(Exit) Action  // default look response for generated exits
    ExitTake  func(Exit) Action  // default take response for generated exits

    // domain defaults
    UseWithFail func(item *Item, h *Hotspot) Action   // unauthored use-with pair
    Player      *Character                            // actor a scene with none gets; who WalkOnClick moves
    WalkOnClick bool                                  // a click on empty floor walks Player (default off)
    Walkboxes   []Polygon                             // floor a scene with none walks
    WindowScale int                                   // window = resolution × this (0 = 4)

    // what the game boots as
    StartStage   *Stage    // the screen the game opens on
    StartScene   *Scene    // the scene a session opens in
    StartScript  *Script   // queued after that scene's OnEnter, once
    FinaleScript *Script   // closing script, queued straight after StartScript
    StartTheme   *Theme    // theme selected at boot (nil -> classic-scumm)
    Flags        []string  // State flags set before the first frame

    // runtime services
    State     *State
    Inventory *Inventory
    Audio     *AudioPlayer
    Input     *Input

    // log buffer
    MaxLogLines int

    // where settings and saves live
    App   string   // the game's name to the machine: its config directory, its localStorage prefix
    Store Store    // nil takes the platform's default for App
    /* unexported runtime fields */
}
```

### 4.1 Construction

```go
func New(spec Game) *Game
func NewGame(title string, w, h int) *Game
```

`New` builds a game from a declared one. Everything a game *is* rather than
*does* is a field, so the whole of it is one literal at the call site instead
of a page of assignments the reader has to hold in their head:

```go
game := inkwell.New(inkwell.Game{
    Title:       "Real World",
    Width:       640,
    Height:      380,
    WindowScale: 2,

    StartScene:  "alley",
    StartScript: "opening",
    StartTheme:  "realworld-93",

    Player:      "paul",
    Walkboxes:   scene.Floor,
    UseWithFail: script.UseWithFail,

    SceneManager:  scene.SceneManager,
    AssetManager:  background.BackgroundManager,
    ScriptManager: script.ScriptManager,
    WidgetManager: widget.WidgetManager,
    ThemeManager:  theme.ThemeManager,
})
```

What the spec leaves out, `New` fills in: an empty manager for every category
the caller did not hand over, the runtime halves (`State`, `Inventory`,
`Audio`, `Input`), the SCUMM verb set (`look`, `use`, `talk`, `take`) when no
verbs were registered, and the four preset themes. Widgets are **not**
auto-registered — call `RegisterDefaultUI(g)` (or one of its siblings)
explicitly, or let `Run` install the default set if `WidgetManager` is empty.

A manager handed over here is used **in place, not copied**, so the engine and
the domain share one registry per category instead of two to keep in step —
and handing one over no longer costs anything. `ThemeManager` used to arrive
holding the four presets, so replacing it threw them away; `New` now *adds*
the presets to whichever `ThemeManager` it ends up with, leaving alone any name
the domain has already taken. `VerbManager` works the same way: hand over a
full one and the SCUMM four are not installed over it.

`NewGame` is `New` for a game that would rather be built in steps — it declares
nothing but the window, and the caller assigns the rest:

```go
g := inkwell.NewGame("Real World", 640, 380)
g.SceneManager = mySceneManager
g.AssetManager = myAssetManager
```

Do it before `Run` — that is where the engine wires up the parts that hold a
registry directly, so whichever managers the `*Game` carries by then are the
ones that get used.

### 4.2 Lifecycle hooks

```go
func (g *Game) StartAt(scene *Scene) *Game    // entry scene
func (g *Game) OnStart(script *Script) *Game  // script run after the scene's OnEnter
func (g *Game) OnFinale(script *Script) *Game // closing script, queued after OnStart
func (g *Game) Validate() error               // cross-check what the types cannot
func (g *Game) Run() error                    // same as inkwell.Run(g)
```

`OnStart` and `OnFinale` name registered scripts rather than taking an action,
so the opening a game plays is content like any other — a `Script` in the
`ScriptManager`, editable without touching the wiring. Both are pointers, so a
script that does not exist cannot be named.

The first three are setters for `StartScene`, `StartScript` and
`FinaleScript`, for a caller building a game in steps rather than declaring
one; they return `*Game` so they chain. A game that hands `New` a literal
names the fields instead and needs none of them:

```go
game := inkwell.New(inkwell.Game{
    …
    StartScene:   start,
    StartScript:  ScriptOpening,
    FinaleScript: finaleScript,   // "" for an ordinary run
})
```

`FinaleScript` is queued directly after `StartScript`, which is what a "boot
straight into the ending" debug flag wants.

### 4.3 Theme accessors

```go
func (g *Game) Theme() *Theme                 // the active theme
func (g *Game) UseTheme(theme *Theme)         // switch theme (runtime ok)
```

Switching is `g.UseTheme(paperNotebook)` with the registered `*Theme`, and the
next frame uses the new colors; the active theme is remembered by its `Name`,
which is what a save file stores.

### 4.4 UI runtime helpers

The current `selectedVerb`, `hoverLabel`, speech and flash message live on
`Game`. They are written by actions / the engine, read by widgets:

```go
func (g *Game) SelectedVerb() string
func (g *Game) SetSelectedVerb(s string)

func (g *Game) HoverLabel() string            // raw target label
func (g *Game) SetHoverLabel(s string)
func (g *Game) HoverHint() string             // composed "verb + target"

func (g *Game) SetSpeech(speaker, text string)
func (g *Game) ClearSpeech()

func (g *Game) FlashLine(text Text)           // ~2s status banner
```

### 4.5 Chat-log API

The engine and the `Say` action push entries here; the `ChatLog` widget
reads them. Set `Game.MaxLogLines` to bound the ring buffer (0 = unbounded).

```go
type LogKind int
const (
    LogAction   // "> look at door" lines (player commands)
    LogResponse // in-world replies / Say output
    LogSystem   // game/system meta
)

type LogMessage struct {
    Speaker string
    Text    string
    Kind    LogKind
}

func (g *Game) LogAction(text string)
func (g *Game) LogResponse(speaker, text string)
func (g *Game) LogSystem(text string)
func (g *Game) Messages() []LogMessage
```

### 4.6 Scene helpers

```go
func (g *Game) CurrentScene() string
func (g *Game) PreviousScene() string
func (g *Game) SceneArea() Rectangle
func (g *Game) SceneHotspots(name string) []Hotspot
func (g *Game) HotspotAt(p Point) *Hotspot
func (g *Game) InvokeVerb(h *Hotspot, verb string)
func (g *Game) CharacterInScene(character *Character) bool

func (g *Game) CurrentStage() *Stage      // the screen showing
func (g *Game) PreviousStage() *Stage     // the one before it
func (g *Game) OnWorldStage() bool        // is the game itself showing?
func (g *Game) Started() bool             // does a session exist?
func (g *Game) GoToStage(stage *Stage)    // show another screen, next frame
func (g *Game) Restart()                  // begin a new session, next frame
func (g *Game) LoadLater(slot int)        // restore a session, next frame
func (g *Game) Quit()                     // close the window
```

`CurrentScene` is where the player is, empty before the first scene is
entered; `PreviousScene` is the one before it, which is where
[`Back()`](#83-built-in-actions) leads. Both survive a save.

`SceneArea` is `SceneRect`, or the whole window when it was never set. It is
what [exit strips](#73-exit) are measured against.

`SceneHotspots` is everything clickable in a scene: its own `Hotspots` first,
then one per `Exit`. Authored hotspots come first because the engine takes the
first area that contains the click, so a painted thing beats the edge strip an
exit sits on wherever the two overlap. The expansion is cached per scene.

`InvokeVerb` is the one path from "the player picked a verb and a target" to
an action on the queue: the hotspot's own `OnUseWith` first, then the item's,
then `OnUseWithAny`, then `UseWithFail` — and with no item selected, the
`On<Verb>` handler, the verb's own `Default`, or the nothing-happens flash.
It is exported because a target does not have to be in the picture. A HUD
widget that owns something clickable hands its own `Hotspot` here and gets
the same logging and the same resolution as a painted one, instead of
re-deciding all of it in its `Tick`.

`CharacterInScene` is used by `CharacterPanel` to auto-hide when its character
isn't an actor in the current scene.

### 4.7 Text drawing and measuring

```go
func (g *Game) DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)

const GlyphW, GlyphH = 6, 16              // the built-in font's cell

func TextWidth(s string) int               // pixel width, rune-counted
func WrapText(s string, maxPx int) []string // word wrap; \n breaks where it stands
func ClipText(s string, maxPx int) string   // cut to width, ellipsis when cut
```

`DrawText` renders in the colour it is given, glyph by glyph out of the
engine's own atlas — each cell is a white glyph over a black drop
shadow, so scaling the blit by the colour tints the glyph and leaves the
shadow black. Every `Theme` text colour is therefore live, in the
built-in widgets and in a domain's own. What the atlas covers, and how
to widen it, is [§5.6](#56-the-font).

These four take a resolved `string`, not a `Text` — they are the sink,
and a `Text` reaches them through `String()`.

The three measuring helpers are what the built-in widgets lay out with,
and every vertical measurement in a HUD should derive from `GlyphH`
rather than a number that happens to look right.

### 4.8 The action queue

```go
func (g *Game) Do(a Action)
```

Queues an action. Queued actions run **one at a time and in order** — the
engine starts the next one on the first frame the previous is done, so a
widget, a trigger and a hotspot click can all hand work in without any of
them being dropped or racing the others. The engine's own hotspot and
`OnEnter` actions go through the same queue.

---

## 5. Text and languages

Every string the player reads is a `Text` — the same string in as many
languages as somebody has written it in.

```go
// inkwell/i18n.text.go

type Translation struct {
    Lang  string
    Text  string
    Voice *Voice   // the recording of this line in this language, if any
}

type Text []Translation

var Lang     string        // the language the game speaks
var Fallback string        // used when Lang has no translation

func (t Text) String() string            // resolve against Lang
func (t Text) In(lang string) string      // resolve against a named language
func (t Text) Empty() bool
func (t Text) Map(fn func(string) string) Text
func Format(format Text, args ...any) Text
```

The wording stays where the thing it names is. A hotspot's label is in
the hotspot; a line of dialog is in the scene that speaks it. Nothing
is lifted out into a catalogue keyed by identifier, because an
identifier is one more name to invent and one more thing to keep in
step — and because a translator reading a scene file can see what the
line is answering.

### 5.1 Writing one

One translation per line, stacked, so the languages read as a table:

```go
Title: inkwell.Text{
    {Lang: "en", Text: "ALLEY — BEHIND NOODLE'S HOUSE"},
    {Lang: "hu", Text: "SIKÁTOR — NOODLE HÁZA MÖGÖTT"},
},
```

```go
inkwell.Say(paul, inkwell.Text{
    {Lang: "en", Text: "Locked. And I'm not guessing four digits."},
    {Lang: "hu", Text: "Zárva. És nem fogok négy számjegyet találgatni."},
})
```

Adding a language is adding a line, and it touches nothing around it.
The zero value says nothing in any language, so a field nobody has
translated behaves exactly like the `""` it used to be, and a `Text`
with one language behaves like a plain string.

Inside a `[]Text` or a `map[string]Text` the element type is implied,
which is what `gofmt -s` leaves behind:

```go
Labels: map[string]inkwell.Text{
    "look": {
        {Lang: "en", Text: "Look"},
        {Lang: "hu", Text: "Néz"},
    },
},
```

### 5.2 Resolving one

`String()` reads the package-level `Lang`, falls back to `Fallback`, and
falls back again to whichever translation was written first — so a line
nobody has translated yet still reaches the player rather than
vanishing. Set both once, before `Run`:

```go
inkwell.Fallback = "en"
inkwell.Lang = *langFlag
```

They are package variables rather than fields on `Game` because a `Text`
resolves itself wherever it is read: inside a widget's `Draw`, inside
the `Say` runner, inside a hover label, inside the hotspot an `Exit`
generates. Threading a `*Game` through all of those would put a
parameter on every one of them for a value there is only ever one of.
`AssetFS` is a package variable for the same reason.

Resolution happens at the moment of display, never at construction, so
content built in an `init()` is unaffected by when the language is set.

### 5.3 Composing

A sentence built around another `Text` has to come out in one language
throughout, which is what `Format` is for: the frame and every `Text`
argument resolve in the same language, and anything that is not a `Text`
is printed as it stands.

```go
var exitLook = inkwell.Text{
    {Lang: "en", Text: "That way: %s."},
    {Lang: "hu", Text: "Arra van: %s."},
}

g.ExitLook = func(e inkwell.Exit) inkwell.Action {
    return inkwell.Say(paul, inkwell.Format(exitLook, e.Label))
}
```

This is not a join, and that is the point: English puts the preposition
between the verb and its object, Hungarian puts it after them. A frame
per language is the only thing that survives contact with a second one.

`Map` is the other composer — it rewrites every translation and keeps
the languages, which is what a derived label wants:

```go
pin := scene.Title.Map(func(s string) string {
    name, _, _ := strings.Cut(s, " — ")
    return name
})
```

### 5.4 What the engine says itself

The handful of sentences inkwell puts on screen without a game having
authored them live in `i18n.strings.go`, in English and Hungarian. They
are variables, so a game that speaks a language the engine has not been
taught can teach it one:

```go
inkwell.MsgNothing = append(inkwell.MsgNothing, inkwell.Translation{
    Lang: "de",
    Text: "Nichts von Bedeutung.",
})
```

| Variable      | What it answers                                  |
| ------------- | ------------------------------------------------ |
| `MsgNothing`  | a verb nothing is bound to                       |
| `MsgNotThat`  | an unpaired use-with, when `UseWithFail` is nil  |
| `MsgNeedItem` | `RequireItem` on an item the player hasn't got   |
| `HintVerb`    | the status line over a plain target              |
| `HintVerbOn`  | the status line while an item is held            |
| `ScoreLabel`  | the `TopBar` score prefix                        |

The four preset verbs carry both languages too, and a game that wants a
third replaces them through `VerbManager.Set`.

### 5.5 Language-stable state

Two things that used to be keyed by wording are not any more, because
wording is no longer a stable identity.

A `Once` dialog choice is recorded against its **position** in its
node, not its text — a choice the player has used up stays used up when
the game is played in another language, and a translator can reword a
line without reopening somebody's save.

`TopBar.NoteVar` reads a `State` var, and a `Condition` compares against
that value, so the value itself has to be identical in every language.
`TopBar.Notes` maps it to what the player reads:

```go
&inkwell.TopBar{
    NoteVar: "note",
    Notes: map[string]inkwell.Text{
        "paused": {
            {Lang: "en", Text: "paused"},
            {Lang: "hu", Text: "szünet"},
        },
    },
}
```

### 5.6 The font

`Text` is only half of it: a game that speaks a second language needs a
font that can draw it. The engine used to render through
`ebitenutil.DebugPrintAt`, whose atlas is 256 cells — U+0000 to U+00FF,
Latin-1 and no further. That is enough for English and stops one
character short of Hungarian: `ő` and `ű` are in Latin Extended-A, and
so is most of what Polish, Czech, Croatian and Turkish need. Anything
past U+00FF drew as a blank gap, em dashes and ellipses included.

inkwell now ships its own atlas, `font.png`, baked by
`internal/fontgen` out of `bitmapfont/v4`:

| Range           | What it is                                   |
| --------------- | -------------------------------------------- |
| U+0020–U+007E   | ASCII                                        |
| U+00A0–U+00FF   | Latin-1 supplement                           |
| U+0100–U+017F   | Latin Extended-A                             |
| U+2010–U+2027   | dashes, curly quotes, dagger, bullet, ellipsis |
| U+20AC          | €                                            |
| U+2190–U+2193   | arrows                                       |

348 glyphs, under 8 KB — the Latin part of a font whose full face is
2.5 MB of embedded CJK, which is the whole of what a European point &
click ever puts on screen. A rune outside the coverage draws a hollow
box, so a gap is visible rather than silent.

`GlyphW` and `GlyphH` are unchanged at 6×16, every cell carries the same
baked drop shadow the debug font had, and the run still starts one pixel
right of `x` — so every layout measured against the old renderer still
measures the same.

To widen the coverage, edit `ranges` in `internal/fontgen/main.go` and
run it; it regenerates both `font.png` and `i18n.glyphs.go`, and refuses
any rune the fixed 6-pixel cell cannot hold.

```
go -C internal/fontgen run .
```

`internal/fontgen` is its own module, so the engine's own dependency
list stays at one entry.

### 5.7 A ladder of lines

```go
// inkwell/i18n.ladder.go

type Ladder struct {
    Rungs []Text     // one line per attempt, in order
    Last  Text       // stands once the rungs run out; empty picks at random
}

func (l Ladder) At(n int) Text
```

A `Ladder` is what somebody says the n-th time the player does the same
fruitless thing. Answering a repeated failure with one line is the surest
way to sound like a machine; a `Ladder` answers the first attempt with the
first rung and keeps a character's patience visibly running out:

```go
var refuse = inkwell.Ladder{
    Rungs: []inkwell.Text{
        {
            {Lang: "en", Text: "No. That's not going to work."},
            {Lang: "hu", Text: "Nem. Ez így nem fog menni."},
        },
        {
            {Lang: "en", Text: "Tried that. It didn't improve."},
            {Lang: "hu", Text: "Ezt már próbáltam. Nem lett jobb."},
        },
    },
    Last: inkwell.Text{
        {Lang: "en", Text: "Leave it. It is not that kind of puzzle."},
        {Lang: "hu", Text: "Hagyd. Nem ilyen fejtörő."},
    },
}

inkwell.Say(paul, refuse.At(state.Talked(key)))
```

`At` counts the first attempt as 1. Past the last rung a `Ladder` picks one
at random, which keeps a stock refusal from turning into a stuck record —
unless it names a `Last`, which stands from then on, because a character who
has given up has given up.

The counter is the game's: `State.NoteTalked` / `State.Talked`
([§10](#10-world-state)) already count how many times a thing has been
worried at, keyed however the game wants — per item, per hotspot, per pair.
`Game.UseWithFail` is the usual caller.

---

## 6. Geometry primitives

```go
// inkwell/util.geometry.go

type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point
func (p Point) Dist(q Point) float64

type Shape interface {
    Contains(p Point) bool
    Bounds() Rectangle
}

type Rectangle struct{ X, Y, W, H float64 }

func Rect(x, y, w, h float64) Rectangle
func (r Rectangle) Contains(p Point) bool
func (r Rectangle) Bounds() Rectangle

type Polygon struct{ Points []Point }

func Poly(pts ...Point) Polygon
func PolyRect(x, y, w, h float64) Polygon     // the rectangle as four corners
func (g Polygon) Contains(pt Point) bool      // ray-casting
func (g Polygon) Bounds() Rectangle           // axis-aligned bbox
```

`PolyRect` is there because a walkbox is a `Polygon` so that a floor can be
any shape, but most floors start as a box — and spelling one out corner by
corner says nothing a reader did not already get from the numbers.

Both `Rectangle` and `Polygon` satisfy `Shape` and can be used as a
`Hotspot.Area`. Internally, hotspots are hit-tested against `Shape.Contains`
each frame.

---

## 7. Entity reference

### 7.1 Asset

```go
// inkwell/asset.def.go

type Asset struct {
    Name string  // logical id used everywhere ("bg/kitchen")
    Path string  // disk path ("assets/bg/kitchen.png")
}

// inkwell/audio.def.go — sound is not an Asset

type Audible interface{ Source() string }   // Voice, Music and Sound

type Voice struct{ Path string }              // a recorded line, on a Translation
type Music struct{ Name, Path string; Loop bool } // a track a Scene or Stage plays
type Sound struct{ Name, Path string }        // a noise, played once by PlaySound
```

Assets are **lazy**: `Register` only stores the spec. The file is opened on
the first `loadedAssets.image(...)` request. Missing files fall back to a
deterministic colored placeholder so the game runs without any art —
useful for prototyping or CI. An `Asset` is a picture; what is heard is a
`Voice`, a `Music` or a `Sound` — see [§14](#14-audio).

```go
kitchen := g.AssetManager.Register(&inkwell.Asset{
    Name: "bg/kitchen",
    Path: "assets/bg/kitchen.png",
})
```

### 7.2 Scene

```go
// inkwell/scene.def.go

type Scene struct {
    Name       string
    Title      Text           // optional human-readable name (TopBar widget reads this)
    Background *Asset         // the picture; nil leaves the theme's backdrop
    Music      *Music         // starts on enter; nil changes nothing
    Hotspots   []Hotspot
    Exits      []Exit         // connections to other scenes
    Walkboxes  []Polygon
    Triggers   []Trigger
    Actors     []SceneActor
    OnEnter    Action
    OnLeave    Action
    Derived    func(g *Game, self *Scene) []Exit  // exits worked out from the rest of the game
}

func (s *Scene) Leads(exits ...Exit) *Scene  // append exits after declaration; see below

type SceneActor struct {
    Character *Character
    At        Point
}
```

`OnEnter` is queued the moment the scene becomes current (after a fade-in);
`OnLeave` runs on the transition out. `Actors` lists which registered
characters are placed in the scene and where their feet start. A `nil`
`Actors` means "the usual cast": `Game.Player` is placed at that
character's `Start`. An empty non-nil slice means an empty scene, which is
how a map or a menu scene opts out.

`Walkboxes` are the same bargain: `nil` falls back to `Game.Walkboxes`, an
empty non-nil slice turns routing off for this scene. They constrain
character movement. When non-empty, `Walk(...)`
routes through a BFS over the polygon adjacency graph (polygons that
share an edge are neighbours), with the midpoint of each shared edge
used as a waypoint. Destinations outside every walkbox are clipped to
the nearest boundary. With no walkboxes the character walks in a
straight line — see [§7.8](#78-character) and [§16.4](#164-character-movement).

`Triggers` fire on the rising edge of their `When` condition. The
engine samples each trigger once per idle frame; see [§7.5](#75-trigger).

`Background` is the picture itself, an `*Asset`, and it need not be registered
anywhere: a scene may nest the literal, or a game may derive one per scene
after registration, the way a game whose backgrounds are named after their
scenes does. Nil leaves the theme's backdrop showing.

**Naming another scene.** `Exit.To` and `GoTo` hold the `*Scene` itself, and
Go's initialisation rule says two package-level variables may not refer to
each other in their initialisers — which two scenes that lead to each other
would. So a scene file declares the scene first and its exits second, with
`Leads`, a method that appends to `Exits` and hands the scene back:

```go
var Alley = SceneManager.Register(&inkwell.Scene{Name: "alley", …})

var _ = Alley.Leads(
    inkwell.Exit{To: Street, Side: inkwell.ExitLeft},
)
```

The second declaration depends on both scenes and neither depends on it, so
the graph may be as tangled as the city is, and every reference is still the
scene.

`Exits` are the scene's connections to other scenes, declared as data. The
engine turns each one into a hotspot — see [§7.3](#73-exit).

`Derived` hands back exits worked out from the rest of the game rather than
written into this scene — the pins on a map screen, which are every location
that says it leads there:

```go
var Selector = SceneManager.Register(&Scene{
    Name:    "selector",
    Derived: selectorPins,
})

func selectorPins(g *inkwell.Game, self *inkwell.Scene) []inkwell.Exit {
    var exits []inkwell.Exit
    for _, s := range g.SceneManager.All() {
        for _, exit := range s.Exits {
            if exit.To == self {
                exits = append(exits, inkwell.Exit{To: s, …})
            }
        }
    }
    return exits
}
```

`Run` asks each scene once, before `Validate`, and appends what comes back to
`Exits`. Everything downstream — validation, hotspot expansion, the debug
menu — therefore sees one kind of exit, and the scene file stays a literal
instead of being patched from the wiring after registration.

### 7.3 Exit

```go
// inkwell/scene.exit.go

type Exit struct {
    To      *Scene      // where it leads; nil means "back the way you came"
    Label   Text        // what the status line calls it
    Side    ExitSide    // where it sits when Area is nil
    Area    Shape       // overrides the edge strip Side would give it
    Needs   string      // flag that has to be set before it opens
    Blocked Action      // what happens while it is not
    OnLook  Action      // overrides Game.ExitLook for this exit
    OnTake  Action      // overrides Game.ExitTake for this exit
}

type ExitSide int
const (
    ExitLeft  ExitSide = iota // off the left edge
    ExitRight                 // off the right edge
    ExitBack                  // into the depth of the picture
    ExitNear                  // out towards the camera
)

func ExitName(to *Scene) string
```

A connection belongs to the scene it leads out of, so it is written in that
scene's own file, right after the scene, with `Leads` — after rather than
inside the literal because the scene it leads to may lead back (see
[§7.2](#72-scene)):

```go
var Alley = SceneManager.Register(&Scene{Name: "alley", …})

var _ = Alley.Leads(
    Exit{
        To:    NoodleHouse,
        Label: Text{{Lang: "en", Text: "back out to the street"}},
        Side:  ExitLeft,
    },
    Exit{
        To:      nil,
        Label:   Text{{Lang: "en", Text: "the fire escape"}},
        Side:    ExitBack,
        Needs:   "has_ladder",
        Blocked: Say(paul, Text{{Lang: "en", Text: "Can't reach it."}}),
    },
)
```

The engine expands each exit into a `Hotspot` with `Cursor: CursorExit`, the
exit's `Label`, `OnUse` bound to `GoTo(To)` — or [`Back()`](#83-built-in-actions)
when `To` is nil — and, when `Needs` is set, the whole travel wrapped in
`If(Flag(Needs), travel, Blocked)`. The exit is visible either way: a locked
door still says it is a door.

`OnLook` and `OnTake` are usually not written per exit. `Game.ExitLook` and
`Game.ExitTake` supply them for every exit in the game, so the phrasing lives
in one place:

```go
g.ExitLook = func(e Exit) Action { return Say(paul, Format(exitLook, e.Label)) }
g.ExitTake = func(e Exit) Action { return Say(paul, exitTake) }
```

**Placement.** With no `Area`, an exit is a strip along one edge of the
picture, sized as a fraction of [`Game.SceneArea()`](#46-scene-helpers) — the
convention every 1990s point & click used, and one field to re-aim at a real
door once the artwork is measured. A game whose HUD covers the foot of the
window sets `Game.SceneRect`, so the strips land inside the painting instead of
under the HUD.

**The graph stays readable.** The generated hotspot is named
`ExitName(To)` — `"exit:noodle_house"`, or `"exit:back"` — so the location
graph can be read straight back out of the registered scenes, and
[`Validate`](#18-validation) rejects an exit that names a scene nobody
registered.

### 7.4 Hotspot

```go
// inkwell/scene.hotspot.go

type Hotspot struct {
    Name   string
    Area   Shape           // Rectangle or Polygon
    Label  Text            // hover hint
    Cursor CursorKind

    // built-in verb handlers
    OnLook Action
    OnUse  Action
    OnTalk Action
    OnTake Action
    OnGive Action

    // when the player has an item selected and clicks this hotspot,
    // the engine looks up the item's Name here first.
    OnUseWith map[string]Action

    // one answer for any item OnUseWith does not name.
    OnUseWithAny Action

    // custom verbs registered through g.VerbManager.
    OnVerb map[string]Action
}

type CursorKind int
const (
    CursorDefault CursorKind = iota
    CursorLook
    CursorUse
    CursorTalk
    CursorTake
    CursorExit
)
```

Hotspots are children of `Scene` — not their own manager. The engine
resolves a click via `hotspot.handler(verbName)` which maps the four
built-in verbs to their `On*` fields and falls back to `OnVerb[verbName]`
for custom verbs.

`OnUseWithAny` answers for an item `OnUseWith` does not name. A slot, a
lock or a person usually has one reply for whatever you put on it, and
enumerating the whole item catalogue to say it is a list that has to be kept
in step with one. It stands between the authored pairs and
`Game.UseWithFail`, so a hotspot with something of its own to say says it,
and everything else still falls through to the game's general refusal.

A hotspot does not have to be in the picture: a HUD widget can own one and
hand it to [`Game.InvokeVerb`](#46-scene-helpers) — see
[`HotspotPanel`](#113-built-in-widgets).

### 7.5 Trigger

```go
// inkwell/scene.trigger.go

type Trigger struct {
    Name string
    When Condition
    Do   Action
    Once bool
}
```

Triggers arm when the scene is entered. Every idle frame (when the
script slot is free), the engine samples each trigger's `When`
condition; on a **rising edge** (false → true) it queues `Do` into the
script slot. `Once: true` makes the trigger fire at most once per
arming. A trigger with `nil` `When` or `nil` `Do` is silently skipped.

Triggers are not sampled while a cutscene or dialog is running, so a
condition that becomes true mid-script is detected on the next idle
frame — the edge is preserved across the busy window. Save/Load resets
trigger state to "freshly armed", matching the behaviour of scene
re-entry.

### 7.6 Item

```go
// inkwell/item.def.go

type Item struct {
    Name        string
    Sprite      *Asset
    Description Text

    OnUseSelf Action
    OnUseWith map[*Item]Action          // this item used on that item
}
```

Items live in the `ItemManager`. When the player picks one up
(`Give(key)`) the item is appended to the `Inventory`, which remembers it by
`Name`. Clicking a hotspot with an item selected resolves the action via, in
order:
1. `hotspot.OnUseWith[item]`
2. `hotspot.OnUseWithAny`
3. Otherwise `Game.UseWithFail(item, hotspot)` is asked for a response and
   the item is deselected; with no hook set the engine flashes "Nem ehhez."

`Item.OnUseWith` is an item used on an *item*: with one in hand, clicking
another in the `InventoryBar` runs whichever of the two pairs them.

`UseWithFail` is the same shape as `ExitLook` and `ExitTake`: the pair the
author never wrote is still content, and the domain is the only one that
knows what its characters say about it.

### 7.7 Inventory

```go
// inkwell/item.inventory.go

type Inventory struct{ /* unexported */ }

func NewInventory() *Inventory
func (i *Inventory) Add(name string)
func (i *Inventory) Remove(name string)
func (i *Inventory) Has(name string) bool
func (i *Inventory) Select(name string)        // "" clears
func (i *Inventory) Selected() string
func (i *Inventory) Items() []string           // copy
```

Not a manager — pure runtime state owned by `Game`. Mutated by actions
(`Give`, `TakeAway`) and by the `InventoryBar` widget on click.

### 7.8 Character

```go
// inkwell/character.def.go

type Character struct {
    Name        string
    Label       Text                         // what the UI shows; empty = Name
    Sprite      *Asset                       // nil draws the placeholder
    Animations  map[string]AnimationClip
    Speed       float64
    SpeechColor color.Color
    Delivery    Delivery                     // DeliveryBubble (default) or DeliveryLog
    Start       Point                        // where Game.Player stands by default
    W, H        float64                      // size hint for placeholder
}

type Delivery int
const (
    DeliveryBubble Delivery = iota  // speech bubble over the scene, plus the log
    DeliveryLog                     // the log alone
)

func (c Character) DisplayName() Text    // Label, or Name

type Talker interface{ /* *Character or *Speaker — what Say and a DialogLine take */ }

func (c *Character) Say(text Text) Action                     // Say(c, text): who, then what
func (c *Character) Reply(answer func(*State) Text) Action
func (s *Speaker) Say(text Text) Action
func (s *Speaker) Reply(answer func(*State) Text) Action
func (g *Game) Speaker(name string) (Speaker, bool)  // a character's voice, or a bodiless one
func (g *Game) SpeakerLabel(name string) Text
func (g *Game) SpeechColor(name string, fallback color.Color) color.Color
func (g *Game) ItemLabel(name string) Text      // Description, or Name
func (g *Game) DeliveryOf(name string) Delivery

type AnimationClip struct {
    Frames    []Rectangle
    FrameTime float64
    Loop      bool
}
```

`Delivery` is where a character reaches the player, and `Say` obeys it: a
`DeliveryLog` character gets the log line without the bubble. A domain never
needs a second spelling of `Say` for that. (It was called `Voice` until a
voice became a recording; see [§14](#14-audio).)

A voice with no body at all is not a character. A radio, a cassette, a
narrator — anything that speaks but never stands in a scene — is a
`Speaker`, registered in `SpeakerManager`: a name, a `Label`, a `Delivery`
and a `Color`, and nothing about walking or being drawn.

```go
type Speaker struct {
    Name  string
    Label Text
    Delivery Delivery
    Color    color.Color
}
```

`Say`, the dialog box, the chat log and the bubble all resolve a
speaker's name through `Game.Speaker(name)`, which looks in the
characters first and the speakers second, so a line does not know or
care which kind said it. `Game.DeliveryOf`, `Game.SpeakerLabel` and
`Game.SpeechColor` are the three answers a widget asks for, and a name
found in neither registry speaks in a bubble under its own name, which
is what a bare `Say` always did.

`Label` is what a HUD prints: a code name, a title, a shout. `Start` is
where `Game.Player` is placed in a scene that lists no actors of its own.
With `Game.WalkOnClick`, `Game.Player` is also who walks: a left click on
empty floor sends that character there, routed through the scene's
walkboxes. It is off by default, so a game says so when it wants it.

When no real sprite is on disk, `drawCharacter` falls back to a stylised
placeholder: humanoid if `W < H`, quadruped otherwise. `SpeechColor` (if
an `RGBA`) colours both the speech-bubble text and the placeholder body.

`Animations` maps a clip name to an `AnimationClip`. The engine
auto-selects clips by name: `"walk"` while the character has an active
path, `"idle"` otherwise. If only one of the two is registered, it
plays in both states. `FrameTime ≤ 0` or `Frames` empty disables the
clip; missing clips fall back to the whole-sprite (or placeholder) draw.

Each `Frames` entry is a source rectangle into the character's `Sprite`
image, blitted via `SubImage` and scaled to the character's `W×H`
footprint. `Loop: true` rewinds to frame 0 after the last frame;
`Loop: false` freezes on the final frame.

Movement is driven by the `Walk` action and `Game.tickCharacters` —
stepping at `Speed` pixels/sec (default 60) along the waypoint list
computed by [walkbox routing](#72-scene).

### 7.9 Dialog

```go
// inkwell/dialog.def.go

type Dialog struct {
    Name  string
    Start string                        // initial node Name; "" = first in slice
    Nodes []DialogNode
}

type DialogNode struct {
    Name    string
    Lines   []DialogLine              // shown one-at-a-time on click
    Choices []DialogChoice
}

type DialogLine struct {
    Speaker Talker                      // a *Character or a *Speaker
    Text    Text
}

type DialogChoice struct {
    Text    Text
    Show    Condition                   // nil = always visible
    Once    bool                        // hide after first pick
    Actions []Action                    // run when chosen
}

func (d Dialog) Node(name string) (DialogNode, bool)
```

The dialog flow:
1. `RunDialog(dialog)` starts the dialog; the runner keeps
   returning `StatusRunning` until the dialog closes.
2. The `DialogBox` widget renders lines one at a time; click advances.
3. After the last line, visible `Choices` are shown.
4. Picking a choice runs `Seq(choice.Actions...)`. `Once` choices remember
   that they fired via `State.NoteTalked`, keyed by the choice's position
   in its node rather than its wording — see [§5.5](#55-language-stable-state).
5. `EndDialog()` closes the conversation; `GotoNode("other")` jumps to
   another node in the same dialog.

### 7.10 Script

```go
// inkwell/script.def.go

type Script struct {
    Name    string
    Actions Action                      // typically Seq(...)
}
```

A `Script` is just a named composite action — useful when you want to
reuse a cutscene (intro, victory, transition) from multiple call sites.
Fire one with `RunScript(script)`.

### 7.11 Verb

```go
// inkwell/verb.def.go

type Verb struct {
    Name    string                      // canonical id, e.g. "look"
    Label   Text                        // display label, e.g. "Look at" / "Nézd"
    Default Action                      // run when a hotspot has no handler
}
```

`New` pre-registers the four SCUMM verbs (`look`, `use`, `talk`,
`take`). Add your own:

```go
g.VerbManager.Register(inkwell.Verb{
    Name: "push",
    Label: inkwell.Text{
        {Lang: "en", Text: "Push"},
        {Lang: "hu", Text: "Lökd"},
    },
    Default: inkwell.Say(player, inkwell.Text{
        {Lang: "en", Text: "It won't move."},
        {Lang: "hu", Text: "Nem mozdul."},
    }),
})
```

The verb-bar and verb-wheel widgets read `VerbManager.Names()` every
frame, so adding a verb at runtime is enough to make it appear.

---

### 7.12 Stage

```go
// inkwell/stage.def.go

type Stage struct {
    Name       string
    Title      Text
    World      bool     // this Stage is the game itself
    Background *Asset   // full-screen image for a Stage that is not the world
    Music      *Music   // a track for this Stage; nil keeps what is playing
    Transition Transition // how this Stage gives way to the next; nil cuts
    Hold       float64  // seconds to show, then move on to Next
    Next       string
    OnEnter    Action
    OnLeave    Action
}
```

A `Stage` is one full-screen thing the game can be showing, and exactly one
of them is showing at a time. A main menu is a Stage, a settings screen is a
Stage, a credits card is a Stage — and so is the game itself, which is the
point: the world is not the thing every other screen floats over, it is one
of the screens.

That is the difference between a Stage and a [Scene](#72-scene). A Scene is a
place in the fiction, and the engine always has one; walking from a room into
an alley is a Scene change *inside* the world. A Stage is what owns the frame.
Off the world Stage the scene is not painted, characters do not walk, triggers
do not fire and a click reaches no hotspot — which is what a game showing its
main menu wants, and what it could not have while the only way to hold the
world still was a script that never finished.

Widgets need no new field to follow this, because `OnStage(stage)` is a
`Condition` like any other:

```go
game := StageManager.Register(&inkwell.Stage{Name: "game", World: true})
mainMenu := StageManager.Register(&inkwell.Stage{Name: "main_menu"})

WidgetManager.Register(&inkwell.Menu{
    Name:   "main_menu",
    When:   inkwell.OnStage(mainMenu),
    Sticky: true,
    Root:   inkwell.MenuItem{ … },
})
```

A HUD hides itself off the world exactly the way it already hides itself for
a cutscene — one condition, no wrapper, no `if` at the top of a `Draw`.

**Changing screens.** `GoToStage`, `Restart` and `LoadSlot` all land at the
top of the *next* frame rather than where they are called. Replacing `State`,
the inventory and every character underneath the action that asked for it is
not something that action should have to survive, and one deferred point
means a menu row can say the whole of what it means in one `Seq` without
having to know what order the engine would otherwise have applied it in:

```go
Do: inkwell.Seq(inkwell.Restart(), inkwell.GoToStage(game)),
```

**Leaving through a Transition.** A `Stage` may name how it gives way:
`Transition` is an `inkwell.Transition`, and a change out of the Stage — a
`GoToStage`, or its `Hold` running out — goes through it. `Fade(seconds)` is
the one the engine ships (`FadeTo` for a colour other than black); nil is a
cut, which is what a menu and its sub-screens want. The whole pending change
lands on the covered frame, a `Restart` or `LoadSlot` asked for in the same
breath included, so the screen still showing during the fade is over the
world it was over. Scene changes go through the same machinery with their own
quarter-second fade.

**A card that moves on.** `Hold` shows a Stage for a number of seconds and
then goes to `Next`, through the Stage's `Transition` — which is all a splash
card or a title sequence is:

```go
StageManager.Register(&inkwell.Stage{
    Name:       "splash_studio",
    Transition: inkwell.Fade(0.5),
    Hold:       2,
    Next:       splashTitle,
})
```

The clock starts once `OnEnter` is done: the hold is a `Wait` and a
`GoToStage` on the same queue, so it lands like any other change, at the top
of a frame. `Validate` rejects a `Next` nobody registered and a `Hold` with
no `Next`.

**Starting a session.** A game whose `StartStage` is the world starts one at
boot: the start scene is entered and `StartScript` is queued. A game that
opens on a menu starts one when the player says so — which is all `Restart`
is. `Started()` reports whether one exists, and is what a "Resume" or "Save"
row is enabled by.

---

## 8. The action system

### 8.1 Action and Runner

```go
// inkwell/action.def.go

type Action interface {
    Start() Runner
}

type Runner interface {
    Tick(ctx *Ctx) Status
}
```

An `Action` is an **immutable spec** of work — safe to store in a
`Hotspot.OnLook` field and use across many clicks. A `Runner` is one
in-flight execution, created by `Start()`, carrying state (elapsed time,
sub-runners, started-flag). The engine never ticks an `Action` directly;
it always asks for a fresh `Runner` per invocation.

Only one runner is active at a time (the "script slot"). New action
queues issued while one is running are dropped (logged at debug level).
A multi-actor scene can still run parallel work via `Par(...)` inside a
single composite action.

### 8.2 Status and Ctx

```go
type Status int
const (
    StatusRunning Status = iota
    StatusDone
    StatusFailed
)

type Ctx struct {
    Game    *Game
    DT      float64                     // seconds since last tick
    Scene   *Scene                      // current scene (may be nil)
    Hotspot *Hotspot                    // populated only for action contexts where applicable
    Item    *Item
}
```

`Status` controls how the parent runner reacts: `StatusDone` advances a
`Seq` to the next step, `StatusFailed` short-circuits the whole branch
(`RequireItem` uses this to abort a sequence when the item is absent).

### 8.3 Built-in actions

| Constructor                              | Behaviour                                                                 |
|------------------------------------------|---------------------------------------------------------------------------|
| `Seq(a ...Action) Action`                | Run children in order. Nested `Seq` are flattened. `nil` children are ignored. |
| `Par(a ...Action) Action`                | Run children simultaneously; completes when every child reports `Done`.   |
| `If(c Condition, then Action, els ...Action) Action` | Run `then` if `c` evaluates true; otherwise `Seq(els...)` (if any).|
| `Wait(seconds float64) Action`           | Block the runner for `seconds`.                                            |
| `Say(speaker Talker, text Text) Action`  | Show the line above the speaker (`SpeechBubble`), append to chat-log. Click-to-skip. Duration scales with the resolved text's length, 1.2s floor. |
| `Reply(speaker Talker, answer func(*State) Text) Action` | `Say` with the line decided when it is spoken: `answer` is handed the `State` and returns the `Text` for this moment, so a speaker can answer from what the player has done without a dialog node per case. |
| `GoTo(scene *Scene) Action`                | Switch the current scene via a fade transition.                            |
| `Back() Action`                          | Return to `PreviousScene()`. No-op when there is nowhere to go back to.    |
| `Walk(character *Character, to Point) Action`| Move a character to `to` at `Character.Speed`. Returns when arrived.       |
| `Give(item *Item) Action`               | Add `item` to inventory.                                                   |
| `TakeAway(item *Item) Action`           | Remove `item` from inventory.                                              |
| `RequireItem(item *Item) Action`        | If the item is not in inventory, flash `MsgNeedItem` and fail.             |
| `SetFlag(name string) Action`            | `State.SetFlag(name)`.                                                     |
| `ClearFlag(name string) Action`          | `State.ClearFlag(name)`.                                                   |
| `SetVar(name string, v any) Action`      | `State.SetVar(name, v)`.                                                   |
| `SetVarIfEmpty(name string, v any) Action` | `State.SetVar(name, v)` only while nothing has written it yet — what a scene's `OnEnter` wants for a first-arrival value. |
| `Noting(name string, v any, inner Action) Action` | Write a var for exactly as long as `inner` runs, and clear it after. |
| `UseTheme(theme *Theme) Action`           | Switch the active theme mid-script, and save the settings.                 |
| `UseLang(lang string) Action`            | Switch the language every `Text` resolves into, and save the settings; reaches the next frame. |
| `GoToStage(stage *Stage) Action`          | Show another [Stage](#712-stage), at the top of the next frame.            |
| `Restart() Action`                       | Begin a new session, at the top of the next frame.                         |
| `LoadSlot(slot int) Action`              | Restore a saved session, at the top of the next frame.                     |
| `SaveSlot(slot int) Action`              | Write the session to a slot, where it is called.                           |
| `Quit() Action`                          | Close the window at the end of the frame.                                  |
| `Flash(text Text) Action`                | Put a line in the status strip for ~2s — `Game.FlashLine` as an action.    |
| `PlayMusic(music *Music) Action`         | Start a track, looping when `Music.Loop`; the same track again is a no-op. |
| `StopMusic() Action`                     | `AudioPlayer.StopMusic`.                                                   |
| `PlaySound(sound *Sound) Action`         | Play a noise once, over whatever else plays.                               |
| `SetVoice(on bool) Action`               | Switch recorded lines on or off, and save the settings.                    |
| `RunDialog(dialog *Dialog) Action`        | Start a registered `Dialog`; runner blocks until the dialog closes.       |
| `EndDialog() Action`                   | Close the active dialog.                                                   |
| `GotoNode(node string) Action`           | Jump to a node within the active dialog.                                   |
| `RunScript(script *Script) Action`          | Run a registered `Script`'s `Actions`.                                     |
| `ShowEnd(text Text) Action`              | Display the end-card overlay; runner never finishes.                       |
| `Custom(fn func(*Ctx) Status) Action`    | Escape hatch — wrap an arbitrary callback as an action.                    |

### 8.4 Custom actions

Implement `Action` and `Runner` directly when you need long-running state
(an animation, a tween, an HTTP poll):

```go
type fadeAction struct{ target string; duration float64 }

func (a *fadeAction) Start() Runner { return &fadeRunner{spec: a} }

type fadeRunner struct {
    spec    *fadeAction
    elapsed float64
}

func (r *fadeRunner) Tick(ctx *Ctx) Status {
    r.elapsed += ctx.DT
    if r.elapsed >= r.spec.duration {
        return StatusDone
    }
    // mutate state, draw a fade, etc.
    return StatusRunning
}
```

For one-shot mutations, `Custom` is shorter:

```go
bumpScore := inkwell.Custom(func(ctx *inkwell.Ctx) inkwell.Status {
    cur := ctx.Game.State.Var("score").(int)
    ctx.Game.State.SetVar("score", cur+10)
    return inkwell.StatusDone
})
```

---

## 9. Conditions

```go
// inkwell/condition.def.go

type Condition interface {
    Eval(ctx *Ctx) bool
}
```

| Constructor                            | Meaning                                          |
|----------------------------------------|--------------------------------------------------|
| `Not(c Condition) Condition`           | `!c`                                             |
| `And(cs ...Condition) Condition`       | All must evaluate true.                          |
| `Or(cs ...Condition) Condition`        | At least one true.                               |
| `Flag(name string) Condition`          | `State.Flag(name)`                               |
| `HasItem(item *Item) Condition`       | `Inventory.Has(name)`                            |
| `SelectedItem(item *Item) Condition`  | `Inventory.Selected() == name`                   |
| `InScene(scene *Scene) Condition`         | The current `Ctx.Scene` is this one.             |
| `OnStage(stage *Stage) Condition`       | The named [Stage](#712-stage) is the one showing. |
| `OnWorld() Condition`                  | The Stage showing is the world, whichever that is. |
| `Started() Condition`                  | A session exists — begun or loaded.              |
| `VarEq(name string, v any) Condition`  | `State.Var(name) == v` (Go equality).            |
| `Given{Flag, NotFlag, Item, NotItem, Var, Is}` | A literal: every set field must hold; `Item` and `NotItem` are `*Item`. `Given{}` holds always. |

Conditions are stateless — they are evaluated lazily, including inside
`If` and `DialogChoice.Show`.

`Given` is the one that is a struct rather than a constructor, and it is
there for the place a condition is read most: a dialog choice, where "only
if the player holds the tape and has not been told yet" is two stacked
fields rather than a call tree.

```go
Show: inkwell.Given{
    Item:    mysteryTape,
    NotFlag: "dex_told",
},
```

Every field set on a `Given` has to hold — `Flag` raised, `NotFlag` not,
`Item` in the inventory, `NotItem` not, `Var` equal to `Is` — so a `Given`
is an `And` of its fields, and `Or(Given{…}, Given{…})` is still how
alternatives read.

---

## 10. World state

```go
// inkwell/state.def.go

type State struct{ /* unexported */ }

func NewState() *State

func (s *State) Flag(name string) bool
func (s *State) SetFlag(name string)
func (s *State) ClearFlag(name string)

func (s *State) Var(name string) any
func (s *State) SetVar(name string, v any)

func (s *State) Visited(name string) int
func (s *State) NoteVisit(name string)             // engine bumps on scene enter

func (s *State) Talked(node string) int
func (s *State) NoteTalked(node string)            // DialogBox bumps on Once-choice pick
```

`Flag`/`SetFlag` are the workhorse for puzzle state ("cupboard_open").
`SetVar` accepts any value — strings for text, numbers for scores, struct
references for richer dashboards. The `TopBar` widget reads vars by name
to render score / time.

`Visited` and `Talked` underpin the "Once" semantics in dialogs and
could be used by triggers once those land.

---

## 11. The widget system

The HUD is a tree of `Widget`s, each registered into `g.WidgetManager`. The
library ships built-in widgets for every classic adventure UI piece;
domain code can register arbitrary new ones — chat panels, minimaps,
hotbars — without touching the library.

### 11.1 Widget interface

```go
// inkwell/widget.def.go

type Widget interface {
    Named
    Tick(ctx *UICtx)
    Draw(dst *ebiten.Image, ctx *UICtx)
}

type Layer int
const (
    LayerScene   Layer = 0    // over the picture, under the HUD
    LayerPanel   Layer = 100  // the plate the HUD sits on
    LayerHUD     Layer = 200  // verbs, inventory, status line, top bar
    LayerSpeech  Layer = 300  // speech bubbles
    LayerDialog  Layer = 400  // the dialog box
    LayerDebug   Layer = 450  // the developer's overlay
    LayerMenu    Layer = 500  // radial verbs, menus
    LayerCurtain Layer = 600  // end cards, fades
    LayerCursor  Layer = 700  // the pointer
)

type Layered interface {
    Layer() Layer
}

func LayerOf(w Widget) Layer  // Layered, else LayerHUD

type Gated interface {
    VisibleWhen() Condition
}

func WidgetVisible(ctx *Ctx, w Widget) bool  // Gated with a false condition = off

type Awake interface {
    TicksDuringScript() bool
}

func WidgetAwake(w Widget) bool  // ticks on the frames a script holds

type UICtx struct {
    Game *Game
    DT   float64
}

type MouseButton int
const (
    MouseButtonLeft MouseButton = iota
    MouseButtonRight
)

type Size struct{ W, H int }

type Align int
const (
    AlignLeft Align = iota
    AlignCenter
    AlignRight
)
```

A widget is responsible for its own input handling, hit-testing, state
and rendering. The library does not impose a layout layer — widgets carry
their own `Bounds` (or compute them dynamically, like `RadialVerbs`).

### 11.2 Z-order and input consumption

- **`Draw` runs from the bottom layer up** — a widget on a higher layer is
  painted on top. Inside one layer, registration order decides.
- **`Tick` runs from the top layer down.** The widget drawn on top gets the
  click first. Each widget calls `ctx.Game.Input.ConsumeLeft()` /
  `ConsumeRight()` to claim the event; widgets below see
  `LeftClicked() == false`.
- **A widget declares its layer**, it does not inherit one from the order it
  was registered in. Every built-in implements `Layered`; a widget that does
  not sits on `LayerHUD`. `Cursor` is on `LayerCursor`, so it always wins on
  visual layer (and effectively never claims clicks) no matter when it was
  registered.
- Layers are what let a domain register its widgets **one file at a time** —
  an `init()` per widget, in whatever order the file names happen to fall —
  without the HUD coming out shuffled.
- A wrapper widget that forwards to an inner one should return
  `LayerOf(inner)` from its own `Layer`, so that wrapping does not move the
  widget.

Every built-in widget also carries a `When Condition`. A widget whose
condition evaluates false neither ticks nor draws nor blocks a click, so a
HUD is hidden for a cutscene by saying so on the widgets rather than by
unregistering them or wrapping each one:

```go
g.WidgetManager.Register(&inkwell.InventoryBar{
    Name: "inventory",
    When: inkwell.VarEq("mode", "play"),
})
```

A domain widget joins in by implementing `VisibleWhen() Condition`; one
that does not implement it is always live, which is what a bare widget
always did.

A running script normally has the frame to itself: no widget ticks until
it is done, which is what keeps a click from racing a cutscene that is
already under way. A widget that implements `Awake` opts out of that —
it ticks on those frames too, **before** the script's runner, so a click
or a key it consumes is not also read as "skip this line".

Two built-ins say yes, and for the same reason. `Menu`, because a menu that
cannot answer a key press while a cutscene plays cannot be closed once one
starts. And `DialogBox`, because a conversation **is** the running script:
`RunDialog` holds the frame for as long as the dialog is up, so the box
has to tick inside it or its lines never advance and no choice is ever
clicked. Nothing else should need it.

After all widgets ticked, the engine offers the (possibly consumed) click
to `handleSceneInput`, which is where hotspot interactions live, and where
a click on nothing walks `Game.Player` when `Game.WalkOnClick` is set. If a widget already consumed the
click, `handleSceneInput` becomes a no-op for that frame.

### 11.3 Built-in widgets

Each widget lives in its own `inkwell/widget.<name>.go` file. Fields with a
zero value fall back to a sensible default.

#### `Panel` (`widget.panel.go`)

A generic colored rectangle, optionally with a stroke border. Useful as a
backdrop for grouping other widgets.

```go
type Panel struct {
    Name        string
    Bounds      Rectangle
    BG          color.Color    // nil → Theme.PanelBG
    Border      int            // 0 = no border
    BorderColor color.Color    // nil → Theme.DialogBorder
    Blocks      bool           // opaque to clicks (see §11.4)
}
```

A plate the HUD is mounted on wants `Blocks`; a rule drawn on top of one
does not, since the plate underneath has already said so.

#### `StatusLine` (`widget.status.go`)

A single line at `Y`. Shows the active flash message (if any) — set by
`Game.FlashLine` — otherwise the composed hover hint (`verb + target` or
`verb + item + target`).

```go
type StatusLine struct {
    Name        string
    Y           int
    Align       Align          // default AlignCenter
    ScreenWidth int            // 0 → Game.Width
}
```

#### `VerbBar` (`widget.verb_bar.go`)

The SCUMM permanent verb panel. Reads `VerbManager.Names()` every frame
and lays the labels into a `Cols × ceil(N/Cols)` grid.

```go
type VerbBar struct {
    Name    string
    Origin  Point
    Cols    int               // default 2
    Button  Size
    Gap     Point
    PanelBG bool              // draw a Theme.PanelBG backdrop behind buttons
}
```

Click on a button → `Game.SetSelectedVerb(name)`. Implements
`BlocksClickAt` so `RadialVerbs` does not pop up over it.

#### `RadialVerbs` (`widget.verb_radial.go`)

Verb-coin or permanent verb-wheel.

```go
type RadialVerbs struct {
    Name    string
    Trigger MouseButton                  // default MouseButtonRight
    Radius  float64                       // default 40

    // AlwaysVisible turns the widget into a permanent wheel at Center.
    // Trigger is ignored in this mode; left-click picks a slice.
    AlwaysVisible bool
    Center        Point

    // Override the rendered label per verb name (e.g. shorten "Használd"
    // → "Tedd" so it fits inside the disk).
    Labels map[string]Text
}
```

In **trigger mode** (default), the wheel is invisible until the trigger
mouse button is pressed somewhere not covered by a `clickBlocker`. The
center latches to the click point. Left-click on a slice picks the verb;
trigger-click again closes without picking.

In **always-visible mode**, the wheel is drawn permanently at `Center`
and left-click on a slice picks the verb. Labels are positioned at
`0.62 × Radius` from the center so short overrides stay inside the rim.

#### `InventoryBar` (`widget.inventory.go`)

Item slots in a grid. Hover sets `Game.SetHoverLabel(item.Description)`;
click toggles selection (or, under the `look` verb, runs a `Say` with
the item description).

```go
type InventoryBar struct {
    Name     string
    Origin   Point
    Slots    int             // total slot count
    Cols     int             // 0 → Slots (single row)
    SlotSize int
    Gap      int
    PanelBG  bool
}
```

Implements `BlocksClickAt` covering the slot grid.

#### `SpeechBubble` (`widget.speech.go`)

Renders whatever `Game.SetSpeech(...)` last set, positioned above the
speaker's head (`character.pos.Y − character.def.H − OffsetY`). Falls back
to a screen-top position if the speaker is unknown.

```go
type SpeechBubble struct {
    Name      string
    MaxWidth  int                // wrap width, default 200
    Padding   int                // bubble padding, default 3
    OffsetY   int                // vertical lift above head
    FallbackY int                // anchor Y when speaker unknown
}
```

Text color follows `Character.SpeechColor` when set, otherwise
`Theme.SpeechDefaultText`.

#### `DialogBox` (`widget.dialog_box.go`)

Dormant unless `Game.dialogActive()`. While active, **consumes every
click** — the rest of the HUD and the scene are inert. The widget renders
either the current line of `runtimeDialog.Node.Lines` (click → next line)
or the visible `Choices` once the lines are exhausted.

It ticks inside the script that is running it — see `Awake` in
[11.2](#112-z-order-and-input-consumption).

```go
type DialogBox struct {
    Name       string
    Bounds     Rectangle      // default Rect(0, 140, Width, 60)
    LineHeight int            // default 14
    Padding    int            // default 6
}
```

Hovered choices use `Theme.DialogChoiceHover`. Once-only choices are
omitted automatically after their first pick (tracked by
`State.NoteTalked`).

#### `EndCard` (`widget.end_card.go`)

A fullscreen overlay shown when `Game.showEndCard(text)` is called (which
the `ShowEnd` action does). Consumes every click so nothing underneath
reacts.

```go
type EndCard struct{ Name string }
```

#### `Cursor` (`widget.cursor.go`)

The crosshair cursor (color from `Theme.CursorColor`). When the inventory
has a selected item with a non-placeholder sprite, that sprite is drawn
at the cursor position instead.

```go
type Cursor struct{ Name string }
```

#### `HotspotDebug` (`widget.hotspot_debug.go`)

Paints a translucent plate over everything on screen a player could point
at, one colour per kind: the scene's hotspots, its exits, the characters
standing in it, the items in the inventory and the targets of every
`HotspotPanel`. The scene's own targets get the element's **name** written
across the plate, because the picture shows them unlabelled; a widget's
targets get the plate alone, because the widget already prints its labels.
F1 toggles it at runtime (override via `ToggleKey`). It sits on `LayerDebug`,
over the scene, the HUD and the dialog, under the menus — so the debug menu
opens on top of it.

```go
type HotspotDebug struct {
    Name      string
    When      Condition
    Enabled   bool
    ToggleKey ebiten.Key                 // 0 → ebiten.KeyF1
    Colors    map[DebugKind]color.Color  // per-kind override; alpha honoured
}

type DebugKind int   // DebugHotspot, DebugExit, DebugCharacter, DebugItem, DebugPanel
```

A widget with clickable things of its own joins the overlay by implementing
`Targeted`: `DebugTargets(g *Game) []DebugTarget`, one `{Kind, Name, Area}`
per thing. `InventoryBar` and `HotspotPanel` do, and a game widget that
embeds one of them inherits it, so the overlay finds a game's slots without
being told about them.

#### `TopBar` (`widget.top_bar.go`)

A horizontal strip across the top with three sections:

- **Left:** `LeftText` override or `Scene.Title` (falling back to `Name`),
  plus `State.Var(NoteVar)` after an em dash when that var holds anything —
  `"ALLEY — paused"`.
- **Center:** score string built from `State.Var(ScoreVar)`. With
  `ScoreMax > 0`, formatted as `"Score: X/MAX"`, else `"Score: X"`.
- **Right:** time string from `State.Var(TimeVar)`.

```go
type TopBar struct {
    Name     string
    When     Condition
    Height   int                  // default 12
    LeftText Text                  // overrides scene title
    ScoreVar string                // State.Var key; "" = no score
    ScoreMax int
    ScoreLabel Text                // overrides the engine's word for it
    TimeVar  string                // State.Var key; "" = no time
    NoteVar  string                // State.Var key; appended to the left text
    Notes    map[string]Text       // what to print for each value NoteVar takes
}
```

The title needs no wiring: a domain that wants "the current scene, and
whatever the game is doing to it" sets `NoteVar` and writes that var.

Empty fields are skipped, so the widget gracefully degrades to two- or
one-section layouts.

#### `Menu` (`widget.menu.go`)

A nested list in a window — see [11.7](#117-the-menu-system) for the tree
model and the keys.

```go
type Menu struct {
    Name       string
    When       Condition
    Root       MenuItem            // first page; its Label titles the window
    Bounds     Rectangle           // zero = the right half of the screen
    ToggleKey  ebiten.Key          // zero = no key, Toggle() only
    Sticky     bool                // the menu IS the screen; cannot be closed
    LineHeight int                 // default GlyphH + 2
    Padding    int                 // default 6
}
```

#### `CharacterPanel` (`widget.character_panel.go`)

A floating info card showing one character's portrait, role badge and
custom stat rows. **Auto-hides** when the character isn't an actor in the
current scene — same registration works across all scenes.

```go
type CharStat struct {
    Label  Text
    VarKey string                  // State.Var(VarKey) is fmt-printed
}

type CharacterPanel struct {
    Name      string
    Bounds    Rectangle
    Character *Character
    Title     Text                  // role badge ("PLAYER", "NPC", "GUIDE", …)
    Stats     []CharStat
}
```

The portrait reuses the placeholder rendering when no sprite is on disk
(humanoid vs. quadruped picked from `W` vs. `H`).

#### `ChatLog` (`widget.chat_log.go`)

Renders the `Game.Messages()` ring buffer with per-kind colors —
`LogAction` lines use `Theme.ChatLogPrompt`, `LogSystem` uses
`Theme.ChatLogSystem`, and a `LogResponse` uses the speaker's own
`Character.SpeechColor` where the character carries one, falling back to
`Theme.ChatLogResponse`. Newest message at the bottom; older messages
scroll out of view as the buffer fills.

```go
type ChatLog struct {
    Name       string
    Bounds     Rectangle
    LineHeight int                 // default 10
    Padding    int                 // default 3
    ShowBorder bool                // draws a 1px border in CharacterPanelBorder

    Speaker    SpeakerStyle        // how a reply names who said it
    Deliveries []Delivery          // keep only replies delivered these ways; nil = all
    Resting    Talker              // who the header carries before anybody speaks
}

type SpeakerStyle int
const (
    SpeakerInline SpeakerStyle = iota  // "DEX: " in front of every reply
    SpeakerHeader                      // the latest speaker once, above a rule
    SpeakerNone                        // nobody is named
)
```

A log can be given over to part of the conversation rather than all of it.
`Deliveries` keeps only the replies delivered one of the ways it lists, which
is how a game with a radio, a tape or a narrator hands that voice a panel of
its own — without the widget having to know what the thing is. The engine
already knows how a character reaches the player
([`Character.Delivery`](#78-character)), so the filter is that fact and not a
second one to keep in step with it:

```go
&inkwell.ChatLog{
    Name:    "channel",
    Bounds:  inkwell.Rect(396, 262, 244, 118),
    Speaker: inkwell.SpeakerHeader,
    Deliveries: []inkwell.Delivery{inkwell.DeliveryLog},
    Resting: "dex",
}
```

`SpeakerHeader` suits a panel given over to one voice at a time: the name
reads as a heading rather than as a word repeated down the margin. Player
commands and system lines are nobody's, so `Deliveries` keeps them either way.

Set `Game.MaxLogLines` to bound the buffer (recommended ~64 for a small
ChatLog).

#### `KeyBinding` (`widget.key_binding.go`)

A key and what it does — the whole widget. It draws nothing and occupies
nothing; it exists so that "Escape opens the main menu" is a registration like
everything else rather than a special case somewhere in the loop.

```go
type KeyBinding struct {
    Name         string
    When         Condition
    Key          ebiten.Key
    Do           Action
    DuringScript bool    // answer on the frames a running script holds
}
```

The key is claimed, so nothing underneath sees it, and `When` is how a binding
belongs to one [Stage](#712-stage). `DuringScript` is what a pause key needs,
since a cutscene is exactly when the player reaches for one.

#### `TextCard` (`widget.text_card.go`)

A block of lines centred on the screen — a credits roll that does not roll, a
title card, an apology. It is the plainest thing a Stage that is not the world
can be made of, and it needs no tree, no rows and no navigation to be one.

```go
type TextCard struct {
    Name       string
    When       Condition
    Lines      []Text    // an empty Text is a blank line
    Align      Align
    LineHeight int       // 0 → GlyphH + 2
    Backdrop   bool      // fill the screen behind the block
}
```

#### `Letterbox` (`widget.letterbox.go`)

The pair of bars that crop the picture top and bottom while a cutscene
plays. It sits on `LayerScene`, so it draws over the picture and under the
HUD — and a game that also hides its HUD for a cutscene gets the whole
convention out of two registrations sharing one condition.

```go
type Letterbox struct {
    Name  string
    When  Condition
    Bar   float64        // each bar's height; 0 → 18
    Color color.Color    // nil → black
}
```

#### `HotspotPanel` (`widget.hotspot_panel.go`)

Makes a region of the HUD answer the verbs, the way a region of the picture
already does. It owns nothing but a hit-test: the hotspots it hands back go
through [`Game.InvokeVerb`](#46-scene-helpers), so a device mounted on the
panel gets look, use, talk and take — and the use-with flow, the log lines
and the nothing-happens fallback — out of the same code a painted hotspot
does.

```go
type HotspotPanel struct {
    Name     string
    When     Condition
    Bounds   Rectangle                    // the region the panel claims
    Hotspots func(g *Game) []Hotspot      // its targets, in screen coordinates
}
```

This is what keeps a game's own HUD from re-implementing the engine. Without
it, a widget with two clickable slots on it has to read the selected verb,
check for a held item, consume the click and answer every combination itself,
in Go, in a switch — a hundred lines of engine re-decided, with the game's
own writing trapped inside it. With it, the slots are `Hotspot` literals like
any other, and the lines live in them.

`Hotspots` is a function rather than a slice because what is mounted on a
panel changes as the game does — a slot answers for whatever is in it — and a
widget that asks each time needs no invalidation.

The panel draws nothing. A game that wants its slots painted embeds it and
declares a `Draw` of its own:

```go
type tapeSlots struct {
    inkwell.HotspotPanel
}

func (slots *tapeSlots) Draw(dst *ebiten.Image, ctx *inkwell.UICtx) { … }
```

### 11.4 The `clickBlocker` contract

Some widgets occupy a fixed area; clicking inside should not pop a
`RadialVerbs` menu over them. They implement the optional library-internal
interface:

```go
type clickBlocker interface {
    BlocksClickAt(p Point) bool
}
```

Built-ins that implement it: `VerbBar`, `InventoryBar`, `DialogBox`,
`TopBar`, `CharacterPanel`, `ChatLog`, `RadialVerbs` (itself, when
visible), `Menu` (its window, while it is open). Implement it on your own widgets to mark them as solid input
zones for the radial menu.

### 11.5 Registration helpers

```go
// inkwell/widget.preset_scumm.go, widget.preset_radial.go, widget.preset_rich.go

func RegisterDefaultUI(g *Game)                              // SCUMM-style verb bar + inventory
func RegisterRadialVerbUI(g *Game)                            // verb-coin (right-click) + wider inventory
func RegisterRichUI(g *Game, player, npc *Character)         // top bar + character panels + chat log + verb wheel
```

Pass `""` for any character name to skip its `CharacterPanel`. All three
helpers register in conventional Z-order (debug overlay back, cursor
front).

### 11.6 Custom widgets

Anything implementing `Widget` plugs in. The library doesn't care about
widget identity beyond the name; it just calls `Tick` and `Draw`. Read
state from `ctx.Game`, consume input via `ctx.Game.Input`, draw with
Ebitengine primitives, and you're done.

```go
type Minimap struct {
    Name   string
    Bounds inkwell.Rectangle
}

func (m *Minimap) GetName() string                                { return m.Name }
func (m *Minimap) Tick(ctx *inkwell.UICtx)                         {}
func (m *Minimap) Draw(dst *ebiten.Image, ctx *inkwell.UICtx) {
    // ... render scene thumbnail, mark NPCs, etc.
}

g.WidgetManager.Register(&Minimap{
    Name:   "minimap",
    Bounds: inkwell.Rect(220, 4, 96, 56),
})
```

### 11.7 The menu system

`Menu` is a nested list in a window: a page of rows, any of which can open
a page of its own, as deep as the tree goes. It is a widget like any other
— on `LayerMenu`, gated by `When`, registered from its own file — and it is
the whole of the menu system. A debug menu, a pause menu and an in-fiction
device screen are this one widget with a different `Root`.

A row and a page are **the same type**. `MenuItem.Label` names the row on
its parent's page and titles the page it opens, which is what makes the
breadcrumb in the title bar fall out for free:

```go
// inkwell/widget.menu.go

type MenuItem struct {
    Label    Text
    Value    func(*Ctx) Text        // right-hand column: a count, a target, "held"
    Items    []MenuItem             // the page this row opens, as authored
    Expand   func(*Ctx) []MenuItem  // the page this row opens, read out of the game
    Do       Action                 // what Enter queues on a leaf
    Enabled  Condition              // false dims the row and refuses Enter
    KeepOpen bool                   // leave the menu up after Do
}

func FixedValue(t Text) func(*Ctx) Text   // a Value that does not depend on the game
```

A row with neither `Items` nor `Expand` is a leaf, and Enter hands its `Do`
to `Game.Do`. The two page fields are one decision: **`Items` is what an
author typed, `Expand` is what the game knows.** An `Expand` page is rebuilt
while it is open, so a list of what is clickable in the current scene
rewrites itself when the player leaves the scene; an `Items` page costs
nothing per frame.

`Value`, `Expand` and `Enabled` all take the same `*Ctx` an action or a
condition takes, so `ctx.Scene` is there without being asked for.

The menu answers what a list has always answered to, and **claims every key
it uses** (see [15](#15-input)) so nothing underneath sees it:

| | |
|---|---|
| up / down, wheel | walk the rows, wrapping |
| PageUp / PageDown / Home / End | jump |
| Enter, right | open the row, or run it |
| Backspace, left, right mouse button | back a level; at the root, close |
| Escape, the `×` in the title bar | close |
| hover | move the cursor — only when the pointer itself moved, so it does not fight the keys |
| left mouse button | open the row under the pointer |
| `←` in the title bar | back a level |

Only clicks inside the window are consumed, so a menu down one side of the
screen leaves the other side playable. `ToggleKey` opens and closes it;
zero binds no key at all and leaves that to `Open`, `Close` and `Toggle`,
which is how a HUD button would drive one.

`Sticky` makes the menu the screen rather than a window over one: it is open
whenever it is visible, Escape and back-at-the-root do nothing, `ToggleKey` is
ignored, and the close box is not drawn. A main menu at boot has nothing behind
it to close onto, and a menu that can be dismissed into an empty frame is a way
to lose a game rather than to leave one. A sticky menu is what a
[Stage](#712-stage) is usually made of, and its key legend says so — it offers
no ESC.

`Menu` implements `Awake` ([11.2](#112-z-order-and-input-consumption)): it
ticks while a script holds the frame, because a menu that goes deaf when a
cutscene starts is a menu that cannot be closed.

Colours come from the dialog tokens of the active `Theme` — `DialogBG`,
`DialogBorder`, `DialogSpeaker` for the title bar, `DialogChoiceBG` behind
the cursor, `DialogChoiceHover` for the row it is on, `TopBarAccent` for
the value column — so a menu needs no theme work to look like the rest of
the game. A disabled row is drawn darkened rather than in a colour of its
own, which needs no token and works in every theme.

Two engine strings sit under the menu: `MenuEmpty`, what a page with no
rows says, and `MenuKeys`, the legend along the foot of the window. Both
are in `i18n.strings.go` with the rest of the engine's voice
([5.4](#54-what-the-engine-says-itself)).

### 11.8 The debug menu

`DebugMenuRoot` is the one use of the menu system the library ships: what
the game is made of, and what is clickable in the scene the player is
standing in. It is **data, not a widget** — the menu is the widget:

```go
g.WidgetManager.Register(&inkwell.Menu{
    Name:      "debug_menu",
    Root:      inkwell.DebugMenuRoot,
    When:      whenDebug,
    ToggleKey: ebiten.KeyD,
})
```

```
DEBUG
├── Scenes         every registered scene, in registration order
│   └── <scene>            value: "here", or ×N visits · Enter teleports
├── Items          every registered item
│   └── <item>             value: "held" / "in hand"
│       ├── give · take away · select
│       ├── use on self    only when the item has an OnUseSelf
│       └── use with → <target>   the item's own OnUseWith pairs
├── Action points  the current scene's hotspots, exits included, in hit-test order
│   └── <hotspot>          value: its CursorKind
│       ├── <verb>         only the verbs the hotspot answers
│       └── use with → <item>   the hotspot's OnUseWith pairs
└── Exits          the current scene's exits, as the graph sees them
    └── exit:<to>          value: which edge it sits on
        ├── go
        └── needs <flag>   value: on / off
            └── set · clear
```

The rows name **registry keys**, not the prose a player reads: a debug menu
is read by whoever writes the content, and `noodle_house` says more there
than `"NOODLE'S HOUSE — KITCHEN"` does, in any language.

Every verb goes through `invokeHotspotVerb`, the same path a click takes,
so the chat log and the status line read as if the player had done it. A
`use with` row hands the item over first — the point is to watch the
pairing run, not to be reminded that the player never picked the thing up.

Nothing registers it, `RegisterDefaultUI` included: a debug menu is
something a game asks for, and gates on whatever its debug builds set. A
game that wants sections of its own builds a root around it:

```go
Root: inkwell.MenuItem{
    Label: inkwell.DebugMenuRoot.Label,
    Items: append(inkwell.DebugMenuRoot.Items, mySection),
},
```

---

## 12. Themes

```go
// inkwell/theme.def.go

type Theme struct {
    Name string

    // panel + status
    PanelBG    color.Color
    StatusText color.Color
    FlashText  color.Color

    // verbs
    VerbButtonBG         color.Color
    VerbButtonSelectedBG color.Color
    VerbButtonText       color.Color

    // inventory
    InventorySlotBG         color.Color
    InventorySlotSelectedBG color.Color

    // speech
    SpeechBubbleBG    color.Color
    SpeechDefaultText color.Color

    // dialog
    DialogBG          color.Color
    DialogBorder      color.Color
    DialogChoiceBG    color.Color
    DialogChoiceHover color.Color
    DialogSpeaker     color.Color
    DialogText        color.Color

    // misc overlays
    EndCardBG      color.Color
    EndCardText    color.Color
    CursorColor    color.Color
    HotspotOutline color.Color
    SceneBackdrop  color.Color

    // rich-UI widgets
    TopBarBG       color.Color
    TopBarText     color.Color
    TopBarAccent   color.Color

    ChatLogBG       color.Color
    ChatLogPrompt   color.Color
    ChatLogResponse color.Color
    ChatLogSystem   color.Color

    CharacterPanelBG     color.Color
    CharacterPanelBorder color.Color
    CharacterPanelTitle  color.Color
}
```

### 12.1 Preset themes

`New` automatically calls `RegisterPresetThemes(g)` and selects `StartTheme`,
or `classic-scumm` when the game names none. `RegisterPresetThemes` adds
rather than replaces — a domain that hands over its own `ThemeManager` keeps
everything already in it, and a preset whose name the domain has taken is left
alone. The four shipped presets:

| Name              | Mood                                                        |
|-------------------|-------------------------------------------------------------|
| `classic-scumm`   | Default SCUMM-era dark blue + warm amber accents.          |
| `sierra-coin`     | Darker purple + orange; designed for the verb-coin layout. |
| `paper-notebook`  | Cream paper + ink + a touch of sepia for the highlights.   |
| `terminal-green`  | Retro phosphor green on black.                              |

### 12.2 Custom themes

Register your own:

```go
g.ThemeManager.Register(inkwell.Theme{
    Name:      "midnight-noir",
    PanelBG:   color.RGBA{6, 6, 12, 255},
    StatusText: color.White,
    // ... ~30 color fields ...
})
g.UseTheme("midnight-noir")
```

Switching at runtime is just `g.UseTheme(name)` — widgets read the active
theme every frame, so the change shows up on the next draw.

---

## 13. Asset pipeline

`AssetManager.Register` only stores the spec (`Name`, `Path`, `Kind`). The
first request for an image goes through `loadedAssets.image(name)`:

1. Look up the file from `Asset.Path`.
2. If the file is missing or unreadable, generate a deterministic
   placeholder color from `fnv32a(name)`. Mark the entry as a placeholder.
3. Cache the decoded `*ebiten.Image` for subsequent calls.

```go
// internal:
type loadedAssets struct {
    images       map[*Asset]*ebiten.Image
    placeholders map[*Asset]bool
}

func (la *loadedAssets) image(asset *Asset) *ebiten.Image
func (la *loadedAssets) isPlaceholder(asset *Asset) bool
```

`isPlaceholder` is used by widgets to decide whether to render a stylised
character placeholder (humanoid / quadruped) or the real sprite.

Supported image formats: PNG and JPEG (registered via blank imports in
`asset.manager.go`). Sound is not an `Asset` at all: a `Voice`, a `Music`
or a `Sound` is decoded on demand by `AudioPlayer` — see [§14](#14-audio). The engine's own font is not an
asset at all: it is the atlas baked into `font.png` — see [§5.6](#56-the-font).

---

## 14. Audio

```go
// inkwell/audio.def.go

type Audible interface{ Source() string }

type Voice struct{ Path string }                    // a recorded line
type Music struct{ Name, Path string; Loop bool }   // a track
type Sound struct{ Name, Path string }              // a noise

// inkwell/asset.audio.go

type AudioPlayer struct{ /* unexported */ }

func NewAudioPlayer() *AudioPlayer
func (a *AudioPlayer) PlayMusic(music *Music)
func (a *AudioPlayer) StopMusic()
func (a *AudioPlayer) PlayVoice(voice *Voice) float64   // seconds it lasts; 0 if nothing played
func (a *AudioPlayer) StopVoice()
func (a *AudioPlayer) PlaySound(sound *Sound)
```

Three things are heard, and they are three types because they are played
three ways:

- **`Voice`** is a recorded line. It sits on a `Translation`, because the
  English and the Hungarian recording are two files, and `Say` plays the one
  the resolved language carries — for exactly as long as the recording lasts,
  in place of the text-length timer. A skipped line stops talking. A `Voice`
  is owned by one line, so it needs no name and no registration:
  `{Lang: "hu", Text: "Katt.", Voice: &Voice{Path: "assets/voice/alley/click.hu.ogg"}}`.
  The player switches voices off in the settings (`Settings.Voice`,
  `SetVoice`), and `Say` falls back to the timer.
- **`Music`** is a track a `Scene` or a `Stage` plays while it shows, looping
  when `Loop` says so. Entering a Scene or Stage whose `Music` is the one
  already playing changes nothing; a Scene or Stage with nil `Music` keeps
  whatever plays. `PlayMusic` is the same thing as an action.
- **`Sound`** is a noise — a door, a latch, a cassette clicking home —
  played once by `PlaySound`, over the music and over another `Sound`.

Backed by Ebiten's `audio` package. The shared `*audio.Context` is
created lazily on the first play, at a sample rate of 48 kHz; existing
audio contexts (e.g. created by a host process) are detected via
`audio.CurrentContext` and reused. Music streams from disk; voices and
sounds are decoded once into PCM, cached by `Source`, and replayed via
`NewPlayerFromBytes`. Finished sound players are pruned on the next
`PlaySound`.

Supported formats: **WAV** (`.wav`), **OGG Vorbis** (`.ogg`), and
**MP3** (`.mp3`), picked by file extension. Anything else fails to
decode and the source is flagged so subsequent plays no-op.

Failure modes (missing file, unsupported codec, no audio device) all
degrade to a debug-log line; the game keeps running with no sound.

---

## 15. Input

```go
// inkwell/input.def.go

type Input struct{ /* unexported */ }

func (i *Input) Pos() (int, int)
func (i *Input) Point() Point
func (i *Input) LeftClicked() bool       // just-pressed AND not consumed
func (i *Input) RightClicked() bool
func (i *Input) ConsumeLeft()
func (i *Input) ConsumeRight()

func (i *Input) Wheel() (float64, float64)         // this frame's travel; nothing consumes it

func (i *Input) KeyPressed(key ebiten.Key) bool    // just-pressed AND not consumed
func (i *Input) ConsumeKey(key ebiten.Key)
func (i *Input) ClaimKey(key ebiten.Key) bool      // KeyPressed + ConsumeKey
```

The engine `poll`s once per frame, snapshotting cursor position, the
just-pressed state for both buttons, the wheel's travel and every key that
went down. Widgets and the engine use the **consume-on-use** pattern: one
click is authoritative — whoever consumes first wins, later checks see
`LeftClicked() == false`.

Keys work the same way, and `ClaimKey` is the shape a widget that binds
half a dozen of them wants: a menu walking a list with the arrow keys
claims them, so nothing else reads them on that frame. `KeyPressed` is
edge-triggered like `LeftClicked` — holding a key answers true once — and
a widget that wants a key held down reads `ebiten.IsKeyPressed` itself.

`Game.Input` is exposed so widgets and custom actions can read it
directly.

---

## 16. The game loop

```go
// inkwell/core.run.go

func Run(g *Game) error                          // toplevel — same as g.Run()
```

`Run` does, in order:

1. `g.Audio.attach(g.AssetManager)` — wire the audio player to whichever
   asset registry the game is carrying by now.
2. `g.Validate()` — cross-check name references between managers.
3. If `WidgetManager` is empty, call `RegisterDefaultUI(g)`.
4. Place the start scene directly (no transition), bump
   `State.NoteVisit`, position registered actors, kick off music.
5. Compose `Seq(scene.OnEnter, OnStart script, OnFinale script)` and queue
   it as the initial action — the first script tick runs them in order.
6. `ebiten.SetWindowSize(Width*WindowScale, Height*WindowScale)` (4 when
   unset), `ebiten.SetWindowTitle(g.Title)`, then
   `ebiten.RunGame(&engine{g})`.

### 16.1 `engine.Update`

```
apply the pending screen change   // GoToStage / Restart / LoadSlot
if quitting → return Termination

poll input
update transition
pump the action queue      // start the next Do() when nothing is running
if transition fading out → return

world = the current Stage is the world

if scriptRunner != nil:
    tick the Awake widgets, top layer down   // a menu, so it stays closable
    tick the runner
    if world: tick characters
    return

clear hoverLabel
if world: tick scene triggers
for w in WidgetManager, top layer down:
    if not WidgetVisible(w): continue
    w.Tick(uictx)        // widgets consume input top-down

if world:
    handleSceneInput()   // hotspot resolution + right-click reset
    tick characters
```

The pending screen change is applied before anything is ticked, so the world a
`Restart` tears down is not one a runner is standing in the middle of. Off the
world [Stage](#712-stage) the scene does not tick and is not drawn: a click
that falls through a menu reaches nothing rather than the alley behind it.

### 16.2 `engine.Draw`

```
fill Theme.SceneBackdrop (if any)
draw scene background image
for c in characters sorted by Y:
    drawCharacter(c)
for w in WidgetManager (by layer, then registration order):
    if not WidgetVisible(w): continue
    w.Draw(screen, uictx)
draw transition overlay
```

### 16.3 `handleSceneInput`

Runs only after every widget had a chance. The flow:

1. Set hover label from the hotspot under the cursor (if any).
2. On right-click: deselect the held item if any, otherwise reset
   `selectedVerb` to `"look"`.
3. On left-click:
   - If an item is selected, try `hotspot.OnUseWith[item]`, then
     `item.OnUseWith[hotspot]`, then `Game.UseWithFail`, and only then a
     "Nem ehhez." flash.
   - Otherwise look up `hotspot.handler(selectedVerb)`. Fall back to the
     verb's `Default` action. Fail with "Semmi említésre méltó."
   - If no hotspot is under the cursor and the click is inside
     `SceneArea`, and `Game.WalkOnClick` is set, walk `Game.Player` to the
     click, through the walkboxes. This calls `walkCharacter` directly rather
     than queueing a `Walk` action: a queued walk would be a script, and
     a script owns the frame, so the next click could not redirect it.
4. On every successful click, push a `LogAction` line so the `ChatLog`
   widget can display it.

### 16.4 Character movement

`Game.tickCharacters` runs every frame. For each moving character it
steps the head of its waypoint queue at `def.Speed` pixels per second
(default 60); on arrival, the waypoint is popped and the character
proceeds to the next one. The character becomes idle when the queue
empties.

The waypoint queue is built when `Walk(name, to)` runs — see
[§7.2](#72-scene) and [`scene.path.go`](#22-project-layout) for the
walkbox routing. The `Walk` runner returns `StatusRunning` until
`characterMoving(name)` reports false.

Each frame also advances the character's animation clock via
`tickAnimation`, which auto-selects `"walk"` (when `moving == true`)
or `"idle"` from `Character.Animations`. The renderer uses the active
clip's source rectangle to blit a sprite-sheet frame, or falls back to
the whole-sprite / placeholder draw when no clip applies.

---

## 17. Scene transitions

```go
// inkwell/transition.def.go

type Transition interface {
    Seconds() float64
    Draw(dst *ebiten.Image, width, height int, cover float64)
}

func Fade(seconds float64) Transition
func FadeTo(seconds float64, to color.Color) Transition
```

A `Transition` is how one picture gives way to another, in two halves: the
old picture is covered, the swap lands on the frame it is fully hidden, and
the new picture is uncovered. `Seconds` is one half; `Draw` paints the cover
for `cover` running 0 → 1 and back. `Fade` is the one the engine ships, and
a wipe or an iris is another implementation of the same two methods. A
[Stage](#712-stage) names its own in `Stage.Transition`; a scene change uses
`Fade(0.25)`.

The one in flight is `core.transition.go`'s unexported `transition`. A scene
change calls `Game.changeScene(name)`, which starts it. At the midpoint:

1. `prevScene.OnLeave` (if any) is queued.
2. `currentScene` is set to the new name.
3. `State.NoteVisit(name)` bumps the visit counter.
4. Actors are placed at their `SceneActor.At` positions.
5. The new scene's `Music` plays.
6. `newScene.OnEnter` is queued.

The fade-from-black plays after the midpoint. During the fade-out half,
`engine.Update` returns early — input is frozen briefly.

The initial scene (`StartScene`) is set **without** a transition so the
intro script runs immediately.

---

## 18. Validation

```go
func (g *Game) Validate() error
```

Run by `Run` before entering the Ebiten loop — after every
[`Scene.Derived`](#72-scene) has been resolved, so a derived exit is checked
like an authored one — and again from your tests. It cross-checks:

- `StartScene` and `StartStage` are set, and registered.
- A `Stage.Hold` has a `Next` to move on to, and that `Next` is registered.
- Every `SceneActor.Character` and every `Exit.To` is registered (a nil `To`
  is allowed — it means "back the way you came").
- `Player`, when set, is a registered character.
- An active theme is selected and registered.

Most of what this used to check is now a type: a scene, an item or a dialog
is referred to by pointer, so a name nobody registered is a compile error
rather than one of these.

Returns the first error wrapping one of the `Err...` sentinels (so callers
can `errors.Is` against them). Duplicate-name registration is caught
eagerly at `Register` time via panic.

---

## 19. Save and load

```go
// inkwell/state.save.go

func (g *Game) Save(slot int) error
func (g *Game) Load(slot int) error
func (g *Game) Slots(count int) []SlotInfo

type SlotInfo struct {
    Slot    int
    Exists  bool
    Scene   string
    SavedAt time.Time
}
```

`Slots` reads the header of every slot from 1 to `count`, which is what a
load or save screen lists. A slot that is missing or unreadable comes back as
one that does not exist — a screen listing saves is not the place to report a
parse error — so a game showing four slots gets four rows, three of which say
"empty".

```go
// inkwell/state.store.go

type Store interface {
    Read(key string) ([]byte, error)   // ErrNotStored for a key nothing has written
    Write(key string, data []byte) error
}

func NewFileStore(app string) *FileStore   // native: os.UserConfigDir()/<app>/<key>.json
func NewMemoryStore() *MemoryStore         // a map, for a test

// inkwell/state.settings.go

type Settings struct {
    Version     int
    Lang        string
    Theme       string
    WindowScale int
    Voice       bool             // recorded lines on or off; on until switched
    Extra       map[string]any   // the game's own, kept untouched
}

func (g *Game) Settings() *Settings
func (g *Game) SaveSettings() error
```

**Where it all goes.** A `Game` keeps what has to outlive the process through
its `Store`: two documents, `settings` and `slots`. `Game.App` names the game
to the machine, and the platform picks the place. Natively that is a
directory under the user's configuration directory — `~/Library/Application
Support/<app>` on macOS, `~/.config/<app>` on Linux, `%AppData%\<app>` on
Windows — with one JSON file per document, written to a temporary file and
renamed into place so a crash mid-write leaves the previous document. In a
browser (`GOOS=js`) it is `localStorage`, one entry per document, prefixed
with the app name. A game that wants neither hands in its own `Store`; a
test hands in a `MemoryStore`.

**Settings.** `Run` reads the `settings` document before the first frame:
the stored language is used when the run has not chosen one (`inkwell.Lang`
left empty), the stored theme when it is registered, the stored window scale
when the game named none — so a `-lang` flag wins for this run without
being written back. `UseLang` and `UseTheme` call `SaveSettings`, so a
choice on a settings screen is kept the moment it is made. `Extra` is the
game's: read and write it through `Settings()`, call `SaveSettings`, and it
round-trips beside the engine's fields.

```json
{"version": 1, "lang": "hu", "theme": "realworld-93", "window_scale": 2, "voice": true}
```

**Slots.** Every save slot lives in the one `slots` document, keyed
`slot<N>`. What a load screen lists — when, where — sits in the clear beside
the session, which is the save's JSON as a base64 string, so listing four
slots decodes no session and the document reads as a table of slots rather
than as four games. `Slots` reads that metadata; `Save` and `Load` encode
and decode one entry.

```json
{
  "version": 1,
  "slots": {
    "slot1": {"saved_at": "…", "scene": "alley", "title": "Real World", "state": "eyJ2ZXJzaW9uIjoy…"}
  }
}
```

The session inside keeps its own schema version (2) independent of the
document's (1). The save captures the **mutable runtime state**:

- `currentScene` and `previousScene`, the active verb, and the active theme.
- Every character's position, target, and moving flag.
- The full `State`: flags, vars, visited and talked counters.
- The inventory item list plus the currently selected slot.
- The time it was written (`saved_at`, schema version 2), so a load screen
  can say when.

Static registrations (items, scenes, characters, dialogs, scripts,
assets, verbs, widgets, themes) are **not** persisted — they are
reconstructed by `domain.Build()` on every launch. References in the
save are validated against the current managers before any runtime
state is touched, so a corrupt or schema-drift save returns an error
without mutating the live `*Game`.

What is **not** saved:

- The in-flight script `Runner` or any active dialog. `Load` cancels
  both, so saves are effectively taken at an idle/ready boundary.
- The animation playback state. After load, characters are placed at
  the saved position; the next idle frame re-derives the clip.
- Walkbox path queues. If a character is still `moving`, `tickCharacters`
  rebuilds a one-step straight-line path from `pos` to `target` so the
  walk finishes.

`Var` values are JSON-marshaled as-is. Numbers round-trip through
`float64`, strings stay strings, and any value that survives
`encoding/json` round-tripping survives the save.

The save format is versioned (`saveVersion`); loading a file from a
future version returns an error rather than silently truncating fields.

---

## 20. Errors

```go
// inkwell/core.errors.go

var (
    ErrUnknownAsset     = errors.New("inkwell: unknown asset")
    ErrUnknownScene     = errors.New("inkwell: unknown scene")
    ErrUnknownCharacter = errors.New("inkwell: unknown character")
    ErrUnknownScript    = errors.New("inkwell: unknown script")
    ErrNoStartScene     = errors.New("inkwell: StartScene not set or unknown scene")
    ErrUnknownStage     = errors.New("inkwell: unknown stage")
    ErrNoStartStage     = errors.New("inkwell: StartStage not set or unknown stage")
    ErrNotStored        = errors.New("inkwell: nothing stored under this key")
)

Every one of them is returned by `Game.Validate` — an error the library
never returns is not a sentinel, it is a name to keep in step, so the list
is exactly what the validator can say.
```

All wrap-able with `fmt.Errorf("%w: …", err, ctx)` and checkable with
`errors.Is`. `Manager.Register` and `MustGet` use `panic` instead of
returning these errors because they signal construction-time bugs that
should crash loudly in development.

---

## 21. Testing

```bash
go build ./...
```

The library is single-package and has no in-tree unit tests at this
milestone — the manager + action primitives are intentionally small
enough that the `inkwell-demo` repo's `build_test.go` (a headless
`Build() + Validate()` smoke test) covers the integration surface:

```go
func TestBuildValidates(t *testing.T) {
    g := domain.Build()
    if err := g.Validate(); err != nil {
        t.Fatalf("validate: %v", err)
    }
}
```

Drop a similar test into any project that consumes the library — it
opens no window and finishes in milliseconds, so it's a cheap CI gate.

Library-internal unit tests are not shipped yet; the building blocks are
small enough to verify through smoke runs
of the demo for the time being.

### 21.1 Useful runtime debug switches

```go
inkwell.DebugLog = true                 // logf output to stderr (script queue, audio, scene change)
&inkwell.HotspotDebug{Enabled: true}    // start with the F1 overlay on
&inkwell.Menu{Root: inkwell.DebugMenuRoot, ToggleKey: ebiten.KeyD}  // see 11.8
```

---

## 22. Project layout

All library source files live at the repo root; no nested package
directory. The naming convention is `category.identifier.go` — `ls widget.*`
or `ls action.*` instantly groups a category, and the category is the word
the entity itself uses: `character.*` holds `Character`, `script.*` holds
`Script`, `dialog.*` holds `Dialog`. There is no `ui.*` group, because
"UI" is not an entity anybody registers — a widget, a theme and a verb are
three categories that happen to be drawn. Below the category file that
holds a group's shared contract, each file holds **one** thing: one action,
one condition, one preset theme, one widget. Adding a built-in action is
adding a file; the tree below writes those runs of files as a glob rather
than listing all twenty-three.

```
inkwell/                        # module git.teletypegames.org/engines/inkwell
├── go.mod
├── go.sum
│
├── core.doc.go                 # package docs
├── core.manager.go             # Manager[T Named], Named, TypeLabel
├── core.game.go                # Game aggregate, its managers and accessors
├── core.transition.go          # the one Transition in flight, scene or Stage
├── core.validate.go            # Game.Validate
├── core.log.go                 # LogKind, LogMessage, the chat-log buffer
├── core.engine.go              # ebiten.Game adapter (Update/Draw/Layout)
├── core.run.go                 # Run()
├── core.errors.go              # sentinel errors
│
├── util.geometry.go            # Point, Rectangle, Polygon, Shape
├── util.log.go                 # DebugLog + logf
│
├── asset.def.go                # Asset: a picture
├── audio.def.go                # Audible: Voice, Music, Sound
├── asset.manager.go            # AssetManager alias + lazy loader
├── asset.audio.go              # AudioPlayer (WAV/OGG/MP3): music, voices, sounds
├── asset.text.go               # drawText / wrapText helpers
├── asset.fs.go                 # AssetFS: an embed.FS or the OS filesystem
│
├── i18n.text.go                # Translation, Text, Format, Lang
├── i18n.strings.go             # the engine's own wording, per language
├── i18n.font.go                # the glyph atlas and its lookup
├── i18n.glyphs.go              # generated: the atlas coverage
├── font.png                    # generated: the atlas itself
├── internal/fontgen/           # the generator (its own module)
│
├── scene.def.go                # Scene, SceneActor
├── scene.manager.go            # SceneManager alias
├── scene.hotspot.go            # Hotspot, CursorKind
├── scene.exit.go               # Exit, ExitSide + edge-strip geometry
├── scene.trigger.go            # Trigger + rising-edge engine sweep
├── scene.path.go               # walkbox routing (BFS over polygon adjacency)
├── scene.change.go             # changeScene: the fade, and what happens mid-fade
│
├── item.def.go                 # Item
├── item.manager.go             # ItemManager alias
├── item.inventory.go           # Inventory
│
├── character.def.go            # Character
├── character.manager.go        # CharacterManager alias
├── character.animation.go      # AnimationClip + tickAnimation
├── character.runtime.go        # runtimeChar: placing, walking, ticking the cast
├── character.draw.go           # the cast on screen, placeholder shapes included
│
├── speaker.def.go              # Speaker: a voice with no body, and how a name resolves to one
├── speaker.manager.go          # SpeakerManager alias
│
├── dialog.def.go               # Dialog, DialogNode, DialogChoice
├── dialog.manager.go           # DialogManager alias
├── dialog.runtime.go           # the conversation in flight
│
├── action.def.go               # Action, Runner, Ctx, Status: the contract
├── action.queue.go             # Game.Do and the one-at-a-time queue
├── action.custom.go            # Custom
├── action.*.go                 # one built-in action per file: say, goto, wait, walk, …
│
├── condition.def.go            # Condition
├── condition.*.go              # one condition per file: flag, has_item, and, or, …
│
├── script.def.go               # Script
├── script.manager.go           # ScriptManager alias
│
├── state.def.go                # State (flags, vars, visited, talked)
├── state.store.go              # Store, MemoryStore: where settings and saves live
├── state.store_native.go       # FileStore under os.UserConfigDir (!js)
├── state.store_js.go           # LocalStore over window.localStorage (js)
├── state.settings.go           # Settings: lang, theme, window scale, the game's Extra
├── state.save.go               # Save/Load: the slots document, one base64 session per slot
│
├── input.def.go                # Input (consume-on-use)
│
├── theme.def.go                # Theme
├── theme.manager.go            # ThemeManager alias
├── theme.presets.go            # RegisterPresetThemes
├── theme.*.go                  # one preset theme per file
│
├── verb.def.go                 # Verb
├── verb.manager.go             # VerbManager alias
├── verb.defaults.go            # the built-in SCUMM verb set
│
├── widget.def.go               # Widget interface, Layer, Gated, Awake, UICtx, Size, Align
├── widget.manager.go           # WidgetManager alias + layer-ordered iterators
├── widget.preset_scumm.go      # RegisterDefaultUI
├── widget.preset_radial.go     # RegisterRadialVerbUI
├── widget.preset_rich.go       # RegisterRichUI
├── widget.verb_bar.go          # VerbBar
├── widget.verb_radial.go       # RadialVerbs
├── widget.inventory.go         # InventoryBar
├── widget.status.go            # StatusLine
├── widget.speech.go            # SpeechBubble
├── widget.dialog_box.go        # DialogBox
├── widget.end_card.go          # EndCard
├── widget.cursor.go            # Cursor
├── widget.hotspot_debug.go     # HotspotDebug
├── widget.menu.go              # Menu + MenuItem: the menu system
├── widget.debug_menu.go        # DebugMenuRoot: the menu system's one built-in use
├── widget.panel.go             # Panel
├── widget.top_bar.go           # TopBar
├── widget.character_panel.go   # CharacterPanel + CharStat
├── widget.chat_log.go          # ChatLog + SpeakerStyle
├── widget.letterbox.go         # Letterbox
├── widget.key_binding.go       # KeyBinding
├── widget.text_card.go         # TextCard
├── stage.def.go                # Stage
├── stage.manager.go            # StageManager
│
├── transition.def.go           # Transition
├── transition.fade.go          # Fade, FadeTo
├── widget.hotspot_panel.go     # HotspotPanel
│
├── LICENSE.md                  # MIT
└── README.md                   # this file
```

The companion demo project (`inkwell-demo`) lives in its own repo at
<ssh://git@git.teletypegames.org:2222/games/inkwell-demo>. It consumes this
package via `require git.teletypegames.org/games/inkwell …` in its `go.mod` —
no in-tree coupling.

---

## 23. License

MIT — see [`LICENSE.md`](LICENSE.md). Copyright © 2026 Teletype Games.
