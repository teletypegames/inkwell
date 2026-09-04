package inkwell

import (
	"sort"
	"strconv"
)

// DebugMenuRoot is the page tree of the built-in debug menu: what the game is
// made of, and what is clickable in the scene the player is standing in. It is
// data rather than a widget — the menu system is the widget, and this is one
// root a game can hand it.
var DebugMenuRoot = MenuItem{
	Label: debugTitle,
	Items: []MenuItem{
		{
			Label:  debugScenes,
			Value:  debugTally(func(ctx *Ctx) int { return ctx.Game.SceneManager.Len() }),
			Expand: debugSceneRows,
		},
		{
			Label:  debugItems,
			Value:  debugTally(func(ctx *Ctx) int { return ctx.Game.ItemManager.Len() }),
			Expand: debugItemRows,
		},
		{
			Label:  debugActionPoints,
			Value:  debugTally(debugActionPointCount),
			Expand: debugActionPointRows,
		},
		{
			Label:  debugExits,
			Value:  debugTally(debugExitCount),
			Expand: debugExitRows,
		},
	},
}

// ----- scenes ------------------------------------------------------------

func debugSceneRows(ctx *Ctx) []MenuItem {
	rows := make([]MenuItem, 0, ctx.Game.SceneManager.Len())
	for _, scene := range ctx.Game.SceneManager.All() {
		scene := scene
		rows = append(rows, MenuItem{
			Label: plain(scene.Name),
			Value: func(ctx *Ctx) Text { return debugVisits(ctx, scene.Name) },
			Do:    GoTo(scene),
		})
	}
	return rows
}

func debugVisits(ctx *Ctx, scene string) Text {
	if ctx.Game.CurrentScene() == scene {
		return debugHere
	}
	if visits := ctx.Game.State.Visited(scene); visits > 0 {
		return plain("×" + strconv.Itoa(visits))
	}
	return Text{}
}

// ----- items -------------------------------------------------------------

func debugItemRows(ctx *Ctx) []MenuItem {
	rows := make([]MenuItem, 0, ctx.Game.ItemManager.Len())
	for _, item := range ctx.Game.ItemManager.All() {
		rows = append(rows, MenuItem{
			Label:  plain(item.Name),
			Value:  debugHeldMark(item),
			Expand: debugItemPage(item),
		})
	}
	return rows
}

func debugItemPage(item *Item) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		rows := []MenuItem{
			{
				Label:    debugGive,
				Do:       Give(item),
				Enabled:  Not(HasItem(item)),
				KeepOpen: true,
			},
			{
				Label:    debugTakeAway,
				Do:       TakeAway(item),
				Enabled:  HasItem(item),
				KeepOpen: true,
			},
			{
				Label:   debugSelect,
				Do:      debugSelectItem(item),
				Enabled: HasItem(item),
			},
		}
		if item.OnUseSelf != nil {
			rows = append(rows, MenuItem{
				Label: debugUseSelf,
				Do:    item.OnUseSelf,
			})
		}
		if len(item.OnUseWith) > 0 {
			rows = append(rows, MenuItem{
				Label:  debugUseWith,
				Value:  debugTally(func(*Ctx) int { return len(item.OnUseWith) }),
				Expand: debugItemPairRows(item),
			})
		}
		return rows
	}
}

func debugItemPairRows(item *Item) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		rows := make([]MenuItem, 0, len(item.OnUseWith))
		for _, target := range debugSortedItems(item.OnUseWith) {
			rows = append(rows, MenuItem{
				Label: plain(target.Name),
				Do:    item.OnUseWith[target],
			})
		}
		return rows
	}
}

func debugSelectItem(item *Item) Action {
	return Custom(func(ctx *Ctx) Status {
		ctx.Game.Inventory.Select(item.Name)
		return StatusDone
	})
}

func debugHeldMark(item *Item) func(*Ctx) Text {
	return func(ctx *Ctx) Text {
		if !ctx.Game.Inventory.Has(item.Name) {
			return Text{}
		}
		if ctx.Game.Inventory.Selected() == item.Name {
			return debugSelected
		}
		return debugHeld
	}
}

// ----- action points -----------------------------------------------------

func debugActionPointRows(ctx *Ctx) []MenuItem {
	scene := ctx.Game.CurrentScene()
	hotspots := ctx.Game.SceneHotspots(scene)
	rows := make([]MenuItem, 0, len(hotspots))
	for _, hotspot := range hotspots {
		name := hotspot.Name
		rows = append(rows, MenuItem{
			Label:  plain(name),
			Value:  FixedValue(plain(debugCursorName(hotspot.Cursor))),
			Expand: debugVerbRows(scene, name),
		})
	}
	return rows
}

func debugActionPointCount(ctx *Ctx) int {
	return len(ctx.Game.SceneHotspots(ctx.Game.CurrentScene()))
}

func debugVerbRows(scene, hotspot string) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		spot := debugHotspot(ctx.Game, scene, hotspot)
		if spot == nil {
			return nil
		}
		var rows []MenuItem
		for _, verb := range ctx.Game.VerbManager.All() {
			if spot.handler(verb.Name) == nil {
				continue
			}
			rows = append(rows, MenuItem{
				Label: plain(verb.Name),
				Do:    debugInvoke(scene, hotspot, verb.Name),
			})
		}
		if len(spot.OnUseWith) > 0 {
			rows = append(rows, MenuItem{
				Label:  debugUseWith,
				Value:  debugTally(func(*Ctx) int { return len(spot.OnUseWith) }),
				Expand: debugPairRows(scene, hotspot),
			})
		}
		return rows
	}
}

func debugPairRows(scene, hotspot string) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		spot := debugHotspot(ctx.Game, scene, hotspot)
		if spot == nil {
			return nil
		}
		rows := make([]MenuItem, 0, len(spot.OnUseWith))
		for _, item := range debugSortedItems(spot.OnUseWith) {
			rows = append(rows, MenuItem{
				Label: plain(item.Name),
				Value: debugHeldMark(item),
				Do:    debugInvokeWith(scene, hotspot, item),
			})
		}
		return rows
	}
}

func debugInvoke(scene, hotspot, verb string) Action {
	return Custom(func(ctx *Ctx) Status {
		if spot := debugHotspot(ctx.Game, scene, hotspot); spot != nil {
			ctx.Game.InvokeVerb(spot, verb)
		}
		return StatusDone
	})
}

func debugInvokeWith(scene, hotspot string, item *Item) Action {
	return Custom(func(ctx *Ctx) Status {
		g := ctx.Game
		g.Inventory.Add(item.Name)
		g.Inventory.Select(item.Name)
		if spot := debugHotspot(g, scene, hotspot); spot != nil {
			g.InvokeVerb(spot, "use")
		}
		return StatusDone
	})
}

func debugHotspot(g *Game, scene, name string) *Hotspot {
	hotspots := g.SceneHotspots(scene)
	for i := range hotspots {
		if hotspots[i].Name == name {
			return &hotspots[i]
		}
	}
	return nil
}

func debugCursorName(kind CursorKind) string {
	switch kind {
	case CursorLook:
		return "look"
	case CursorUse:
		return "use"
	case CursorTalk:
		return "talk"
	case CursorTake:
		return "take"
	case CursorExit:
		return "exit"
	default:
		return ""
	}
}

// ----- exits -------------------------------------------------------------

func debugExitRows(ctx *Ctx) []MenuItem {
	scene, ok := ctx.Game.SceneManager.Get(ctx.Game.CurrentScene())
	if !ok {
		return nil
	}
	rows := make([]MenuItem, 0, len(scene.Exits))
	for _, exit := range scene.Exits {
		rows = append(rows, MenuItem{
			Label:  plain(ExitName(exit.To)),
			Value:  FixedValue(plain(debugSideName(exit.Side))),
			Expand: debugExitPage(exit),
		})
	}
	return rows
}

func debugExitPage(exit Exit) func(*Ctx) []MenuItem {
	return func(*Ctx) []MenuItem {
		travel := GoTo(exit.To)
		if exit.To == nil {
			travel = Back()
		}
		rows := []MenuItem{{
			Label: debugGo,
			Do:    travel,
		}}
		if exit.Needs == "" {
			return rows
		}
		flag := exit.Needs
		return append(rows, MenuItem{
			Label: Format(debugNeeds, plain(flag)),
			Value: func(ctx *Ctx) Text { return debugFlagMark(ctx, flag) },
			Items: []MenuItem{
				{
					Label:    debugSet,
					Do:       SetFlag(flag),
					Enabled:  Not(Flag(flag)),
					KeepOpen: true,
				},
				{
					Label:    debugClear,
					Do:       ClearFlag(flag),
					Enabled:  Flag(flag),
					KeepOpen: true,
				},
			},
		})
	}
}

func debugExitCount(ctx *Ctx) int {
	scene, ok := ctx.Game.SceneManager.Get(ctx.Game.CurrentScene())
	if !ok {
		return 0
	}
	return len(scene.Exits)
}

func debugFlagMark(ctx *Ctx, flag string) Text {
	if ctx.Game.State.Flag(flag) {
		return debugOn
	}
	return debugOff
}

func debugSideName(side ExitSide) string {
	switch side {
	case ExitLeft:
		return "left"
	case ExitRight:
		return "right"
	case ExitBack:
		return "back"
	default:
		return "near"
	}
}

// ----- helpers -----------------------------------------------------------

func debugTally(count func(*Ctx) int) func(*Ctx) Text {
	return func(ctx *Ctx) Text { return plain(strconv.Itoa(count(ctx))) }
}

func debugSortedItems(m map[*Item]Action) []*Item {
	items := make([]*Item, 0, len(m))
	for item := range m {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items
}

var (
	debugTitle = Text{
		{Lang: "en", Text: "DEBUG"},
		{Lang: "hu", Text: "DEBUG"},
	}
	debugScenes = Text{
		{Lang: "en", Text: "Scenes"},
		{Lang: "hu", Text: "Jelenetek"},
	}
	debugItems = Text{
		{Lang: "en", Text: "Items"},
		{Lang: "hu", Text: "Tárgyak"},
	}
	debugActionPoints = Text{
		{Lang: "en", Text: "Action points"},
		{Lang: "hu", Text: "Akciópontok"},
	}
	debugExits = Text{
		{Lang: "en", Text: "Exits"},
		{Lang: "hu", Text: "Kijáratok"},
	}
	debugHere = Text{
		{Lang: "en", Text: "here"},
		{Lang: "hu", Text: "itt"},
	}
	debugHeld = Text{
		{Lang: "en", Text: "held"},
		{Lang: "hu", Text: "nálam"},
	}
	debugSelected = Text{
		{Lang: "en", Text: "in hand"},
		{Lang: "hu", Text: "kézben"},
	}
	debugGive = Text{
		{Lang: "en", Text: "give"},
		{Lang: "hu", Text: "add ide"},
	}
	debugTakeAway = Text{
		{Lang: "en", Text: "take away"},
		{Lang: "hu", Text: "vedd el"},
	}
	debugSelect = Text{
		{Lang: "en", Text: "select"},
		{Lang: "hu", Text: "válaszd ki"},
	}
	debugUseSelf = Text{
		{Lang: "en", Text: "use on self"},
		{Lang: "hu", Text: "használd magában"},
	}
	debugUseWith = Text{
		{Lang: "en", Text: "use with"},
		{Lang: "hu", Text: "használd ezzel"},
	}
	debugGo = Text{
		{Lang: "en", Text: "go"},
		{Lang: "hu", Text: "menj"},
	}
	debugNeeds = Text{
		{Lang: "en", Text: "needs %s"},
		{Lang: "hu", Text: "kell hozzá: %s"},
	}
	debugSet = Text{
		{Lang: "en", Text: "set"},
		{Lang: "hu", Text: "állítsd be"},
	}
	debugClear = Text{
		{Lang: "en", Text: "clear"},
		{Lang: "hu", Text: "töröld"},
	}
	debugOn = Text{
		{Lang: "en", Text: "on"},
		{Lang: "hu", Text: "be"},
	}
	debugOff = Text{
		{Lang: "en", Text: "off"},
		{Lang: "hu", Text: "ki"},
	}
)
