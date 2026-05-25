# UIPLAN.md — dinamikus widget-alapú UI réteg

## Cél

A `pncdsl` jelenlegi kezelőfelülete hard-kódolt komponensek halmaza: van egy
fix verb-bar, fix inventory-bar, fix dialog-box, fix speech-bubble. Cserélni,
átalakítani őket csak az `ui.*.go` fájlok átírásával lehet.

Át kell térni egy **dinamikus widget-rendszerre**, ami:

- minden HUD-elemet egységesen, **`Widget` interfészként** kezel,
- engedi a `domain`-nak hogy **saját widget-eket adjon hozzá** (pl. egy chat
  panel) anélkül, hogy a library kódjához nyúlna,
- támogatja az alternatív input-módokat — pl. **verb-coin** (radiális menü
  másodlagos kattintásra) a SCUMM-szerű verb-bar helyett,
- a jelenlegi „look and feel"-t **preset Theme-ekként** csomagolja, amik
  egyetlen sorral aktiválhatók,
- mindezt a meglévő **`Manager.Register(...)` mintával** integrálja, hogy ne
  legyen idegen test a libraryban.

## Az alapfogalom: `Widget`

Egy widget bármi, ami:

1. névvel hivatkozható (lookup, ki/be-kapcsolás),
2. frame-enként frissít magán állapotot és input-ot (Tick),
3. frame-enként rajzol magát a képernyőre (Draw).

```go
// pncdsl/ui.widget.go
type Widget interface {
    Named                            // GetName() string
    Tick(ctx *UICtx)                 // állapot + input feldolgozás
    Draw(dst *ebiten.Image, ctx *UICtx)
}

type UICtx struct {
    Game *Game
    DT   float64
}
```

A widget **maga felel a saját geometriájáért** — nincs külső layout-rendszer,
ami pozícionálja. Egy widget lehet:

- **statikus** (mindig ugyanazon a fix téglalapon van, pl. `VerbBar`),
- **dinamikus** (a bounds-ot futás közben számolja, pl. a `RadialVerbs`
  a click-pozícióból, vagy a `SpeechBubble` a beszélő karakter feje fölött),
- **modális** (csak bizonyos állapotban látszik / fogad inputot,
  pl. `DialogBox` amikor egy `Dialogue` aktív).

A bounds tehát **method, nem mező** — ha kellene, a widget egy `Bounds()`
implementál belsőleg, de a library nem követeli meg. Ez a dynamic layout DSL
lényege: nincs nagy `UILayout` struct fix mezőkkel; a layout **a widget-fa
maga**.

## `UIManager`: a megszokott regiszter-minta

A többi entitáshoz hasonlóan a widgetek is egy generic Manager-ben élnek:

```go
// pncdsl/ui.manager.go
type UIManager = Manager[Widget]
```

```go
g.UIManager.Register(&pncdsl.VerbBar{Name: "verbs"})
g.UIManager.Register(&pncdsl.InventoryBar{Name: "inventory"})
g.UIManager.Register(&pncdsl.RadialVerbs{Name: "verb_coin", Trigger: pncdsl.MouseButtonRight})
g.UIManager.Register(&domain.ChatPanel{Name: "chat", Bounds: pncdsl.Rect(220, 4, 96, 130)})
```

Hideg újdonság a `Manager`-ben: a `Widget` egy **interfész**, nem konkrét
struct. A `Manager[T Named]` ezt eddig is kezelte (T = `Widget`); csak most
használjuk ki először.

**Z-order**: a regisztráció sorrendje. Korábban regisztrált widget hátrább van.

- `Draw` **regisztráció sorrendben** fut → későbbi widget rajzol felülre.
- `Tick` **fordított sorrendben** fut → felül lévő widget kapja először az
  inputot. Ha a top-most widget `ctx.Game.input.ConsumeLeft()`-ot hív, az
  alatta lévők `input.LeftClicked()` hívása már `false`-ot ad — így a
  click „elnyelődik". Ez ugyanaz a consume-on-use pattern, mint amit a
  `Say` action és a hotspot-klikk-resolve használ.

## Beépített widget-ek

A library a következő `Widget` implementációkat szállítja (mindegyik egy-egy
`ui.<név>.go` fájlban):

| Widget          | Mit csinál                                                       |
|-----------------|------------------------------------------------------------------|
| `Panel`         | Generic konténer — háttér-szín, opcionális keret. Csoportosításhoz. |
| `StatusLine`    | Hover-felirat + 2s-os flash üzenetek központozva.                |
| `VerbBar`       | SCUMM-szerű verb-gomb grid (mai default).                        |
| `RadialVerbs`   | Verb-coin: másodlagos kattintásra felugró radiális ige-menü.     |
| `InventoryBar`  | Item slotok, sorba/grid-be rendezve.                             |
| `SpeechBubble`  | A `Say` action által beállított beszéd buborék a beszélő fölött. |
| `DialogBox`     | Az aktív `Dialogue` választói + sorai (modális).                 |
| `EndCard`       | A `ShowEnd` által beállított befejező overlay (modális).         |
| `Cursor`        | Egér-kurzor, kiválasztott item sprite-jával ha van.              |
| `HotspotDebug`  | Hotspot-téglalapok overlay-je, fejlesztéshez (`Enabled` toggle). |

Mindegyik **konfigurálható mezőkkel** (struct literal) és **felülírható** —
ha a domain saját `VerbBar`-t akar, nem regisztrálja a beépítettet, hanem a
sajátját.

Példa egy beépített widget mezőire:

```go
type VerbBar struct {
    Name    string
    Origin  Point
    Cols    int
    Button  Size
    Gap     Point
    // belső állapot — nem kell érintened
    buttons []verbButton
}
func (v *VerbBar) GetName() string { return v.Name }
func (v *VerbBar) Tick(ctx *UICtx)                 { /* hit-test + verb-szelekció */ }
func (v *VerbBar) Draw(dst *ebiten.Image, c *UICtx) { /* gombrajzolás */ }
```

### A `RadialVerbs` widget

Az új verb-coin működése:

```go
type RadialVerbs struct {
    Name    string
    Trigger MouseButton    // alap: MouseButtonRight
    Radius  float64        // alap: 40
    // belső állapot
    visible bool
    center  Point
}
```

Lifecycle:

1. Tick-ben figyel a `Trigger`-re. Ha a játékos a scene-területen (nem UI-on)
   triggert nyom, megnyitja: `visible = true`, `center = mouse.Pos`. Konzumálja
   a kattintást, hogy ne menjen tovább a hotspot-resolve-ra.
2. Amíg `visible`: minden frame megrajzolja a `VerbManager` ige-listáját
   `len(verbs)` szelet formájában a `center` körül `Radius` sugárral.
3. Bal-klikk egy szeletre: kiválasztja az adott igét (`g.selectedVerb = ...`),
   `visible = false`. Klikk kívülre: csak bezár.
4. ESC ugyanígy bezár.

A radiális menü tehát **dinamikus pozíciójú widget** — minden frame
újraszámolt szelet-bounds-ok, a domain nem kell hogy pozícionálja.

## Custom widget — `ChatPanel` példa

A domain a `Widget` interfész implementálásával bármilyen Ebiten-eszközt be
tud kötni. Egy chat-panel kb. 40 sorból:

```go
// domain/widget.chat.go
package domain

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
    p "pncdsl/pncdsl"
)

type ChatPanel struct {
    Name     string
    Bounds   p.Rectangle
    MaxLines int

    lines []string
}

func (c *ChatPanel) GetName() string { return c.Name }

func (c *ChatPanel) Push(line string) {
    c.lines = append(c.lines, line)
    if len(c.lines) > c.MaxLines {
        c.lines = c.lines[len(c.lines)-c.MaxLines:]
    }
}

func (c *ChatPanel) Tick(ctx *p.UICtx) {
    // pl. scrollozás egér-kerékkel — szabadon
}

func (c *ChatPanel) Draw(dst *ebiten.Image, ctx *p.UICtx) {
    th := ctx.Game.Theme()    // aktív téma
    vector.DrawFilledRect(dst,
        float32(c.Bounds.X), float32(c.Bounds.Y),
        float32(c.Bounds.W), float32(c.Bounds.H),
        th.PanelBG, false,
    )
    for i, ln := range c.lines {
        ctx.Game.DrawText(dst, ln,
            int(c.Bounds.X)+2,
            int(c.Bounds.Y)+2+i*14,
            th.StatusText,
        )
    }
}
```

Bejegyzés a `domain/game.go` `Build()`-jébe:

```go
g.UIManager.Register(&ChatPanel{
    Name:     "chat",
    Bounds:   p.Rect(220, 4, 96, 130),
    MaxLines: 8,
})
```

A chat tartalmát egy új action vagy a domain-kód push-olja:

```go
g.UIManager.MustGet("chat").(*ChatPanel).Push("[server] Welcome!")
```

Vagy egy action-konstruktor: `p.Custom(func(ctx) { ... })`.

## Theme rendszer

A jelenlegi színek/színskála a `Theme` struct-ba kerül. **A téma egy
manager-be regisztrált entitás** — több téma lehet regisztrálva, és egyik
aktív:

```go
// pncdsl/ui.theme.go
type Theme struct {
    Name string

    PanelBG    color.Color
    StatusText color.Color
    FlashText  color.Color

    VerbButtonBG         color.Color
    VerbButtonSelectedBG color.Color
    VerbButtonText       color.Color

    InventorySlotBG          color.Color
    InventorySlotSelectedBG  color.Color

    SpeechBubbleBG    color.Color
    SpeechDefaultText color.Color

    DialogBG          color.Color
    DialogBorder      color.Color
    DialogChoiceBG    color.Color
    DialogChoiceHover color.Color
    DialogSpeaker     color.Color
    DialogText        color.Color

    EndCardBG   color.Color
    EndCardText color.Color

    CursorColor    color.Color
    HotspotOutline color.Color
}

func (t Theme) GetName() string   { return t.Name }
func (t Theme) TypeLabel() string { return "theme" }

type ThemeManager = Manager[Theme]
```

A `Game` mezője:

```go
type Game struct {
    // ...
    ThemeManager *ThemeManager
    activeTheme  string
}

func (g *Game) Theme() Theme            { return g.ThemeManager.MustGet(g.activeTheme) }
func (g *Game) UseTheme(name string)    { g.activeTheme = name }
```

A widgetek a `ctx.Game.Theme()` hívásból olvassák ki a színeket — **minden
frame friss**, így futás közbeni téma-csere zökkenőmentes.

### Szállított preset-ek

A library a következő témákat regisztrálja induláskor (vagy a domain hívja
őket `pncdsl.RegisterPresetThemes(g)`-vel):

| Preset                | Mit reprezentál                                                       |
|-----------------------|------------------------------------------------------------------------|
| `classic-scumm`       | A jelenlegi sötét-kékes paletta, SCUMM-szerű, alapértelmezett.          |
| `sierra-coin`         | Sötétebb, a verb-coin használathoz hangolt (nincs alsó verb-bar).      |
| `paper-notebook`      | Világos, krémszínű, mint egy papír — modern indie-stílus.              |
| `terminal-green`      | Klasszikus zöld-fekete terminál-look, retro hacker hangulat.           |

A `classic-scumm` **pontosan a jelenlegi színekkel** indul, így a meglévő
demo változtatás nélkül ugyanúgy néz ki, csak már a Theme-rendszeren át.

```go
// pncdsl/ui.theme_presets.go — a built-in presetek
func RegisterPresetThemes(g *Game) {
    g.ThemeManager.Register(Theme{Name: "classic-scumm", /* … */ })
    g.ThemeManager.Register(Theme{Name: "sierra-coin",    /* … */ })
    // …
}
```

`NewGame()` automatikusan hívja, és `UseTheme("classic-scumm")`-ot állít be —
backwards-compatible default.

## Lifecycle — pontosabban

```go
// engine.Update
func (e *engine) Update() error {
    g.input.poll()
    g.transition.update(dt)
    if g.transition.active && g.transition.out { return nil }
    if g.scriptRunner != nil { /* tick runner */ }

    uictx := &UICtx{Game: g, DT: dt}

    // 1) Top-down input feldolgozás. A felül lévő widget (reverse reg-order)
    //    kapja először az inputot; a saját Tick-jében dönthet a Consume-ról.
    for _, w := range g.UIManager.reversed() {
        w.Tick(uictx)
    }

    // 2) Ha az input még nincs konzumálva, a scene/hotspot kap esélyt:
    e.handleSceneInput()
}

// engine.Draw
func (e *engine) Draw(screen *ebiten.Image) {
    // háttér + karakterek
    g.UIManager.Each(func(w Widget) { w.Draw(screen, uictx) })
    g.transition.draw(screen, ...)
}
```

A `Cursor` widget mindig **utolsóként regisztrálva** — biztosan a tetejére
rajzol és előbb kap inputot, mint bármi más, ha kell.

A modális widgetek (`DialogBox`, `EndCard`) belsőleg figyelik az aktív
állapotot (`g.dialogueActive()`, `g.UI.endCard != ""`), és csak akkor
fogadnak inputot / rajzolnak. Ha aktívak, **konzumálnak minden inputot**,
hogy a háttér widgetek és a scene-input ne reagáljon.

## Mit kapcsol ki / mire vált ki

A jelenlegi `pncdsl/ui.def.go` központi `UI` struct **megszűnik mint
adattároló** — minden ami benne volt (`hoverLabel`, `flash`, `speech`,
`dialog`, `endCard`, `verbButtons`) **átköltözik az adott widgetbe** mint
belső állapot.

A `Game.UI` mező marad mint convenience-handle (pl. `g.UI.FlashLine(...)`
továbbra is hívható — átdelegál a `StatusLine` widgetnek a manager-en
keresztül), de ez egy vékony adapter, nem központi state.

## Fájlok

```
pncdsl/
├── ui.widget.go          # új — Widget interface, UICtx, MouseButton, Size
├── ui.manager.go         # új — UIManager alias, reversed() iterátor
├── ui.theme.go           # új — Theme struct, ThemeManager alias
├── ui.theme_presets.go   # új — classic-scumm, sierra-coin, paper-notebook, terminal-green
│
├── ui.panel.go           # új — Panel widget (generic container)
├── ui.status.go          # új — StatusLine widget (volt: ui.def hover/flash)
├── ui.verb_bar.go        # átnevezve: ui.verb.go → ui.verb_bar.go (most widget)
├── ui.verb_radial.go     # új — RadialVerbs (verb-coin)
├── ui.inventory.go       # új — InventoryBar widget (volt: ui.def-ben)
├── ui.speech.go          # → SpeechBubble widget interface-é alakítva
├── ui.dialog_box.go      # átnevezve: dialog.box.go (most widget)
├── ui.end_card.go        # új — EndCard widget
├── ui.cursor.go          # → Cursor widget interface-é alakítva
├── ui.hotspot_debug.go   # új — HotspotDebug widget (toggle-elhető)
│
└── ui.def.go             # zsugorodik egy vékony Game.UI adapterre vagy törölhető
```

Az `ui.verb_manager.go` (verb-manager alias) **változatlan marad** — a
`RadialVerbs` és a `VerbBar` is ugyanazt a `VerbManager`-t olvassa.

## Domain-oldali használat (a végén így néz ki)

```go
// domain/game.go
func Build() *p.Game {
    g := p.NewGame("Reggeli Kávé", 320, 200)

    registerAssets(g)
    definePlayer(g); defineCat(g)
    // … itemek, dialógusok, scriptek, scene-ek …

    // SCUMM-szerű alapértelmezett UI:
    g.UIManager.Register(&p.HotspotDebug{Name: "debug"})         // legalul (z=0)
    g.UIManager.Register(&p.VerbBar{Name: "verbs",
        Origin: p.Point{X: 4, Y: 152}, Cols: 2,
        Button: p.Size{W: 60, H: 14}})
    g.UIManager.Register(&p.InventoryBar{Name: "inventory",
        Origin: p.Point{X: 132, Y: 152}, Slots: 8, Cols: 8, SlotSize: 22, Gap: 2})
    g.UIManager.Register(&p.StatusLine{Name: "status", Y: 142})
    g.UIManager.Register(&p.SpeechBubble{Name: "speech"})
    g.UIManager.Register(&p.DialogBox{Name: "dialog"})
    g.UIManager.Register(&p.EndCard{Name: "endcard"})
    g.UIManager.Register(&p.Cursor{Name: "cursor"})              // legfelül

    g.UseTheme("classic-scumm")
    g.StartAt("bedroom").OnStart(p.RunScript("intro"))
    return g
}
```

Ha verb-coin kell verb-bar helyett, ennyi a változás:

```go
// töröld a VerbBar regisztrációt, helyette:
g.UIManager.Register(&p.RadialVerbs{Name: "verbs", Radius: 40})
g.UseTheme("sierra-coin")
```

És ha chat is kell:

```go
g.UIManager.Register(&ChatPanel{Name: "chat",
    Bounds: p.Rect(220, 4, 96, 130), MaxLines: 8})
```

## A „minden alapból regisztrálva" rövidítés

A fenti boilerplate cseppet sok ahhoz, hogy minden domain újra leírja. A
`pncdsl.RegisterDefaultUI(g)` egysoros hívás regisztrálja a SCUMM-szerű
default készletet:

```go
g := p.NewGame(...)
p.RegisterDefaultUI(g)   // verbs + inventory + status + speech + dialog + endcard + cursor

// a domain csak felülír / hozzáad:
g.UIManager.MustGet("verbs").(*p.VerbBar).Cols = 4
g.UIManager.Register(&ChatPanel{...})
```

A `NewGame` **nem** hívja automatikusan — explicit kell, hogy a domain
lássa hova kerülnek a default widgetek (és könnyen kivehessen onnan).

## Validáció

`Game.Validate()` kiegészül:

- Minden regisztrált `Widget.GetName()` egyedi (ezt amúgy a `Register`
  panic-cal kiszűri).
- `ActiveTheme` létezik a `ThemeManager`-ben.
- Legalább egy `SpeechBubble` és egy `DialogBox` regisztrálva, mert a
  `Say` és `RunDialogue` action ezekre számít. (Vagy: ezek a built-in actionök
  graceful no-op-pá válnak ha hiányoznak — eldöntendő.)

## Lépésenkénti implementáció

1. **Widget interface + UIManager + Theme + ThemeManager** — csontváz, még
   üres widget-ekkel. NewGame regisztrálja a manager-eket; a régi UI
   párhuzamosan él tovább.
2. **`Panel`, `StatusLine`, `SpeechBubble`, `Cursor`** — refaktoráld át
   ezeket widgetekre. A régi `UI.SetSpeech`/`UI.FlashLine` delegál a
   widgetekhez.
3. **`VerbBar` widget** — a régi `rebuildVerbButtons` + draw + hit-test
   kerül bele. Az engine `handleFreeInput` ezt használja.
4. **`InventoryBar`, `DialogBox`, `EndCard`, `HotspotDebug`** — a maradék
   beépített widgetek.
5. **`Theme` preset-ek** — minden hard-kódolt szín kikerül a Theme-be.
   `RegisterPresetThemes` regisztrálja a négyet.
6. **`RadialVerbs`** — új widget, verb-coin működéssel.
7. **`RegisterDefaultUI`** convenience + dokumentáció.
8. **Demo átállítása** — `domain/game.go` használja az új API-t. Smoke test
   marad zöld.
9. **Régi `ui.def.go` központi state kitörlése** vagy zsugorítása vékony
   adapterre.

Egy-egy lépés után a játék fut, a smoke test zöld — semmi atomi nagy commit.

## Out of scope (most)

- **Egyéni font** — a `Game.DrawText` továbbra is ebitenutil debug fontot
  használ. A Theme-ben lesz `Font` mező placeholder, de a swap később.
- **Animált widgetek** — pl. tweenelt megjelenés a `DialogBox`-nak. Felület
  van rá (Tick), implementáció későbbre.
- **Drag and drop az inventory-ban** — most csak klikk.
- **Több aktív theme egyszerre** (pl. UI overlay külön téma). Egy aktív téma.
- **Layout-szerver / declarative layout DSL** (mint a CSS flexbox). A widget
  maga felel a saját pozíciójáért. Ha kell csoportosítás, `Panel` widget +
  children API jöhet a 2. iterációban.

## Nyitott kérdés

A `Game.UI` (jelenlegi convenience-handle, `g.UI.FlashLine` stb.) maradjon,
vagy minden hívás közvetlenül `g.UIManager.MustGet("status").(*StatusLine).Flash(...)`?

Javaslom: maradjon mint **shortcut**, de mögötte a widget-rendszer fut. A
hivatkozást egy `(g *Game) Status() *StatusLine` típusú getter is megteheti
a publikus API-n, hogy ne kelljen type-assertelni a domain kódjában.
