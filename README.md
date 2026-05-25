# pncdsl

A point-and-click adventure game framework for Go, built on top of
[Ebitengine](https://ebitengine.org/). You describe your game with
**declarative struct literals**: every entity (item, scene, character,
dialogue, script, asset, verb, theme, widget) is registered through its
matching `XxxManager.Register(...)` call, and the library handles the rest —
input, rendering, state, dialog trees, cutscenes, HUD.

## What it is

Inspired by the classic LucasArts SCUMM era (Maniac Mansion, Monkey
Island, Day of the Tentacle). The framework doesn't try to be a generic
"engine" — it's a thin Go-level DSL layer on top of Ebitengine that gives
you the fixed skeleton of the adventure genre:

- scenes with hotspots and click zones,
- a fully configurable HUD as a tree of `Widget`s — verb bar, inventory,
  speech, dialog box, end card, cursor, hotspot debug, or your own,
- pluggable verb input modes — built-in `VerbBar` (SCUMM-style) and
  `RadialVerbs` (verb-coin, right-click radial menu),
- inventory, item-on-item and item-on-hotspot interactions,
- dialog trees with conditional choices,
- a script/cutscene system with composite actions (`Seq`, `Par`, `If`, `Wait`),
- swappable color themes (4 presets shipped: SCUMM, Sierra, Paper, Terminal),
- world state (flags, vars), validation.

The approach deliberately collapses to **one single pattern**: every
piece is registered via `Manager.Register(Entity{Name: "..."})`. No
builders, no fluent chains, no registry sprawl.

## Quick start

Prerequisites: Go 1.24+ (for generic type aliases) and a working OpenGL
context.

```bash
git clone <repo>
cd pncdsl
go run .
```

Opens an Ebiten window at 1280×800 (internal resolution 320×200, scaled
4×). No asset files are required — the library generates deterministic
colored placeholders for anything missing under `assets/`. Press **F1**
in-game to toggle the hotspot debug overlay.

## One-minute example

```go
// main.go
package main

import (
    "log"
    "pncdsl/pncdsl"
)

func main() {
    g := pncdsl.NewGame("Sample", 320, 200)

    g.AssetManager.Register(pncdsl.Asset{Name: "bg/kitchen", Path: "assets/bg/kitchen.png", Kind: pncdsl.AssetImage})

    g.ItemManager.Register(pncdsl.Item{
        Name: "key", Sprite: "spr/key", Description: "a rusty key",
    })

    g.SceneManager.Register(pncdsl.Scene{
        Name:       "kitchen",
        Background: "bg/kitchen",
        Hotspots: []pncdsl.Hotspot{
            {
                Name:   "drawer",
                Area:   pncdsl.Rect(40, 80, 60, 40),
                Label:  "drawer",
                OnLook: pncdsl.Say("player", "A drawer. Wonder what's inside?"),
                OnUse:  pncdsl.Seq(pncdsl.Give("key"), pncdsl.Say("player", "A key!")),
            },
        },
    })

    g.CharacterManager.Register(pncdsl.Character{Name: "player", W: 28, H: 62})

    pncdsl.RegisterDefaultUI(g)        // SCUMM-style HUD widgets
    g.UseTheme("classic-scumm")        // or sierra-coin / paper-notebook / terminal-green

    g.StartAt("kitchen")
    if err := pncdsl.Run(g); err != nil {
        log.Fatal(err)
    }
}
```

## The central pattern — `Manager`

Every name-addressable entity goes into the game through one shared type:

```go
type Named interface { GetName() string }

type Manager[T Named] struct { /* ... */ }
func (m *Manager[T]) Register(v T)
func (m *Manager[T]) Get(name string) (T, bool)
func (m *Manager[T]) MustGet(name string) T
func (m *Manager[T]) Has(name string) bool
func (m *Manager[T]) Names() []string       // insertion order
func (m *Manager[T]) SortedNames() []string // alphabetical
func (m *Manager[T]) Remove(name string)
```

Each entity kind gets a generic alias (`ItemManager = Manager[Item]`,
`SceneManager = Manager[Scene]`, `UIManager = Manager[Widget]`, …) so
usage is uniform:

```go
g.ItemManager.Register(pncdsl.Item{Name: "key", ...})
g.SceneManager.Register(pncdsl.Scene{Name: "kitchen", ...})
g.UIManager.Register(&pncdsl.RadialVerbs{Name: "verbs"})
g.ThemeManager.Register(pncdsl.Theme{Name: "my-theme", ...})
```

A duplicate or empty `Name` panics — that's a construction-time bug,
not a runtime error.

Full design rationale: [`PLAN.md`](PLAN.md).

## HUD as a widget tree

The whole HUD is built from `Widget` instances registered into
`g.UIManager`. The interface is minimal — three methods:

```go
type Widget interface {
    GetName() string
    Tick(ctx *UICtx)              // input + state per frame
    Draw(dst *ebiten.Image, ctx *UICtx)
}
```

`Tick` runs in **reverse registration order** so the top widget claims
input first; `Draw` runs in **registration order** so later widgets paint
on top. Built-in widgets the library ships:

| Widget         | Purpose                                                |
|----------------|--------------------------------------------------------|
| `Panel`        | Generic background panel (group/frame other widgets)   |
| `StatusLine`   | Hover hint + 2s flash messages                         |
| `VerbBar`      | SCUMM-style permanent verb grid                        |
| `RadialVerbs`  | Verb-coin: secondary-click radial menu                 |
| `InventoryBar` | Item slots, click to select/inspect                    |
| `SpeechBubble` | Renders the line set by the `Say` action              |
| `DialogBox`    | Active conversation with conditional choices           |
| `EndCard`      | Fullscreen overlay for `ShowEnd`                       |
| `Cursor`       | Mouse cursor + selected-item sprite                    |
| `HotspotDebug` | F1-toggleable hotspot outlines                         |

Convenience helpers:

```go
pncdsl.RegisterDefaultUI(g)     // SCUMM-style: verb bar + inv + dialog + speech + cursor
pncdsl.RegisterRadialVerbUI(g)  // Verb-coin instead of verb bar
```

### Adding your own widget

Anything that implements the `Widget` interface plugs in — a chat panel,
minimap, hotbar, scrollable journal, whatever:

```go
type ChatPanel struct {
    Name   string
    Bounds pncdsl.Rectangle
    lines  []string
}

func (c *ChatPanel) GetName() string                       { return c.Name }
func (c *ChatPanel) Tick(ctx *pncdsl.UICtx)                {}
func (c *ChatPanel) Draw(dst *ebiten.Image, ctx *pncdsl.UICtx) {
    th := ctx.Game.Theme()
    vector.DrawFilledRect(dst, /* … */, th.PanelBG, false)
    for i, ln := range c.lines {
        ctx.Game.DrawText(dst, ln, int(c.Bounds.X)+2, int(c.Bounds.Y)+2+i*14, th.StatusText)
    }
}

g.UIManager.Register(&ChatPanel{Name: "chat", Bounds: pncdsl.Rect(220, 4, 96, 130)})
```

## Themes

Colors live on a registered `Theme` entity, addressed by name. Four
presets are pre-registered by `NewGame`; `classic-scumm` is selected by
default. Switch at any time (including at runtime):

```go
g.UseTheme("paper-notebook")    // light cream + ink
g.UseTheme("terminal-green")    // retro CRT
```

Register your own theme by copying any preset and tweaking the fields:

```go
g.ThemeManager.Register(pncdsl.Theme{
    Name:      "midnight-noir",
    PanelBG:   color.RGBA{6, 6, 12, 255},
    // ... ~25 color fields ...
})
g.UseTheme("midnight-noir")
```

## File layout

```
pncdsl/
├── main.go                  # entry point: pncdsl.Run(domain.Build())
├── pncdsl/                  # the library — theme-prefixed file names
│   ├── core.*.go            # Game, engine, manager, dsl, errors
│   ├── scene.*.go           # Scene, Hotspot, Trigger, Camera, Transition
│   ├── item.*.go            # Item, Inventory
│   ├── actor.*.go           # Character, Animation
│   ├── dialog.*.go          # Dialogue (entity + manager)
│   ├── action.*.go          # Action/Runner, built-in actions, Condition
│   ├── state.*.go           # State, Save (stub)
│   ├── ui.widget.go         # Widget interface, UICtx, Size, Align
│   ├── ui.manager.go        # UIManager alias
│   ├── ui.theme*.go         # Theme + preset themes
│   ├── ui.<widget>.go       # one file per built-in widget
│   ├── ui.defaults.go       # RegisterDefaultUI, RegisterRadialVerbUI
│   ├── asset.*.go           # Registry, lazy load, Audio (stub), Text
│   └── util.*.go            # geometry, timer, log
├── domain/                  # the concrete game — one file per entity
│   ├── game.go              # Build()
│   ├── scene.bedroom.go     # ↔ g.SceneManager.Register(...Name: "bedroom"...)
│   ├── item.key.go          # ↔ g.ItemManager.Register(...Name: "key"...)
│   └── ...
├── PLAN.md                  # detailed design document
├── UIPLAN.md                # widget-system design
├── DEMO.md                  # demo game overview  ◀────────
├── GFX.md                   # asset prompts for image-AI
└── README.md                # this file
```

Files under `domain/` always follow `theme.identifier.go` (e.g.
`scene.kitchen.go`, `character.player.go`) — `ls domain/scene.*`
instantly lists every location.

## The demo

The repo ships with a mini-game (`domain/`) titled **"Morning Coffee"**.
Two scenes, three items, one NPC with a dialog, one intro and one ending
cutscene — just enough to exercise every library feature without the
demo outgrowing the library.

Full walkthrough and file mapping: [`DEMO.md`](DEMO.md).

If you want to generate art for the game, prompts for image AIs are in
[`GFX.md`](GFX.md).

## Current status

| Area                              | State                                                    |
|-----------------------------------|----------------------------------------------------------|
| Manager registries                | ✅ Item, Scene, Character, Dialogue, Script, Asset, Verb, UI, Theme |
| Hotspot + verb interaction        | ✅                                                       |
| Inventory, use-with-item          | ✅ (both `hotspot.OnUseWith` and `item.OnUseWith`)       |
| Dialog tree + conditional choices | ✅ (`Show` condition, `Once` flag)                       |
| `Action`/`Runner` engine          | ✅ Seq / Par / If / Wait / Say / Walk / Give / GoTo / …  |
| Condition library                 | ✅ Flag, HasItem, SelectedItem, InScene, VarEq, Not/And/Or |
| Widget HUD (`Widget` interface)   | ✅ 10 built-in widgets, custom widgets via interface     |
| Verb-coin (`RadialVerbs`)         | ✅ secondary-click radial menu                           |
| Theme system + presets            | ✅ 4 presets, runtime switching                          |
| Scene transitions (fade)          | ✅ on scene change                                       |
| Asset placeholders                | ✅ deterministic color for missing files                 |
| Stylized character placeholder    | ✅ humanoid / quadruped shape when no sprite art         |
| Audio (PlayMusic/PlaySound)       | 🟡 log-only — Ebiten audio backend not wired up yet      |
| Animation (sprite sheet)          | 🟡 `AnimationClip` field exists, rendering not yet       |
| Walkbox + A\*                     | 🟡 straight-line movement only                           |
| Save / Load (JSON)                | 🟡 API in place, implementation stubbed                  |
| Trigger processing                | 🟡 struct exists, engine doesn't run them yet            |
| Custom fonts                      | 🟡 ebitenutil debug font for now                          |

## Testing

```bash
go test ./...
```

`domain/build_test.go` is a headless smoke test: it calls `Build()` and
then `Validate()` to cross-check every name reference between managers.
No Ebiten window is opened.

## Further reading

- [`DEMO.md`](DEMO.md) — walkthrough of the "Morning Coffee" demo

## License

MIT — see [`LICENSE.md`](LICENSE.md). Copyright © 2026 Teletype Games.
