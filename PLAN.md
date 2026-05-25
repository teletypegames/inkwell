# PLAN.md — `pncdsl` point-and-click adventure library

## Cél

Egy Go nyelvű, [Ebitengine](https://ebitengine.org/) tetejére épülő könyvtár, amellyel
klasszikus point-and-click kalandjátékok (Lucasarts/Sierra stílus) gyorsan, **deklaratív
struct literal**-ekkel összerakhatók. A library a `pncdsl` modulban él, a konkrét játék
pedig a `domain` package-ben a regisztráló hívásokon keresztül definiálódik. A `main.go`
mindössze behúzza a domaint, példányosítja a `pncdsl.Game`-et és elindítja.

**Vezérlő minta — Manager-regisztráció.** Minden névvel hivatkozható entitás (item,
scene, character, dialogue, script, asset, verb) egy-egy globális (Game-en lógó)
**`XxxManager`**-be kerül a **`Register(...)`** metódussal, és a `Get(name)` adja
vissza. Például `g.ItemManager.Register(pncdsl.Item{Name: "key", ...})` és onnantól
`g.ItemManager.Get("key")`. Mindegyik manager pontosan ugyanúgy viselkedik —
egységes konvenció, nincs entitásonkénti egyedi API.

## Felhasználói élmény (DSL ízelítő)

A `domain/` minden építőkövét **külön fájlba** tesszük, a fájlnév az adott entitás
típusát és azonosítóját tükrözi (`téma.azonosító.go`). A `domain/game.go` csak
behúzza az építő függvényeket egy `Build()`-be:

```go
// domain/game.go
package domain

import p "pncdsl/pncdsl"

func Build() *p.Game {
    g := p.NewGame("Reggeli Kávé", 320, 200)

    // assetek (a domain.asset.go-ban felsorolva)
    registerAssets(g)

    // karakterek
    definePlayer(g)
    defineCat(g)

    // itemek
    defineKey(g)
    defineBeans(g)
    defineMug(g)

    // dialógusok
    defineCatMorning(g)

    // scenek
    defineBedroom(g)
    defineKitchen(g)

    // scriptek / cutscene-ek
    defineIntroScript(g)
    defineVictoryScript(g)

    g.StartAt("bedroom").OnStart(p.RunScript("intro"))
    return g
}
```

```go
// domain/scene.kitchen.go
package domain

import p "pncdsl/pncdsl"

func defineKitchen(g *p.Game) {
    g.SceneManager.Register(p.Scene{
        Name:       "kitchen",
        Background: "bg/kitchen",
        Music:      "mus/calm",
        Hotspots: []p.Hotspot{
            {
                Name:  "cupboard",
                Area:  p.Rect(40, 60, 80, 120),
                Label: "szekrény",
                OnLook: p.Say("player", "Egy konyhaszekrény. Zárva."),
                OnUse: p.If(p.Flag("cupboard_open"),
                    p.Say("player", "Már nyitva van."),
                    p.Seq(
                        p.RequireItem("key"),
                        p.Say("player", "Klikk."),
                        p.SetFlag("cupboard_open"),
                        p.Give("beans"),
                        p.TakeAway("key"),
                    ),
                ),
            },
            {
                Name:   "cat",
                Area:   p.Rect(200, 130, 30, 30),
                Label:  "Cirmos",
                OnTalk: p.RunDialogue("cat_morning"),
            },
            {
                Name:  "door",
                Area:  p.Rect(280, 60, 30, 120),
                OnUse: p.GoTo("bedroom"),
            },
        },
    })
}
```

```go
// domain/item.key.go
package domain

import p "pncdsl/pncdsl"

func defineKey(g *p.Game) {
    g.ItemManager.Register(p.Item{
        Name:        "key",
        Sprite:      "spr/key",
        Description: "rozsdás kulcs",
    })
}
```

Az építőkövek tehát **plain struct literal**-ek, amiket a megfelelő `*Manager`-be
regisztrálunk. Nincs fluent builder, nincs `.Done()` — a struct mezői önmagukban
deklaratívak, az IDE autocomplete pedig kvázi-DSL-szerű élményt ad.

```go
// main.go
package main

import (
    "pncdsl/domain"
    "pncdsl/pncdsl"
)

func main() {
    pncdsl.Run(domain.Build())
}
```

## Mappa- és fájlstruktúra

**Konvenció**: minden Go-fájl neve `téma.azonosító.go`. A pont szeparátor szándékos —
így `ls`-ben a fájlok témánként szépen csoportosulnak. A library oldalon a téma a
komponens funkcióját jelöli (`core`, `scene`, `actor`, `dialog`, `action`, `state`,
`ui`, `input`, `asset`, `util`); a domain oldalon a téma az entitás típusa (`scene`,
`item`, `character`, `dialog`, `script`), az azonosító pedig pontosan az a string,
amivel a DSL-ben hivatkozunk rá (`scene.kitchen.go` ↔ `g.DefineScene("kitchen")`).

```
pncdsl/
├── go.mod
├── PLAN.md
├── main.go                          # belépési pont
├── assets/                          # képek, hangok, fontok (futásidőben olvasva)
│   ├── bg/  spr/  snd/  mus/  fnt/
│
├── pncdsl/                          # maga a library
│   ├── core.doc.go                  # package docs, design rationale
│   ├── core.game.go                 # Game (root aggregátum), Run()
│   ├── core.engine.go               # Ebiten adapter: Update/Draw/Layout
│   ├── core.manager.go              # Manager[T Named] — KÖZPONTI minta
│   ├── core.dsl.go                  # toplevel helperek (Rect, Seq, Par, Say...)
│   ├── core.errors.go               # tipizált hibák, sentinel értékek
│   │
│   ├── scene.def.go                 # Scene struct (Name + mezők)
│   ├── scene.manager.go             # SceneManager (alias: Manager[Scene])
│   ├── scene.hotspot.go             # Hotspot struct (Scene gyermeke)
│   ├── scene.trigger.go             # Trigger (Scene gyermeke)
│   ├── scene.transition.go          # Fade, Cut, Wipe a scenek között
│   ├── scene.camera.go              # Camera (pan/follow/shake)
│   │
│   ├── item.def.go                  # Item struct (Name + mezők)
│   ├── item.manager.go              # ItemManager (alias: Manager[Item])
│   ├── item.inventory.go            # Inventory, slotok, drag&drop
│   │
│   ├── actor.def.go                 # Character struct (Name + mezők)
│   ├── actor.manager.go             # CharacterManager (alias: Manager[Character])
│   ├── actor.animation.go           # SpriteSheet, AnimationClip, Animator
│   │
│   ├── dialog.def.go                # Dialogue, DialogueNode, Choice
│   ├── dialog.manager.go            # DialogueManager (alias: Manager[Dialogue])
│   ├── dialog.box.go                # DialogBox UI a Dialogue megjelenítéshez
│   │
│   ├── action.def.go                # Action interface + beépített konstruktorok
│   ├── action.condition.go          # Condition interface + beépítettek
│   ├── action.script.go             # Script struct (Name + Actions)
│   ├── action.manager.go            # ScriptManager (alias: Manager[Script])
│   │
│   ├── state.def.go                 # World state, Flags, Vars, History
│   ├── state.save.go                # SaveSlot, serialize/restore
│   │
│   ├── ui.def.go                    # UI root, layout
│   ├── ui.verb.go                   # Verb struct + VerbBar
│   ├── ui.verb_manager.go           # VerbManager (alias: Manager[Verb])
│   ├── ui.cursor.go                 # Cursor sprite-ek és állapotok
│   ├── ui.speech.go                 # SpeechBubble renderelés
│   │
│   ├── input.def.go                 # Egér, billentyű, hover/click resolve
│   │
│   ├── asset.def.go                 # Asset struct (Name, Path, Kind)
│   ├── asset.manager.go             # AssetManager (alias: Manager[Asset]) + lazy load
│   ├── asset.audio.go               # AudioPlayer, Music, Sfx mixing
│   ├── asset.text.go                # FontLoader, TextRenderer, Wrapper
│   │
│   ├── util.geometry.go             # Point, Rect, Polygon, Path
│   ├── util.timer.go                # Tickek, Tween, Delay
│   └── util.log.go                  # opcionális debug overlay
│
└── domain/                          # a konkrét játék tartalma — lásd Demo szakasz
    ├── game.go                      # Build() — sorban hívja a define*() függvényeket
    ├── domain.asset.go              # registerAssets(g) — minden asset egy helyen
    ├── domain.flag.go               # flag/var konstansok (string típusos enum)
    │
    ├── scene.bedroom.go             # defineBedroom(g)
    ├── scene.kitchen.go             # defineKitchen(g)
    │
    ├── item.key.go                  # defineKey(g)
    ├── item.beans.go                # defineBeans(g)
    ├── item.mug.go                  # defineMug(g)
    │
    ├── character.player.go          # definePlayer(g)
    ├── character.cat.go             # defineCat(g)
    │
    ├── dialog.cat_morning.go        # defineCatMorning(g)
    │
    ├── script.intro.go              # defineIntroScript(g)
    └── script.victory.go            # defineVictoryScript(g)
```

A domain-konvenció előnye: amikor új scene-t/itemet teszek hozzá, **csak egy új fájlt
hozok létre + egy sort a `Build()`-ben**. Nincs egyetlen ezer soros `items.go`. Ha a
jövőben többszáz entitás lesz, a fájllista önmagában dokumentáció: `ls domain/scene.*`
megmutatja az összes helyszínt.

## Központi struktúrák

### `Manager[T]` — a regisztráló minta (`pncdsl/core.manager.go`)

Minden névvel hivatkozható entitást egyetlen generikus `Manager` kezel. A library
egyetlen helyen definiálja, és típus-aliasokkal kap nevet entitásonként — így minden
manager egységesen viselkedik (`Register`, `Get`, `MustGet`, `Has`, `Names`, `Each`).

```go
// pncdsl/core.manager.go
package pncdsl

type Named interface {
    GetName() string
}

type Manager[T Named] struct {
    items map[string]T
}

func NewManager[T Named]() *Manager[T] {
    return &Manager[T]{items: make(map[string]T)}
}

func (m *Manager[T]) Register(v T) {
    name := v.GetName()
    if name == "" {
        panic("pncdsl: Register: empty Name")
    }
    if _, dup := m.items[name]; dup {
        panic("pncdsl: Register: duplicate name: " + name)
    }
    m.items[name] = v
}

func (m *Manager[T]) Get(name string) (T, bool) {
    v, ok := m.items[name]
    return v, ok
}

func (m *Manager[T]) MustGet(name string) T {
    v, ok := m.items[name]
    if !ok {
        panic("pncdsl: MustGet: unknown name: " + name)
    }
    return v
}

func (m *Manager[T]) Has(name string) bool        { _, ok := m.items[name]; return ok }
func (m *Manager[T]) Names() []string             { /* sorted keys */ }
func (m *Manager[T]) Each(fn func(T))             { /* iterate */ }
```

Minden entitás implementálja a `Named` interfészt egyetlen kis metódussal:

```go
type Item struct {
    Name        string
    Sprite      string
    Description string
    // ... handlerek
}
func (i Item) GetName() string { return i.Name }
```

Entitásonkénti **nevesített manager-típus** (generic alias, Go 1.24+):

```go
// pncdsl/item.manager.go
type ItemManager = Manager[Item]

// pncdsl/scene.manager.go
type SceneManager = Manager[Scene]

// pncdsl/actor.manager.go
type CharacterManager = Manager[Character]

// pncdsl/dialog.manager.go
type DialogueManager = Manager[Dialogue]

// pncdsl/action.manager.go
type ScriptManager = Manager[Script]

// pncdsl/asset.manager.go
type AssetManager = Manager[Asset]

// pncdsl/ui.verb_manager.go
type VerbManager = Manager[Verb]
```

A használat ezért minden entitáson **pontosan ugyanígy néz ki**:

```go
g.ItemManager.Register(p.Item{Name: "key", Sprite: "spr/key"})
g.ItemManager.Get("key")           // (Item, bool)
g.ItemManager.MustGet("key")       // Item, panic ha nincs
g.SceneManager.Register(p.Scene{Name: "kitchen", ...})
g.AssetManager.Register(p.Asset{Name: "bg/kitchen", Path: "assets/bg/kitchen.png"})
```

**Manager-listája:** `ItemManager`, `SceneManager`, `CharacterManager`,
`DialogueManager`, `ScriptManager`, `AssetManager`, `VerbManager`. (A `Hotspot`,
`Trigger`, `Walkbox`, `DialogueNode` nem kap saját managert — ezek mindig egy
parent struct mezőiként élnek, scoped a Scene-hez ill. a Dialogue-hoz.)

### `Game` — a gyökér aggregátum (`pncdsl/core.game.go`)

A `Game` az összes managert mezőként hordozza, plusz a futásidejű állapotot:

```go
type Game struct {
    Title         string
    Width, Height int

    // regisztrálható entitások — minden mező egyforma alakú
    ItemManager      *ItemManager
    SceneManager     *SceneManager
    CharacterManager *CharacterManager
    DialogueManager  *DialogueManager
    ScriptManager    *ScriptManager
    AssetManager     *AssetManager
    VerbManager      *VerbManager

    // futásidejű állapot
    State     *State
    Inventory *Inventory
    UI        *UI
    Audio     *AudioPlayer
    Camera    *Camera
    Cursor    *Cursor

    current *Scene
    startID string
    onStart Action
}

func NewGame(title string, w, h int) *Game {
    return &Game{
        Title:            title,
        Width:            w,
        Height:           h,
        ItemManager:      NewManager[Item](),
        SceneManager:     NewManager[Scene](),
        CharacterManager: NewManager[Character](),
        DialogueManager:  NewManager[Dialogue](),
        ScriptManager:    NewManager[Script](),
        AssetManager:     NewManager[Asset](),
        VerbManager:      NewManager[Verb](),
        State:            NewState(),
        Inventory:        NewInventory(),
        // ...
    }
}
```

A `Game` többi metódusa minimális — nem builder, csak konfig:

- `func NewGame(title string, w, h int) *Game`
- `func (g *Game) StartAt(sceneName string) *Game`
- `func (g *Game) OnStart(a Action) *Game`
- `func (g *Game) Validate() error`
- `func (g *Game) Run() error`
- `func (g *Game) Save(slot int) error` / `func (g *Game) Load(slot int) error`

A toplevel `pncdsl.Run(g *Game) error` annyi, hogy `g.Validate()` után belép az
Ebiten ciklusba.

### `Asset` (`pncdsl/asset.def.go`)

```go
type Asset struct {
    Name string       // pl. "bg/kitchen"
    Path string       // pl. "assets/bg/kitchen.png"
    Kind AssetKind    // Image | Audio | Font
}
func (a Asset) GetName() string { return a.Name }
```

Az `AssetManager` lazy: `Register` csak metaadatot tárol, az első tényleges
hozzáférésnél tölti be és cache-eli (a betöltött resource az `AssetManager`
belső térképében, nem az `Asset` structban — utóbbi marad immutable konfig).

### `Scene` (`pncdsl/scene.def.go`)

```go
type Scene struct {
    Name       string
    Background string          // Asset.Name
    Music      string          // Asset.Name
    Hotspots   []Hotspot
    Walkboxes  []Polygon
    Triggers   []Trigger
    Actors     []SceneActor    // { CharacterName string; At Point }
    OnEnter    Action
    OnLeave    Action
}
func (s Scene) GetName() string { return s.Name }
```

Hotspot / Trigger / Walkbox értékek inline, ugyanabban a struct literalban —
nem kell külön regisztrálni őket; egy Hotspot eleve scene-scoped.

### `Hotspot` (`pncdsl/scene.hotspot.go`)

```go
type Hotspot struct {
    Name   string
    Area   Shape          // Rect | Polygon
    Label  string         // a kurzor melletti felirat hover-kor
    Cursor CursorKind

    // verbenkénti handler — egyszerű mezők, nem map
    OnLook Action
    OnUse  Action
    OnTalk Action
    OnTake Action
    OnGive Action

    // egyéni verbek (`g.VerbManager.Register(Verb{Name: "push"})`-hoz)
    OnVerb map[string]Action
}
```

### `Item` (`pncdsl/item.def.go`)

```go
type Item struct {
    Name        string
    Sprite      string                 // Asset.Name
    Description string

    OnUseSelf Action                   // pl. „nem tudom mit kezdjek vele"
    OnUseWith map[string]Action        // targetName -> action (hotspot vagy másik item)
}
func (i Item) GetName() string { return i.Name }
```

### `Inventory` (`pncdsl/item.inventory.go`)

```go
type Inventory struct {
    items    []string
    capacity int
    selected string
}

func (i *Inventory) Add(name string)
func (i *Inventory) Remove(name string)
func (i *Inventory) Has(name string) bool
func (i *Inventory) Select(name string)
func (i *Inventory) Selected() string
```

(Az `Inventory` nem manager — futásidejű állapot, nem regisztráció.)

### `Character` (`pncdsl/actor.def.go`)

```go
type Character struct {
    Name        string
    Sprite      string                         // Asset.Name
    Animations  map[string]AnimationClip       // idle, walk_left, walk_right, talk
    Speed       float64
    SpeechColor color.Color
    Start       Point                          // alappozíció (scene-en belül felülírhatja Actor)
}
func (c Character) GetName() string { return c.Name }
```

Futásidejű mozgás-API a `Game`-en keresztül érhető el (nem a Character struct
módosításán át, hiszen az regisztrált konfig):
`g.WalkCharacter("player", Point{...})`, `g.FaceCharacter("player", dir)`.

### `Dialogue` (`pncdsl/dialog.def.go`)

```go
type Dialogue struct {
    Name  string
    Start string                       // kezdő node Name
    Nodes []DialogueNode               // mezőként, nem manager — Dialogue-scoped
}
func (d Dialogue) GetName() string { return d.Name }

type DialogueNode struct {
    Name    string
    Lines   []DialogueLine             // ki mit mond a node belépésén
    Choices []DialogueChoice
}

type DialogueLine struct {
    Speaker string                     // Character.Name
    Text    string
}

type DialogueChoice struct {
    Text    string
    Show    Condition                  // nil = mindig látszik
    Once    bool
    Actions []Action
}
```

### `Script` (`pncdsl/action.script.go`)

```go
type Script struct {
    Name    string
    Actions Action                     // tipikusan Seq(...)
}
func (s Script) GetName() string { return s.Name }
```

### `Verb` (`pncdsl/ui.verb.go`)

```go
type Verb struct {
    Name    string                     // "look", "use", "talk", "take", ...
    Label   string                     // megjelenő szöveg a VerbBar-on
    Default Action                     // ha a hotspot nem definiál erre handlert
}
func (v Verb) GetName() string { return v.Name }
```

### `Action` interface (`pncdsl/action.def.go`)

```go
type Action interface {
    Run(ctx *Ctx) Status              // Running | Done | Failed
}

type Ctx struct {
    Game    *Game
    Scene   *Scene
    Hotspot *Hotspot
    Item    *Item
    dt      float64
}
```

Beépített actionök (sima konstruktorfüggvények, hogy DSL-ből szépen olvasódjon):

- `Say(speakerID, text string)`
- `Wait(seconds float64)`
- `GoTo(sceneID string)`
- `Walk(charID string, to Point)`
- `Give(itemID string)` / `TakeAway(itemID string)`
- `SetFlag(name string)` / `ClearFlag(name string)` / `SetVar(name string, v any)`
- `PlaySound(assetID string)` / `PlayMusic(assetID string)` / `StopMusic()`
- `RunDialogue(dialogueID string)` / `EndDialogue()` / `GotoNode(nodeID string)`
- `RunScript(scriptID string)`
- `Seq(a ...Action)` — sorban (default builder-szemantika)
- `Par(a ...Action)` — párhuzamosan
- `If(c Condition, then Action, else_ ...Action)`
- `RequireItem(itemID string)` — short-circuit ha nincs nála (visszadob egy "nem tudom használni"-t)
- `Custom(fn func(*Ctx) Status)` — escape hatch a domainnek

### `Condition` (`pncdsl/action.condition.go`)

```go
type Condition interface{ Eval(*Ctx) bool }
```

Beépítettek: `HasItem(name)`, `Flag(name)`, `Not(c)`, `And(c...)`, `Or(c...)`,
`VarEq(name, v)`, `InScene(name)`, `SelectedItem(name)`.

### `State` (`pncdsl/state.def.go`)

```go
type State struct {
    flags   map[string]bool
    vars    map[string]any
    visited map[string]int   // scene name -> hányszor lépett be
    talked  map[string]int   // dialogue node -> hányszor futott le (a `Once` choice-okhoz)
}
```

Metódusok: `Flag/SetFlag/ClearFlag`, `Var/SetVar`, `Visited(name)`, `TalkedTo(name)`,
`Snapshot()`, `Restore(snap)`. (Ez sem manager — runtime állapot.)

### `UI` (`pncdsl/ui.def.go`, `pncdsl/dialog.box.go`)

```go
type UI struct {
    verbBar      *VerbBar
    invBar       *InventoryBar
    speech       *SpeechBubble
    dialogBox    *DialogBox
    hoverLabel   string
}
```

Konfigurálható: layout, font, színek. Default a SCUMM-szerű alsó verbsáv + jobboldali
inventory; ki is kapcsolható ha valaki modern verb-coin variánst akar.

### `Engine` — Ebiten adapter (`pncdsl/core.engine.go`)

```go
type engine struct {
    g *Game
}

func (e *engine) Update() error
func (e *engine) Draw(screen *ebiten.Image)
func (e *engine) Layout(outW, outH int) (int, int)
```

Az `Update` belerajzol:
1. input → resolve (melyik hotspoton van a kurzor, milyen verb aktív)
2. aktív cutscene/dialogue lépése
3. karakterek mozgása, animation tick
4. trigger feltételek
5. UI állapot frissítése

A `Draw` rétegei alulról felfele: háttér → walkbox debug (csak ha bekapcsolva) →
karakterek (y-sorted) → speech bubble-k → UI → cursor → tranzíció overlay.

## Asset és audio kezelés

Az `AssetManager` egy hagyományos `Manager[Asset]` — a `Register` csak az
útvonalat és típust tárolja, az első `Load(name)` hívásra olvas be a diszkről
és cache-el (a betöltött `*ebiten.Image` / audio buffer egy belső térképben él,
az `Asset` struct maga immutable konfig marad).

```go
g.AssetManager.Register(p.Asset{Name: "bg/kitchen", Path: "assets/bg/kitchen.png", Kind: p.AssetImage})
g.AssetManager.Register(p.Asset{Name: "mus/calm",   Path: "assets/mus/calm.ogg",   Kind: p.AssetAudio})
```

`AudioPlayer` egy zenecsatornát + N sfx csatornát kezel, fade-in/out támogatással
a `PlayMusic` action-höz.

## Save / Load

`SaveSlot` JSON-be sorosít:
- aktuális scene `Name`-je
- karakterek pozíciói
- `State` (flags, vars, visited, talked)
- inventory tartalom (item `Name`-ek listája)
- futó script PC (ha cutscene közben menteni szabad — opcionálisan tilthatóan)

A managerek **tartalmát nem mentjük** — azok a `Build()`-ből minden indításkor
újra előállnak. Csak a runtime-állapotot persistáljuk.

API: `Game.Save(slot int) error`, `Game.Load(slot int) error`.

## Hibakezelés és validáció

`Game.Validate() error` a `Run()` előtt fut. Ellenőrzi a managerek közötti
név-hivatkozások konzisztenciáját:

- minden `Scene.Background`/`Music` és `Character.Sprite` szerepel-e az `AssetManager`-ben,
- minden `GoTo(name)` célja él-e a `SceneManager`-ben,
- `RunDialogue(name)` és `Dialogue` `Choice.Actions`-ben szereplő `GotoNode` ill.
  `RunScript(name)` célok léteznek-e,
- minden `Hotspot.OnUseWith` map kulcs (item/hotspot név) feloldható,
- `StartAt(name)` scene létezik,
- nincs két azonos nevű regisztráció (ezt amúgy a `Register` panic-cal megfogja).

Hibák tipizáltak: `ErrUnknownAsset`, `ErrUnknownScene`, `ErrUnknownDialogue`,
`ErrUnknownDialogueNode`, `ErrUnknownScript`, `ErrUnknownItem`, `ErrDuplicateName`.

## Tesztelési stratégia

- A `pncdsl/*.go` minden nem-render részére táblavezérelt unit teszt
  (`action_test.go`, `dialogue_test.go`, `state_test.go`).
- A `domain/`-hoz egy headless smoke test, ami `Build()` után `Validate()`-et hív —
  így a játék statikus integritását CI-ben is le lehet csekkolni Ebiten ablak nélkül.
- A renderhez `ebiten/ebitenutil` segítségével offscreen image diff — opcionális,
  csak ha sok regressziónk lenne.

## Demo játék: „Reggeli Kávé"

A library kipróbálására egy szándékosan apró, ~5 perces játék, ami a fontos
mechanikákat **egyenként** gyakorolja anélkül, hogy a domain túlnőne a libraryn.

### Sztori egy mondatban

Frissen ébredsz, és muszáj kávét főznöd, mielőtt a macskád szétszedné a konyhát.

### Helyszínek

| ID         | Háttér              | Zene        | Tartalom                                                   |
|------------|---------------------|-------------|------------------------------------------------------------|
| `bedroom`  | `bg/bedroom.png`    | `mus/wakeup`| ágy, éjjeliszekrény, ajtó                                  |
| `kitchen`  | `bg/kitchen.png`    | `mus/calm`  | szekrény, polc, kávéfőző, macska, ajtó                     |

### Itemek

| ID      | Sprite          | Honnan szerezhető                       | Mire jó                                  |
|---------|-----------------|------------------------------------------|-------------------------------------------|
| `key`   | `spr/key.png`   | `bedroom.nightstand` → Take             | `kitchen.cupboard` kinyitása              |
| `beans` | `spr/beans.png` | `kitchen.cupboard` (kinyitás után)      | `kitchen.coffee_machine`-be tölteni       |
| `mug`   | `spr/mug.png`   | `kitchen.shelf` → Take                  | `kitchen.coffee_machine`-re tenni         |

### Karakterek

- **`player`** — az általunk irányított PC. `WalkTo` célokra mozog, beszél (`Say`).
- **`cat`** — az NPC, a konyhában ül. Csak `Talk` interakció, dialog-fa.

### Dialógus (`dialog.cat_morning.go`)

```
[start]
cat: Mióóóóóóóóu. Késő van.
  > „Mindjárt megy a kávé." → SetFlag("promised_coffee") → end
  > „Te etted meg a babot?" → goto cat_beans
  > „Hagyj békén." → end

[cat_beans]   (csak ha NEM flag `cupboard_open`)
cat: A szekrényben van. Mindig ott szokott lenni.
  > „Köszi." → end
```

### Hotspot-tábla

**`scene.bedroom.go`**

| Hotspot      | Look                              | Use                                                              |
|--------------|-----------------------------------|------------------------------------------------------------------|
| `bed`        | „A párnám még meleg."             | `Say("player", "Most keltem fel.")`                              |
| `nightstand` | „Az éjjeliszekrényemen egy kulcs." | egyszer: `Give("key")` + `SetFlag("key_taken")`; utána: „Üres." |
| `door`       | „Az ajtó a konyhába vezet."       | `GoTo("kitchen")`                                                |

**`scene.kitchen.go`**

| Hotspot          | Look                          | Use                                                                                                                                                  |
|------------------|-------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| `cupboard`       | „Konyhaszekrény. Zárva."     | ha `cupboard_open`: „Már nyitva."; különben `RequireItem("key")` → `SetFlag("cupboard_open")` + `Give("beans")` + `TakeAway("key")`                  |
| `shelf`          | „A polc tele bögrékkel."     | egyszer: `Give("mug")`; utána: „Tiszta már mind."                                                                                                    |
| `coffee_machine` | „A jó öreg kávéfőző."        | progresszív: 0/2 → kell `beans` és `mug`; 1/2 → kell a másik; 2/2 → `RunScript("victory")`                                                            |
| `cat` (NPC)      | „Cirmos, a házikedvencem."  | (Talk) `RunDialogue("cat_morning")`                                                                                                                  |
| `door`           | „Vissza a hálóba."           | `GoTo("bedroom")`                                                                                                                                    |

A „progresszív" használat egy kis state-trükkel: a kávéfőző két flaget figyel
(`coffee_has_beans`, `coffee_has_mug`); a `Use` ág `If` action-ökkel szétágazik
aszerint, hogy melyik item van kiválasztva az inventory-ban (`SelectedItem`
condition), és csak a megfelelőt fogadja el. Konkrétan így néz ki a hotspot
regisztrációja:

```go
// domain/scene.kitchen.go (részlet)
{
    Name: "coffee_machine",
    Area: p.Rect(140, 70, 50, 60),
    Label: "kávéfőző",
    OnLook: p.Say("player", "A jó öreg kávéfőző."),
    OnUse: p.Seq(
        p.If(p.And(p.SelectedItem("beans"), p.Not(p.Flag("coffee_has_beans"))),
            p.Seq(p.SetFlag("coffee_has_beans"), p.TakeAway("beans"),
                  p.Say("player", "Bab betöltve."))),
        p.If(p.And(p.SelectedItem("mug"), p.Not(p.Flag("coffee_has_mug"))),
            p.Seq(p.SetFlag("coffee_has_mug"), p.TakeAway("mug"),
                  p.Say("player", "Bögre a helyén."))),
        p.If(p.And(p.Flag("coffee_has_beans"), p.Flag("coffee_has_mug")),
            p.RunScript("victory")),
    ),
}
```

És így néz ki egy item regisztrációja, ami egy másik item-re vagy hotspotra reagál:

```go
// domain/item.key.go
func defineKey(g *p.Game) {
    g.ItemManager.Register(p.Item{
        Name:        "key",
        Sprite:      "spr/key",
        Description: "rozsdás kulcs",
        OnUseWith: map[string]p.Action{
            "cupboard": p.Seq(p.SetFlag("cupboard_open"), p.Give("beans"),
                              p.TakeAway("key"), p.Say("player", "Klikk.")),
        },
    })
}
```

### Scriptek

**`script.intro.go`** — az `OnStart`-ra kötve, `bedroom` scene-ben:

```
Wait(0.5) → fade in →
Say("player", "Brr, hideg van.") →
Say("player", "Egy kávé kéne...") →
Walk("player", Point{160, 140})
```

**`script.victory.go`** — amikor a `coffee_machine` mindkét flaget látja:

```
Say("player", "Készen is van!") →
PlaySound("snd/coffee_done") →
Wait(0.8) →
Say("cat", "Mióóóóu. *boldog macska*") →
Say("player", "Reggeli kávé: megmentve.") →
fade out → ShowEnd("Vége — kösz hogy játszottál!")
```

### Flag/var táblázat (`domain.flag.go`)

```go
const (
    FlagKeyTaken       = "key_taken"
    FlagCupboardOpen   = "cupboard_open"
    FlagCoffeeHasBeans = "coffee_has_beans"
    FlagCoffeeHasMug   = "coffee_has_mug"
    FlagPromisedCoffee = "promised_coffee"
)
```

### Mit gyakorol a demo

| Library képesség        | Hol jelenik meg a demoban                                  |
|--------------------------|------------------------------------------------------------|
| `Scene` + `GoTo`         | bedroom ↔ kitchen ajtó                                     |
| `Hotspot` + verbek       | minden helyszínen több hotspot, Look/Use/Talk              |
| `Inventory` + `Give`     | key, beans, mug felvétele                                  |
| `RequireItem`            | szekrény nyitása csak kulccsal                             |
| Use-with-item            | kávéfőző két különböző itemet kér, progresszíven           |
| `Flag` + `If`            | szekrény „már nyitva", éjjeliszekrény „már üres"           |
| `Dialogue` + `Choice`    | macska dialógusfája feltételes ággal                      |
| `Script` / cutscene      | intro (mozgás+beszéd) és victory (több actor)              |
| `Walk` + walkbox         | player A* a kattintott pontig                              |
| `PlayMusic` / `PlaySound`| szobánkénti zene, victory hangeffekt                       |
| `Save` / `Load`          | mid-game mentés, pl. a szekrény kinyitása után             |
| `Validate()`             | CI smoke test a `Build()` outputra                         |

Ami **nincs benne** szándékosan: több aktor egyszerre mozogva, kamera-pan,
inventory-tárgyak kombinálása (item+item). Ezeket egy esetleges 2. demo (pl.
"A pajta") gyakorolná, de a libraryban már első körben benne van a kapacitás.

## Implementációs ütemterv

1. **Csontváz** — `Game`, `Scene`, `Hotspot`, `AssetRegistry`, Ebiten `engine`, egér +
   háttér + egy darab Look kattintás. Egy szoba, három hotspot, működő `Say`.
2. **Verbek és inventory** — `VerbBar`, `Inventory`, `Item`, `UseOn`. SCUMM-szerű flow.
3. **Karakterek és walkboxok** — `Character`, `Walkbox`, A* az útkereséshez.
4. **Dialógus rendszer** — `Dialogue`, `DialogBox`, `DialogueChoice`, `Once`.
5. **Cutscene / Script** — `Action` interface, `Seq`/`Par`/`If`, `RunScript`.
6. **State, Flags, Save/Load** — perzisztencia.
7. **Polish** — tranzíciók, audio fade, font wrapping, debug overlay, validáció.

Minden lépés végén egy minimális `domain/` példajáték validálja, hogy a DSL valóban
kényelmes-e arra, amire ki van találva.
