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
5. [Geometry primitives](#5-geometry-primitives)
6. [Entity reference](#6-entity-reference)
   - 6.1 [Asset](#61-asset)
   - 6.2 [Scene](#62-scene)
   - 6.3 [Exit](#63-exit)
   - 6.4 [Hotspot](#64-hotspot)
   - 6.5 [Trigger](#65-trigger)
   - 6.6 [Item](#66-item)
   - 6.7 [Inventory](#67-inventory)
   - 6.8 [Character](#68-character)
   - 6.9 [Dialogue](#69-dialogue)
   - 6.10 [Script](#610-script)
   - 6.11 [Verb](#611-verb)
7. [The action system](#7-the-action-system)
   - 7.1 [Action and Runner](#71-action-and-runner)
   - 7.2 [Status and Ctx](#72-status-and-ctx)
   - 7.3 [Built-in actions](#73-built-in-actions)
   - 7.4 [Custom actions](#74-custom-actions)
8. [Conditions](#8-conditions)
9. [World state](#9-world-state)
10. [The widget system](#10-the-widget-system)
    - 10.1 [Widget interface](#101-widget-interface)
    - 10.2 [Z-order and input consumption](#102-z-order-and-input-consumption)
    - 10.3 [Built-in widgets](#103-built-in-widgets)
    - 10.4 [The clickBlocker contract](#104-the-clickblocker-contract)
    - 10.5 [Registration helpers](#105-registration-helpers)
    - 10.6 [Custom widgets](#106-custom-widgets)
11. [Themes](#11-themes)
12. [Asset pipeline](#12-asset-pipeline)
13. [Audio](#13-audio)
14. [Input](#14-input)
15. [The game loop](#15-the-game-loop)
16. [Scene transitions](#16-scene-transitions)
17. [Validation](#17-validation)
18. [Save and load](#18-save-and-load)
19. [Errors](#19-errors)
20. [Testing](#20-testing)
21. [Project layout](#21-project-layout)
22. [License](#22-license)

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
items, scenes, characters, dialogues, scripts, assets, verbs, widgets and
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
    g := inkwell.NewGame("Sample", 320, 200)

    g.AssetManager.Register(inkwell.Asset{
        Name: "bg/room", Path: "assets/bg/room.png", Kind: inkwell.AssetImage,
    })
    g.CharacterManager.Register(inkwell.Character{Name: "player", W: 28, H: 62})
    g.SceneManager.Register(inkwell.Scene{
        Name:       "room",
        Background: "bg/room",
        Actors:     []inkwell.SceneActor{{CharacterName: "player", At: inkwell.Point{X: 160, Y: 140}}},
        Hotspots: []inkwell.Hotspot{{
            Name:   "door",
            Area:   inkwell.Rect(260, 50, 40, 90),
            Label:  "door",
            OnLook: inkwell.Say("player", "A wooden door."),
        }},
    })

    g.StartAt("room")
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
type ItemManager       = Manager[Item]
type SceneManager      = Manager[Scene]
type CharacterManager  = Manager[Character]
type DialogueManager   = Manager[Dialogue]
type ScriptManager     = Manager[Script]
type AssetManager      = Manager[Asset]
type VerbManager       = Manager[Verb]
type UIManager         = Manager[Widget]
type ThemeManager      = Manager[Theme]
```

### 3.1 Public API

| Method                 | Behaviour                                                  |
|------------------------|------------------------------------------------------------|
| `Register(v T)`        | Adds `v` to the registry. Panics on empty or duplicate `Name`. |
| `Set(v T)`             | Replaces an entry in place, keeping its position; registers it when the name is new. |
| `Get(name) (T, bool)`  | Looks up by name. The boolean is `false` if absent.        |
| `MustGet(name) T`      | Same as `Get`, but panics on missing names.                |
| `Has(name) bool`       | True if the name is registered.                            |
| `Len() int`            | Number of registered entries.                              |
| `Names() []string`     | Returns names in **insertion order** (used for widget Z-order). |
| `All() []T`            | Returns every entry in **insertion order**.                |
| `SortedNames() []string` | Returns names alphabetically.                            |
| `Each(fn func(T))`     | Iterates in insertion order.                               |
| `Remove(name string)`  | Drops a registration; silent no-op if unknown.             |

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
    DialogueManager  *DialogueManager
    ScriptManager    *ScriptManager
    AssetManager     *AssetManager
    VerbManager      *VerbManager
    UIManager        *UIManager
    ThemeManager     *ThemeManager

    // exits
    SceneRect Rectangle          // the part of the window the picture occupies
    ExitLook  func(Exit) Action  // default look response for generated exits
    ExitTake  func(Exit) Action  // default take response for generated exits

    // runtime services
    State     *State
    Inventory *Inventory
    Audio     *AudioPlayer
    Camera    *Camera
    Input     *Input

    // log buffer
    MaxLogLines int
    /* unexported runtime fields */
}
```

### 4.1 Construction

```go
func NewGame(title string, w, h int) *Game
```

Initialises every manager (empty), registers the default SCUMM verb set
(`look`, `use`, `talk`, `take`), installs all four preset themes, and
selects `classic-scumm` as the active theme. Widgets are **not**
auto-registered — call `RegisterDefaultUI(g)` (or one of its siblings)
explicitly, or let `Run` install the default set if `UIManager` is empty.

### 4.2 Lifecycle hooks

```go
func (g *Game) StartAt(name string) *Game     // entry scene
func (g *Game) OnStart(a Action) *Game        // action run after the scene's OnEnter
func (g *Game) Validate() error               // cross-check name references
func (g *Game) Run() error                    // same as inkwell.Run(g)
```

`StartAt` and `OnStart` return `*Game` so they chain at the end of `Build`.

### 4.3 Theme accessors

```go
func (g *Game) Theme() Theme                  // active theme value
func (g *Game) UseTheme(name string)          // switch theme (runtime ok)
```

The active theme is addressed by name, not by pointer — switching is just
`g.UseTheme("paper-notebook")` and the next frame uses the new colors.

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

func (g *Game) FlashLine(text string)         // ~2s status banner
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
func (g *Game) CharacterInScene(name string) bool
```

`CurrentScene` is where the player is, empty before the first scene is
entered; `PreviousScene` is the one before it, which is where
[`Back()`](#73-built-in-actions) leads. Both survive a save.

`SceneArea` is `SceneRect`, or the whole window when it was never set. It is
what [exit strips](#63-exit) are measured against.

`SceneHotspots` is everything clickable in a scene: its own `Hotspots` first,
then one per `Exit`. Authored hotspots come first because the engine takes the
first area that contains the click, so a painted thing beats the edge strip an
exit sits on wherever the two overlap. The expansion is cached per scene.

`CharacterInScene` is used by `CharacterPanel` to auto-hide when its character
isn't an actor in the current scene.

### 4.7 Text drawing hook

```go
func (g *Game) DrawText(dst *ebiten.Image, s string, x, y int, c color.Color)
```

Single indirection so widgets can hand a `color.Color` today and the
library can swap to `text/v2` later without changing call sites.

---

## 5. Geometry primitives

```go
// inkwell/util.geometry.go

type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point
func (p Point) Sub(q Point) Point
func (p Point) Dist(q Point) float64

type Shape interface {
    Contains(p Point) bool
    Bounds() Rectangle
}

type Rectangle struct{ X, Y, W, H float64 }

func Rect(x, y, w, h float64) Rectangle
func (r Rectangle) Contains(p Point) bool
func (r Rectangle) Bounds() Rectangle
func (r Rectangle) Center() Point

type Polygon struct{ Points []Point }

func Poly(pts ...Point) Polygon
func (g Polygon) Contains(pt Point) bool      // ray-casting
func (g Polygon) Bounds() Rectangle           // axis-aligned bbox
```

Both `Rectangle` and `Polygon` satisfy `Shape` and can be used as a
`Hotspot.Area`. Internally, hotspots are hit-tested against `Shape.Contains`
each frame.

---

## 6. Entity reference

### 6.1 Asset

```go
// inkwell/asset.def.go

type AssetKind int
const (
    AssetImage AssetKind = iota
    AssetAudio
    AssetFont
)

type Asset struct {
    Name string  // logical id used everywhere ("bg/kitchen")
    Path string  // disk path ("assets/bg/kitchen.png")
    Kind AssetKind
}
```

Assets are **lazy**: `Register` only stores the spec. The file is opened on
the first `loadedAssets.image(...)` request. Missing files fall back to a
deterministic colored placeholder so the game runs without any art —
useful for prototyping or CI.

```go
g.AssetManager.Register(inkwell.Asset{
    Name: "bg/kitchen",
    Path: "assets/bg/kitchen.png",
    Kind: inkwell.AssetImage,
})
```

### 6.2 Scene

```go
// inkwell/scene.def.go

type Scene struct {
    Name       string
    Title      string         // optional human-readable name (TopBar widget reads this)
    Background string         // Asset.Name
    Music      string         // Asset.Name (optional)
    Hotspots   []Hotspot
    Exits      []Exit         // connections to other scenes
    Walkboxes  []Polygon
    Triggers   []Trigger
    Actors     []SceneActor
    OnEnter    Action
    OnLeave    Action
}

type SceneActor struct {
    CharacterName string
    At            Point
}
```

`OnEnter` is queued the moment the scene becomes current (after a fade-in);
`OnLeave` runs on the transition out. `Actors` lists which registered
characters are placed in the scene and where their feet start.

`Walkboxes` constrain character movement. When non-empty, `Walk(...)`
routes through a BFS over the polygon adjacency graph (polygons that
share an edge are neighbours), with the midpoint of each shared edge
used as a waypoint. Destinations outside every walkbox are clipped to
the nearest boundary. With no walkboxes the character walks in a
straight line — see [§6.8](#68-character) and [§15.4](#154-character-movement).

`Triggers` fire on the rising edge of their `When` condition. The
engine samples each trigger once per idle frame; see [§6.5](#65-trigger).

`Exits` are the scene's connections to other scenes, declared as data. The
engine turns each one into a hotspot — see [§6.3](#63-exit).

### 6.3 Exit

```go
// inkwell/scene.exit.go

type Exit struct {
    To      string      // target scene; empty means "back the way you came"
    Label   string      // what the status line calls it
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

func ExitName(to string) string
```

A connection belongs to the scene it leads out of, so it is written in that
scene's own literal rather than in a map kept somewhere else:

```go
Scene{
    Name: "alley",
    Exits: []Exit{
        {To: "noodle_house", Label: "back out to the street", Side: ExitLeft},
        {To: "", Label: "the fire escape", Side: ExitBack,
         Needs: "has_ladder", Blocked: Say("paul", "Can't reach it.")},
    },
}
```

The engine expands each exit into a `Hotspot` with `Cursor: CursorExit`, the
exit's `Label`, `OnUse` bound to `GoTo(To)` — or [`Back()`](#73-built-in-actions)
when `To` is empty — and, when `Needs` is set, the whole travel wrapped in
`If(Flag(Needs), travel, Blocked)`. The exit is visible either way: a locked
door still says it is a door.

`OnLook` and `OnTake` are usually not written per exit. `Game.ExitLook` and
`Game.ExitTake` supply them for every exit in the game, so the phrasing lives
in one place:

```go
g.ExitLook = func(e Exit) Action { return Say("paul", "That way: "+e.Label+".") }
g.ExitTake = func(e Exit) Action { return Say("paul", "It's a way out, not a thing.") }
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
[`Validate`](#17-validation) rejects an exit that names a scene nobody
registered.

### 6.4 Hotspot

```go
// inkwell/scene.hotspot.go

type Hotspot struct {
    Name   string
    Area   Shape           // Rectangle or Polygon
    Label  string          // hover hint
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

### 6.5 Trigger

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

### 6.6 Item

```go
// inkwell/item.def.go

type Item struct {
    Name        string
    Sprite      string                  // Asset.Name
    Description string

    OnUseSelf Action
    // targetName → action. Target can be a hotspot Name or another item Name.
    OnUseWith map[string]Action
}
```

Items live in the `ItemManager`. When the player picks one up
(`Give("key")`) the item is appended to the `Inventory`. Clicking a
hotspot with an item selected resolves the action via, in order:
1. `hotspot.OnUseWith[item.Name]`
2. `item.OnUseWith[hotspot.Name]`
3. Otherwise the engine flashes "Nem ehhez." and deselects.

### 6.7 Inventory

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

### 6.8 Character

```go
// inkwell/actor.def.go

type Character struct {
    Name        string
    Sprite      string                       // Asset.Name (placeholder if missing)
    Animations  map[string]AnimationClip
    Speed       float64
    SpeechColor color.Color
    Start       Point
    W, H        float64                      // size hint for placeholder
}

type AnimationClip struct {
    Frames    []Rectangle
    FrameTime float64
    Loop      bool
}
```

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
computed by [walkbox routing](#62-scene).

### 6.9 Dialogue

```go
// inkwell/dialog.def.go

type Dialogue struct {
    Name  string
    Start string                        // initial node Name; "" = first in slice
    Nodes []DialogueNode
}

type DialogueNode struct {
    Name    string
    Lines   []DialogueLine              // shown one-at-a-time on click
    Choices []DialogueChoice
}

type DialogueLine struct {
    Speaker string                      // Character.Name
    Text    string
}

type DialogueChoice struct {
    Text    string
    Show    Condition                   // nil = always visible
    Once    bool                        // hide after first pick
    Actions []Action                    // run when chosen
}

func (d Dialogue) Node(name string) (DialogueNode, bool)
```

The dialog flow:
1. `RunDialogue("name")` action starts the dialogue; the runner keeps
   returning `StatusRunning` until the dialog closes.
2. The `DialogBox` widget renders lines one at a time; click advances.
3. After the last line, visible `Choices` are shown.
4. Picking a choice runs `Seq(choice.Actions...)`. `Once` choices remember
   that they fired via `State.NoteTalked`.
5. `EndDialogue()` closes the conversation; `GotoNode("other")` jumps to
   another node in the same dialogue.

### 6.10 Script

```go
// inkwell/action.script.go

type Script struct {
    Name    string
    Actions Action                      // typically Seq(...)
}
```

A `Script` is just a named composite action — useful when you want to
reuse a cutscene (intro, victory, transition) from multiple call sites.
Fire one with `RunScript("name")`.

### 6.11 Verb

```go
// inkwell/ui.verb.go

type Verb struct {
    Name    string                      // canonical id, e.g. "look"
    Label   string                      // display label, e.g. "Nézd"
    Default Action                      // run when a hotspot has no handler
}
```

`NewGame` pre-registers the four SCUMM verbs (`look`, `use`, `talk`,
`take`). Add your own:

```go
g.VerbManager.Register(inkwell.Verb{
    Name: "push", Label: "Lökd",
    Default: inkwell.Say("player", "Nem mozdul."),
})
```

The verb-bar and verb-wheel widgets read `VerbManager.Names()` every
frame, so adding a verb at runtime is enough to make it appear.

---

## 7. The action system

### 7.1 Action and Runner

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

### 7.2 Status and Ctx

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

### 7.3 Built-in actions

| Constructor                              | Behaviour                                                                 |
|------------------------------------------|---------------------------------------------------------------------------|
| `Seq(a ...Action) Action`                | Run children in order. Nested `Seq` are flattened. `nil` children are ignored. |
| `Par(a ...Action) Action`                | Run children simultaneously; completes when every child reports `Done`.   |
| `If(c Condition, then Action, els ...Action) Action` | Run `then` if `c` evaluates true; otherwise `Seq(els...)` (if any).|
| `Wait(seconds float64) Action`           | Block the runner for `seconds`.                                            |
| `Say(speaker, text string) Action`       | Show the line above the speaker (`SpeechBubble`), append to chat-log. Click-to-skip. Duration scales with text length, 1.2s floor. |
| `GoTo(scene string) Action`              | Switch the current scene via a fade transition.                            |
| `Back() Action`                          | Return to `PreviousScene()`. No-op when there is nowhere to go back to.    |
| `Walk(character string, to Point) Action`| Move a character to `to` at `Character.Speed`. Returns when arrived.       |
| `Give(item string) Action`               | Add `item` to inventory.                                                   |
| `TakeAway(item string) Action`           | Remove `item` from inventory.                                              |
| `RequireItem(item string) Action`        | If the item is not in inventory, flash "Ehhez kell egy …" and fail.        |
| `SetFlag(name string) Action`            | `State.SetFlag(name)`.                                                     |
| `ClearFlag(name string) Action`          | `State.ClearFlag(name)`.                                                   |
| `SetVar(name string, v any) Action`      | `State.SetVar(name, v)`.                                                   |
| `PlayMusic(asset string) Action`         | Hand off to `AudioPlayer.PlayMusic`. (Currently logged only.)              |
| `StopMusic() Action`                     | `AudioPlayer.StopMusic`.                                                   |
| `PlaySound(asset string) Action`         | `AudioPlayer.PlaySound`.                                                   |
| `RunDialogue(name string) Action`        | Start a registered `Dialogue`; runner blocks until the dialog closes.       |
| `EndDialogue() Action`                   | Close the active dialog.                                                   |
| `GotoNode(node string) Action`           | Jump to a node within the active dialog.                                   |
| `RunScript(name string) Action`          | Run a registered `Script`'s `Actions`.                                     |
| `ShowEnd(text string) Action`            | Display the end-card overlay; runner never finishes.                       |
| `Custom(fn func(*Ctx) Status) Action`    | Escape hatch — wrap an arbitrary callback as an action.                    |

### 7.4 Custom actions

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

## 8. Conditions

```go
// inkwell/action.condition.go

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
| `HasItem(name string) Condition`       | `Inventory.Has(name)`                            |
| `SelectedItem(name string) Condition`  | `Inventory.Selected() == name`                   |
| `InScene(name string) Condition`       | The current `Ctx.Scene` matches.                 |
| `VarEq(name string, v any) Condition`  | `State.Var(name) == v` (Go equality).            |

Conditions are stateless — they are evaluated lazily, including inside
`If` and `DialogueChoice.Show`.

---

## 9. World state

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

`Visited` and `Talked` underpin the "Once" semantics in dialogues and
could be used by triggers once those land.

---

## 10. The widget system

The HUD is a tree of `Widget`s, each registered into `g.UIManager`. The
library ships built-in widgets for every classic adventure UI piece;
domain code can register arbitrary new ones — chat panels, minimaps,
hotbars — without touching the library.

### 10.1 Widget interface

```go
// inkwell/ui.widget.go

type Widget interface {
    Named
    Tick(ctx *UICtx)
    Draw(dst *ebiten.Image, ctx *UICtx)
}

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

### 10.2 Z-order and input consumption

- **`Tick` runs in reverse registration order.** Widgets registered later
  (drawn on top) get the click first. Each widget calls
  `ctx.Game.Input.ConsumeLeft()` / `ConsumeRight()` to claim the event;
  later widgets see `LeftClicked() == false`.
- **`Draw` runs in registration order** — registered last → painted on top.
- The `Cursor` widget is registered last by convention so it always wins
  on visual layer (and effectively never claims clicks).

After all widgets ticked, the engine offers the (possibly consumed) click
to `handleSceneInput`, which is where hotspot interactions live. If a
widget already consumed the click, `handleSceneInput` becomes a no-op for
that frame.

### 10.3 Built-in widgets

Each widget lives in its own `inkwell/ui.<name>.go` file. Fields with a
zero value fall back to a sensible default.

#### `Panel` (`ui.panel.go`)

A generic colored rectangle, optionally with a stroke border. Useful as a
backdrop for grouping other widgets.

```go
type Panel struct {
    Name        string
    Bounds      Rectangle
    BG          color.Color    // nil → Theme.PanelBG
    Border      int            // 0 = no border
    BorderColor color.Color    // nil → Theme.DialogBorder
}
```

#### `StatusLine` (`ui.status.go`)

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

#### `VerbBar` (`ui.verb_bar.go`)

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

#### `RadialVerbs` (`ui.verb_radial.go`)

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
    Labels map[string]string
}
```

In **trigger mode** (default), the wheel is invisible until the trigger
mouse button is pressed somewhere not covered by a `clickBlocker`. The
center latches to the click point. Left-click on a slice picks the verb;
trigger-click again closes without picking.

In **always-visible mode**, the wheel is drawn permanently at `Center`
and left-click on a slice picks the verb. Labels are positioned at
`0.62 × Radius` from the center so short overrides stay inside the rim.

#### `InventoryBar` (`ui.inventory.go`)

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

#### `SpeechBubble` (`ui.speech.go`)

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

#### `DialogBox` (`ui.dialog_box.go`)

Dormant unless `Game.dialogueActive()`. While active, **consumes every
click** — the rest of the HUD and the scene are inert. The widget renders
either the current line of `runtimeDialog.Node.Lines` (click → next line)
or the visible `Choices` once the lines are exhausted.

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

#### `EndCard` (`ui.end_card.go`)

A fullscreen overlay shown when `Game.showEndCard(text)` is called (which
the `ShowEnd` action does). Consumes every click so nothing underneath
reacts.

```go
type EndCard struct{ Name string }
```

#### `Cursor` (`ui.cursor.go`)

The crosshair cursor (color from `Theme.CursorColor`). When the inventory
has a selected item with a non-placeholder sprite, that sprite is drawn
at the cursor position instead.

```go
type Cursor struct{ Name string }
```

#### `HotspotDebug` (`ui.hotspot_debug.go`)

Stroke-outlines every hotspot in the current scene. F1 toggles it at
runtime (override via `ToggleKey`).

```go
type HotspotDebug struct {
    Name      string
    Enabled   bool
    ToggleKey ebiten.Key       // 0 → ebiten.KeyF1
}
```

#### `TopBar` (`ui.top_bar.go`)

A horizontal strip across the top with three sections:

- **Left:** `LeftText` override or `Scene.Title` (falling back to `Name`).
- **Center:** score string built from `State.Var(ScoreVar)`. With
  `ScoreMax > 0`, formatted as `"Score: X/MAX"`, else `"Score: X"`.
- **Right:** time string from `State.Var(TimeVar)`.

```go
type TopBar struct {
    Name     string
    Height   int                  // default 12
    LeftText string                // overrides scene title
    ScoreVar string                // State.Var key; "" = no score
    ScoreMax int
    TimeVar  string                // State.Var key; "" = no time
}
```

Empty fields are skipped, so the widget gracefully degrades to two- or
one-section layouts.

#### `CharacterPanel` (`ui.character_panel.go`)

A floating info card showing one character's portrait, role badge and
custom stat rows. **Auto-hides** when the character isn't an actor in the
current scene — same registration works across all scenes.

```go
type CharStat struct {
    Label  string
    VarKey string                  // State.Var(VarKey) is fmt-printed
}

type CharacterPanel struct {
    Name      string
    Bounds    Rectangle
    Character string                // Character.Name
    Title     string                // role badge ("PLAYER", "NPC", "GUIDE", …)
    Stats     []CharStat
}
```

The portrait reuses the placeholder rendering when no sprite is on disk
(humanoid vs. quadruped picked from `W` vs. `H`).

#### `ChatLog` (`ui.chat_log.go`)

Renders the `Game.Messages()` ring buffer with per-kind colors —
`LogAction` lines use `Theme.ChatLogPrompt`, `LogResponse` uses
`Theme.ChatLogResponse`, `LogSystem` uses `Theme.ChatLogSystem`. Newest
message at the bottom; older messages scroll out of view as the buffer
fills.

```go
type ChatLog struct {
    Name       string
    Bounds     Rectangle
    LineHeight int                 // default 10
    Padding    int                 // default 3
    ShowBorder bool                // draws a 1px border in CharacterPanelBorder
}
```

Set `Game.MaxLogLines` to bound the buffer (recommended ~64 for a small
ChatLog).

### 10.4 The `clickBlocker` contract

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
visible). Implement it on your own widgets to mark them as solid input
zones for the radial menu.

### 10.5 Registration helpers

```go
// inkwell/ui.defaults.go

func RegisterDefaultUI(g *Game)                              // SCUMM-style verb bar + inventory
func RegisterRadialVerbUI(g *Game)                            // verb-coin (right-click) + wider inventory
func RegisterRichUI(g *Game, playerName, npcName string)      // top bar + character panels + chat log + verb wheel
```

Pass `""` for any character name to skip its `CharacterPanel`. All three
helpers register in conventional Z-order (debug overlay back, cursor
front).

### 10.6 Custom widgets

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

g.UIManager.Register(&Minimap{
    Name:   "minimap",
    Bounds: inkwell.Rect(220, 4, 96, 56),
})
```

---

## 11. Themes

```go
// inkwell/ui.theme.go

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

### 11.1 Preset themes

`NewGame` automatically calls `RegisterPresetThemes(g)` and
`UseTheme("classic-scumm")`. The four shipped presets:

| Name              | Mood                                                        |
|-------------------|-------------------------------------------------------------|
| `classic-scumm`   | Default SCUMM-era dark blue + warm amber accents.          |
| `sierra-coin`     | Darker purple + orange; designed for the verb-coin layout. |
| `paper-notebook`  | Cream paper + ink + a touch of sepia for the highlights.   |
| `terminal-green`  | Retro phosphor green on black.                              |

### 11.2 Custom themes

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

## 12. Asset pipeline

`AssetManager.Register` only stores the spec (`Name`, `Path`, `Kind`). The
first request for an image goes through `loadedAssets.image(name)`:

1. Look up the file from `Asset.Path`.
2. If the file is missing or unreadable, generate a deterministic
   placeholder color from `fnv32a(name)`. Mark the entry as a placeholder.
3. Cache the decoded `*ebiten.Image` for subsequent calls.

```go
// internal:
type loadedAssets struct {
    images       map[string]*ebiten.Image
    placeholders map[string]bool
}

func (la *loadedAssets) image(am *AssetManager, name string) *ebiten.Image
func (la *loadedAssets) isPlaceholder(name string) bool
```

`isPlaceholder` is used by widgets to decide whether to render a stylised
character placeholder (humanoid / quadruped) or the real sprite.

Supported image formats: PNG and JPEG (registered via blank imports in
`asset.manager.go`). Audio assets (`AssetAudio`) are decoded on demand
by `AudioPlayer` — see [§13](#13-audio). Fonts (`AssetFont`) share the
same registry but are not loaded by the library at this milestone.

---

## 13. Audio

```go
// inkwell/asset.audio.go

type AudioPlayer struct{ /* unexported */ }

func NewAudioPlayer() *AudioPlayer
func (a *AudioPlayer) PlayMusic(name string)
func (a *AudioPlayer) StopMusic()
func (a *AudioPlayer) PlaySound(name string)
```

Backed by Ebiten's `audio` package. The shared `*audio.Context` is
created lazily on the first `PlayMusic` / `PlaySound`, at a sample rate
of 48 kHz; existing audio contexts (e.g. created by a host process) are
detected via `audio.CurrentContext` and reused. The `Asset.Kind` must
be `AssetAudio` for the player to consider it.

- **Music** streams from disk and loops infinitely. `PlayMusic(name)`
  on the currently-playing track no-ops, so it's safe to call from a
  scene's `OnEnter` on every visit.
- **Sound effects** are decoded once, cached as raw PCM, and replayed
  via `NewPlayerFromBytes` for near-zero-latency triggering. Finished
  players are pruned on the next `PlaySound` call.

Supported formats: **WAV** (`.wav`), **OGG Vorbis** (`.ogg`), and
**MP3** (`.mp3`), picked by file extension. Anything else fails to
decode and the asset is flagged so subsequent plays no-op.

Failure modes (missing file, unsupported codec, no audio device) all
degrade to a debug-log line; the game keeps running with no sound. The
`PlayMusic` / `PlaySound` / `StopMusic` actions route through the same
methods, so a domain authored against the original stub still works.

---

## 14. Input

```go
// inkwell/input.def.go

type Input struct{ /* unexported */ }

func (i *Input) Pos() (int, int)
func (i *Input) Point() Point
func (i *Input) LeftClicked() bool       // just-pressed AND not consumed
func (i *Input) RightClicked() bool
func (i *Input) ConsumeLeft()
func (i *Input) ConsumeRight()
```

The engine `poll`s once per frame, snapshotting cursor position and the
just-pressed state for both buttons. Widgets and the engine use the
**consume-on-use** pattern: one click is authoritative — whoever consumes
first wins, later checks see `LeftClicked() == false`.

`Game.Input` is exposed so widgets and custom actions can read it
directly.

---

## 15. The game loop

```go
// inkwell/core.dsl.go

func Run(g *Game) error                          // toplevel — same as g.Run()
```

`Run` does, in order:

1. `g.Validate()` — cross-check name references between managers.
2. If `UIManager` is empty, call `RegisterDefaultUI(g)`.
3. Place the start scene directly (no transition), bump
   `State.NoteVisit`, position registered actors, kick off music.
4. Compose `Seq(scene.OnEnter, game.OnStart)` and queue it as the initial
   action — the first script tick runs both in order.
5. `ebiten.SetWindowSize(Width*4, Height*4)`,
   `ebiten.SetWindowTitle(g.Title)`, then `ebiten.RunGame(&engine{g})`.

### 15.1 `engine.Update`

```
poll input
update transition
if transition fading out → return

if scriptRunner != nil:
    tick the runner
    tick characters
    return

clear hoverLabel
for w in reversed(UIManager):
    w.Tick(uictx)        // widgets consume input top-down

handleSceneInput()       // hotspot resolution + right-click reset
tick characters
```

### 15.2 `engine.Draw`

```
fill Theme.SceneBackdrop (if any)
draw scene background image
for c in characters sorted by Y:
    drawCharacter(c)
for w in UIManager (registration order):
    w.Draw(screen, uictx)
draw transition overlay
```

### 15.3 `handleSceneInput`

Runs only after every widget had a chance. The flow:

1. Set hover label from the hotspot under the cursor (if any).
2. On right-click: deselect the held item if any, otherwise reset
   `selectedVerb` to `"look"`.
3. On left-click:
   - If an item is selected, try `hotspot.OnUseWith[item]`, then
     `item.OnUseWith[hotspot]`. Fail with a "Nem ehhez." flash.
   - Otherwise look up `hotspot.handler(selectedVerb)`. Fall back to the
     verb's `Default` action. Fail with "Semmi említésre méltó."
4. On every successful click, push a `LogAction` line so the `ChatLog`
   widget can display it.

### 15.4 Character movement

`Game.tickCharacters` runs every frame. For each moving character it
steps the head of its waypoint queue at `def.Speed` pixels per second
(default 60); on arrival, the waypoint is popped and the character
proceeds to the next one. The character becomes idle when the queue
empties.

The waypoint queue is built when `Walk(name, to)` runs — see
[§6.2](#62-scene) and [`scene.path.go`](#21-project-layout) for the
walkbox routing. The `Walk` runner returns `StatusRunning` until
`characterMoving(name)` reports false.

Each frame also advances the character's animation clock via
`tickAnimation`, which auto-selects `"walk"` (when `moving == true`)
or `"idle"` from `Character.Animations`. The renderer uses the active
clip's source rectangle to blit a sprite-sheet frame, or falls back to
the whole-sprite / placeholder draw when no clip applies.

---

## 16. Scene transitions

```go
// inkwell/scene.transition.go (unexported)

type transition struct{ /* unexported */ }
```

A scene change calls `Game.changeScene(name)`, which starts a fade-to-
black overlay (`duration = 0.25s`). At the midpoint:

1. `prevScene.OnLeave` (if any) is queued.
2. `currentScene` is set to the new name.
3. `State.NoteVisit(name)` bumps the visit counter.
4. Actors are placed at their `SceneActor.At` positions.
5. The new scene's `Music` plays.
6. `newScene.OnEnter` is queued.

The fade-from-black plays after the midpoint. During the fade-out half,
`engine.Update` returns early — input is frozen briefly.

The initial scene (from `StartAt`) is set **without** a transition so the
intro script runs immediately.

---

## 17. Validation

```go
func (g *Game) Validate() error
```

Run by `Run` before entering the Ebiten loop, and again from your tests.
It cross-checks:

- `StartAt` points to a registered scene.
- Every scene has a non-empty `Background` referencing a registered
  `Asset`.
- Optional `Scene.Music` (if non-empty) references a registered `Asset`.
- Every `SceneActor.CharacterName` is a registered character.
- Every `Scene.Exit.To` names a registered scene (empty is allowed — it
  means "back the way you came").
- An active theme is selected and registered.

Returns the first error wrapping one of the `Err...` sentinels (so callers
can `errors.Is` against them). Duplicate-name registration is caught
eagerly at `Register` time via panic.

---

## 18. Save and load

```go
// inkwell/state.save.go

func (g *Game) Save(slot int) error
func (g *Game) Load(slot int) error
```

Slots are written as JSON files under `g.SaveDir` (default `saves/`,
relative to the working directory), one file per slot named
`slot<N>.json`. The save captures the **mutable runtime state**:

- `currentScene` and `previousScene`, the active verb, and the active theme.
- Every character's position, target, and moving flag.
- The full `State`: flags, vars, visited and talked counters.
- The inventory item list plus the currently selected slot.

Static registrations (items, scenes, characters, dialogues, scripts,
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

## 19. Errors

```go
// inkwell/core.errors.go

var (
    ErrUnknownAsset           = errors.New("inkwell: unknown asset")
    ErrUnknownScene           = errors.New("inkwell: unknown scene")
    ErrUnknownItem            = errors.New("inkwell: unknown item")
    ErrUnknownCharacter       = errors.New("inkwell: unknown character")
    ErrUnknownDialogue        = errors.New("inkwell: unknown dialogue")
    ErrUnknownDialogueNode    = errors.New("inkwell: unknown dialogue node")
    ErrUnknownScript          = errors.New("inkwell: unknown script")
    ErrUnknownVerb            = errors.New("inkwell: unknown verb")
    ErrDuplicateName          = errors.New("inkwell: duplicate name")
    ErrNoStartScene           = errors.New("inkwell: StartAt not set or unknown scene")
    ErrSceneMissingBackground = errors.New("inkwell: scene has no background")
)
```

All wrap-able with `fmt.Errorf("%w: …", err, ctx)` and checkable with
`errors.Is`. `Manager.Register` and `MustGet` use `panic` instead of
returning these errors because they signal construction-time bugs that
should crash loudly in development.

---

## 20. Testing

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

### 20.1 Useful runtime debug switches

```go
inkwell.DebugLog = true                 // logf output to stderr (script queue, audio, scene change)
&inkwell.HotspotDebug{Enabled: true}    // start with the F1 overlay on
```

---

## 21. Project layout

All library source files live at the repo root; no nested package
directory. The naming convention is `theme.identifier.go` — `ls core.*`
or `ls ui.*` instantly groups related sources.

```
inkwell/                            # module git.teletypegames.org/games/inkwell
├── go.mod
├── go.sum
│
├── core.doc.go                    # package docs
├── core.manager.go                # Manager[T Named], Named, TypeLabel
├── core.game.go                   # Game aggregate, runtime state, helpers
├── core.engine.go                 # ebiten.Game adapter (Update/Draw/Layout)
├── core.dsl.go                    # Run()
├── core.errors.go                 # sentinel errors
│
├── util.geometry.go               # Point, Rectangle, Polygon, Shape
├── util.timer.go                  # Timer helper
├── util.log.go                    # DebugLog + logf
│
├── asset.def.go                   # Asset, AssetKind
├── asset.manager.go               # AssetManager alias + lazy loader
├── asset.audio.go                 # AudioPlayer (WAV/OGG/MP3, looping music)
├── asset.text.go                  # drawText / wrapText helpers
│
├── scene.def.go                   # Scene, SceneActor
├── scene.manager.go               # SceneManager alias
├── scene.hotspot.go               # Hotspot, CursorKind
├── scene.exit.go                  # Exit, ExitSide + edge-strip geometry
├── scene.trigger.go               # Trigger + rising-edge engine sweep
├── scene.path.go                  # walkbox routing (BFS over polygon adjacency)
├── scene.transition.go            # fade-to-black overlay (internal)
├── scene.camera.go                # Camera (stub identity transform)
│
├── item.def.go                    # Item
├── item.manager.go                # ItemManager alias
├── item.inventory.go              # Inventory
│
├── actor.def.go                   # Character
├── actor.manager.go               # CharacterManager alias
├── actor.animation.go             # AnimationClip + tickAnimation
│
├── dialog.def.go                  # Dialogue, DialogueNode, DialogueChoice
├── dialog.manager.go              # DialogueManager alias
│
├── action.def.go                  # Action, Runner, Ctx, Status + every built-in action
├── action.condition.go            # Condition interface + combinators
├── action.script.go               # Script entity
├── action.manager.go              # ScriptManager alias
│
├── state.def.go                   # State (flags, vars, visited, talked)
├── state.save.go                  # Save/Load JSON persistence
│
├── input.def.go                   # Input (consume-on-use)
│
├── ui.widget.go                   # Widget interface, UICtx, Size, Align
├── ui.manager.go                  # UIManager alias + reversed/ordered iterators
├── ui.theme.go                    # Theme + ThemeManager
├── ui.theme_presets.go            # 4 preset themes
├── ui.defaults.go                 # RegisterDefaultUI/RadialVerbUI/RichUI
├── ui.verb.go                     # Verb + VerbManager
├── ui.verb_bar.go                 # VerbBar widget
├── ui.verb_radial.go              # RadialVerbs widget
├── ui.inventory.go                # InventoryBar widget
├── ui.status.go                   # StatusLine widget
├── ui.speech.go                   # SpeechBubble widget
├── ui.dialog_box.go               # DialogBox widget
├── ui.end_card.go                 # EndCard widget
├── ui.cursor.go                   # Cursor widget
├── ui.hotspot_debug.go            # HotspotDebug widget
├── ui.panel.go                    # Panel widget
├── ui.top_bar.go                  # TopBar widget
├── ui.character_panel.go          # CharacterPanel widget + CharStat
├── ui.chat_log.go                 # ChatLog widget
│
├── LICENSE.md                     # MIT
└── README.md                      # this file
```

The companion demo project (`inkwell-demo`) lives in its own repo at
<ssh://git@git.teletypegames.org:2222/games/inkwell-demo>. It consumes this
package via `require git.teletypegames.org/games/inkwell …` in its `go.mod` —
no in-tree coupling.

---

## 22. License

MIT — see [`LICENSE.md`](LICENSE.md). Copyright © 2026 Teletype Games.
