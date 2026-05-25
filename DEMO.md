# DEMO.md — "Morning Coffee"

The `domain/` folder contains a 5-minute mini-game that exercises every
important `pncdsl` feature in one sitting — without the demo itself
outgrowing the library. The code is a good place to learn from, because
each mechanic appears in **exactly one small file**.

## Story

You've just woken up in your bedroom. The cat (Whiskers) is already
meowing in the kitchen, and you'd better make coffee before she tears
the whole apartment apart.

## Running it

```bash
go run .
```

A window opens at 1280×800 (internal resolution 320×200, 4× scale). No
asset files are needed — every graphic is a colored placeholder /
stylized humanoid until you replace it with real art (see
[`GFX.md`](GFX.md)).

Controls:

- **left click** on a hotspot/UI button: run the active verb
- **left click** on an inventory item: select it (the cursor picks it up)
- **left click** on a hotspot while an item is selected: use-with
- **right click**: deselect the held item, or reset the verb to "Look"
- click during dialog: advance

## Walkthrough (spoilers)

1. **Bedroom** — the `intro` script runs two `Say` lines from the
   `player`.
2. Click the **nightstand** with the "Take" verb → you receive `key`.
3. Click the **door** → you switch to the kitchen.
4. **Kitchen** — Whiskers is sitting in the corner. Try "Talk" on her —
   a conditional dialog opens; one of the branches is only available
   while the cupboard hasn't been opened yet.
5. Click the **shelf** with "Take" → you receive `mug`.
6. Click the **cupboard** with "Use" (key in inventory) → it opens, you
   automatically receive `beans`, and the key is consumed. *Or* select
   the key in the inventory and click the cupboard — `OnUseWith` does
   the same thing.
7. Click the **coffee machine** with `beans` → beans loaded.
8. Click the coffee machine with `mug` → mug in place.
9. Once both flags are set, the coffee machine fires the `victory`
   script: two characters talk, a sound plays (stubbed), end card → done.

## Scenes

| ID         | Background       | Music        | Contents                                                  |
|------------|------------------|--------------|-----------------------------------------------------------|
| `bedroom`  | `bg/bedroom`     | `mus/wakeup` | bed, nightstand (key), door                               |
| `kitchen`  | `bg/kitchen`     | `mus/calm`   | cupboard, shelf, coffee machine, cat (NPC), door          |

## Items

| ID      | Acquired from                                | Used for                                                            |
|---------|----------------------------------------------|---------------------------------------------------------------------|
| `key`   | `bedroom.nightstand` → Take                  | Use on `kitchen.cupboard`, or as a selected item via `OnUseWith`    |
| `beans` | `kitchen.cupboard` (after unlocking)         | Selected, used on `kitchen.coffee_machine`                          |
| `mug`   | `kitchen.shelf` → Take                       | Selected, used on `kitchen.coffee_machine`                          |

## Characters

- **`player`** — the player character. Walks (`Walk`), talks (`Say`).
  Size 28×62 (W×H), roughly the SCUMM-style 30% screen height.
- **`cat`** — Whiskers, the household pet. Present in the kitchen, only
  responds to `Talk`, has a dialog tree.

Until you provide real sprites, they're rendered as a stylized humanoid
(skin-tone head + colored torso + dark trousers + feet) and quadruped
(side-view with ears, head, eye, tail) placeholder shape.

## Dialog (`dialog.cat_morning.go`)

```
[start]
cat: Meeeoooow. It's late.
  > "Coffee's coming right up."         → SetFlag(promised_coffee), end
  > "Did you eat the beans?"             (only while NOT cupboard_open) → goto beans
  > "Leave me alone."                    → end

[beans]
cat: They're in the cupboard. They're always in the cupboard.
  > "Thanks."                            → end
```

The `Show: p.Not(p.Flag(FlagCupboardOpen))` condition makes the "Did
you eat the beans?" choice disappear once you've opened the cupboard
— this demonstrates conditional choice visibility
(`DialogueChoice.Show`).

## Cutscenes

**`script.intro.go`** — wired up to `OnStart`, runs in the bedroom:

```go
p.Seq(
    p.Wait(0.4),
    p.Say("player", "Brr, it's cold."),
    p.Say("player", "I need coffee before my head explodes."),
    p.Walk("player", p.Point{X: 200, Y: 145}),
)
```

**`script.victory.go`** — fired when the coffee machine sees both flags:

```go
p.Seq(
    p.Say("player", "All done!"),
    p.PlaySound("snd/coffee_done"),
    p.Wait(0.6),
    p.Say("cat", "Mrrrrow. *content cat noises*"),
    p.Say("player", "Morning coffee: saved."),
    p.ShowEnd("The End — thanks for playing!"),
)
```

## Flag/var table (`domain.flag.go`)

```go
const (
    FlagKeyTaken       = "key_taken"
    FlagCupboardOpen   = "cupboard_open"
    FlagMugTaken       = "mug_taken"
    FlagCoffeeHasBeans = "coffee_has_beans"
    FlagCoffeeHasMug   = "coffee_has_mug"
    FlagPromisedCoffee = "promised_coffee"
)
```

Naming the flags as constants is worth it because the type system then
catches typos — `pncdsl.Flag("cupbord_open")` would otherwise silently
return `false` forever.

## What the demo exercises

| Library feature                    | Where it shows up                                              |
|------------------------------------|----------------------------------------------------------------|
| `Scene` + `GoTo`                   | bedroom ↔ kitchen door, with fade transition                   |
| `Hotspot` with every verb          | each scene has multiple hotspots: Look / Use / Talk / Take     |
| `Inventory` + `Give` / `TakeAway`  | picking up and consuming key, beans, mug                       |
| `RequireItem` semantics            | cupboard checks `HasItem("key")`                               |
| `OnUseWith` (from the hotspot)     | cupboard + key                                                 |
| `OnUseWith` (from the item)        | coffee machine + beans / mug                                   |
| `Flag` + `If`                      | cupboard "already open", shelf "already took one"              |
| `Dialogue` + `Choice.Show`         | cat's dialog tree, vanishing branch driven by a flag           |
| `Script` / cutscene                | intro (walk+say) + victory (multi-actor, sound, end card)      |
| `Walk` with straight-line stepping | every interaction is preceded by a walk target                 |
| `ShowEnd`                          | last action of the victory script                              |
| `Validate()`                       | `domain/build_test.go` CI-friendly smoke test                  |

## Files in `domain/`

```
domain/
├── game.go                  # Build() — calls every define*() in order
├── domain.asset.go          # registerAssets(g) — all assets in one place
├── domain.flag.go           # flag-name constants
│
├── character.player.go      # definePlayer(g)
├── character.cat.go         # defineCat(g)
│
├── item.key.go              # defineKey(g) — OnUseWith: {"cupboard": ...}
├── item.beans.go            # defineBeans(g)
├── item.mug.go              # defineMug(g)
│
├── dialog.cat_morning.go    # defineCatMorning(g) — 2 nodes, 3+1 choices
│
├── scene.bedroom.go         # defineBedroom(g)
├── scene.kitchen.go         # defineKitchen(g)
│
├── script.intro.go          # defineIntroScript(g)
├── script.victory.go        # defineVictoryScript(g)
│
└── build_test.go            # Build() + Validate() smoke test
```

Adding a new scene/item: one new file + one line in `Build()` inside
`game.go`. There's no central registry list to edit.

## How to extend it

- **New scene**: create `scene.<name>.go` with a `define<Name>(g *p.Game)`
  function that calls `g.SceneManager.Register(...)`. Add a call in
  `Build()`. If any hotspot uses `GoTo("<name>")`, `Validate()` will
  cross-check the link for you.
- **New item**: same pattern, `item.<name>.go`, register into
  `g.ItemManager`. Map keys inside `OnUseWith` must exactly match the
  target hotspot's `Name` field.
- **New dialog**: `dialog.<name>.go` — `DialogueChoice.Show` accepts
  anything that implements `Condition` (`Flag`, `HasItem`, `Not(...)`,
  etc.).
- **New cutscene**: `script.<name>.go`, register into
  `g.ScriptManager`, then fire it anywhere with
  `p.RunScript("<name>")`.

## Graphics

The library substitutes a **deterministic colored rectangle** for every
missing image and renders characters as a stylized humanoid/quadruped
shape. The game therefore runs without any asset on disk. When you drop
a generated image into `assets/bg/<name>.png` or
`assets/spr/<name>.png`, the library switches to it immediately — no
code change required.

Prompts and positioning tips for an image AI: [`GFX.md`](GFX.md).

## Back to the library

For the framework structure and the Manager pattern itself, see
[`README.md`](README.md) and [`PLAN.md`](PLAN.md).
