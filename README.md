# pncdsl

A point-and-click adventure game framework for Go, built on top of
[Ebitengine](https://ebitengine.org/). You describe your game with
**declarative struct literals**: every entity (item, scene, character,
dialogue, script, asset, verb) is registered through its matching
`XxxManager.Register(...)` call, and the library handles the rest —
input, rendering, state, dialog trees, cutscenes.

## What it is

Inspired by the classic LucasArts SCUMM era (Maniac Mansion, Monkey
Island, Day of the Tentacle). The framework doesn't try to be a generic
"engine" — it's a thin Go-level DSL layer on top of Ebitengine that gives
you the fixed skeleton of the adventure genre:

- scenes with hotspots and click zones,
- a verb-bar UI: Look / Use / Talk / Take,
- inventory, item-on-item and item-on-hotspot interactions,
- dialog trees with conditional choices,
- a script/cutscene system with composite actions (`Seq`, `Par`, `If`, `Wait`),
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
colored placeholders for anything missing under `assets/`.

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
func (m *Manager[T]) Names() []string
```

Each entity kind gets a generic alias (`ItemManager = Manager[Item]`,
`SceneManager = Manager[Scene]`, …) so usage is uniform:

```go
g.ItemManager.Register(pncdsl.Item{Name: "key", ...})
g.SceneManager.Register(pncdsl.Scene{Name: "kitchen", ...})
g.AssetManager.Register(pncdsl.Asset{Name: "bg/kitchen", Path: "..."})
```

A duplicate or empty `Name` panics — that's a construction-time bug,
not a runtime error.

Full design rationale: [`PLAN.md`](PLAN.md).

## File layout

```
pncdsl/
├── main.go                  # entry point: pncdsl.Run(domain.Build())
├── pncdsl/                  # the library — theme-prefixed file names
│   ├── core.*.go            # Game, engine, manager, dsl, errors
│   ├── scene.*.go           # Scene, Hotspot, Trigger, Camera, Transition
│   ├── item.*.go            # Item, Inventory
│   ├── actor.*.go           # Character, Animation
│   ├── dialog.*.go          # Dialogue, DialogBox
│   ├── action.*.go          # Action/Runner, built-in actions, Condition
│   ├── state.*.go           # State, Save (stub)
│   ├── ui.*.go              # VerbBar, Cursor, SpeechBubble
│   ├── asset.*.go           # Registry, lazy load, Audio (stub), Text
│   └── util.*.go            # geometry, timer, log
├── domain/                  # the concrete game — one file per entity
│   ├── game.go              # Build()
│   ├── scene.bedroom.go     # ↔ g.SceneManager.Register(...Name: "bedroom"...)
│   ├── item.key.go          # ↔ g.ItemManager.Register(...Name: "key"...)
│   └── ...
├── PLAN.md                  # detailed design document
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

| Area                          | State                                                    |
|-------------------------------|----------------------------------------------------------|
| Manager registries            | ✅ complete (Item, Scene, Character, Dialogue, Script, Asset, Verb) |
| Hotspot + verb interaction    | ✅                                                       |
| Inventory, use-with-item      | ✅ (both `hotspot.OnUseWith` and `item.OnUseWith`)       |
| Dialog tree + conditional choices | ✅ (`Show` condition, `Once` flag)                    |
| `Action`/`Runner` engine      | ✅ Seq / Par / If / Wait / Say / Walk / Give / GoTo / …   |
| Condition library             | ✅ Flag, HasItem, SelectedItem, InScene, VarEq, Not/And/Or |
| Scene transitions (fade)      | ✅ on scene change                                       |
| Asset placeholders            | ✅ deterministic color for missing files                 |
| Stylized character placeholder | ✅ humanoid / quadruped shape when no sprite art       |
| Audio (PlayMusic/PlaySound)   | 🟡 log-only — Ebiten audio backend not wired up yet      |
| Animation (sprite sheet)      | 🟡 `AnimationClip` field exists, rendering not yet       |
| Walkbox + A\*                 | 🟡 straight-line movement only                           |
| Save / Load (JSON)            | 🟡 API in place, implementation stubbed                  |
| Trigger processing            | 🟡 struct exists, engine doesn't run them yet            |

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
