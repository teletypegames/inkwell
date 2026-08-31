package inkwell

import (
	"sort"
	"strconv"
)

// DebugMenuRoot is the page tree of the built-in debug menu: what the game is
// made of, and what is clickable in the scene the player is standing in. It is
// data rather than a widget — the menu system is the widget, and this is one
// use of it:
//
//	WidgetManager.Register(&inkwell.Menu{
//		Name:      "debug_menu",
//		Root:      inkwell.DebugMenuRoot,
//		When:      whenDebug,
//		ToggleKey: ebiten.KeyD,
//	})
//
// Its rows name registry keys, not the prose a player reads: a debug menu is
// read by whoever writes the content, and `noodle_house` says more there than
// "NOODLE'S HOUSE — KITCHEN" does, in any language.
//
// Nothing registers it, RegisterDefaultUI included. A game that wants it says
// so, and gates it on whatever its debug builds set.
//
// A game can add sections of its own by building a Root around it:
//
//	Root: inkwell.MenuItem{
//		Label: inkwell.DebugMenuRoot.Label,
//		Items: append(inkwell.DebugMenuRoot.Items, mySection),
//	}
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

// ----- scenes -----------------------------------------------------------

// debugSceneRows is every registered scene, in registration order, and
// entering one teleports there. The value column marks the scene the player
// is in and counts the visits to the others.
func debugSceneRows(ctx *Ctx) []MenuItem {
	rows := make([]MenuItem, 0, ctx.Game.SceneManager.Len())
	for _, name := range ctx.Game.SceneManager.Names() {
		scene := name
		rows = append(rows, MenuItem{
			Label: plain(scene),
			Value: func(ctx *Ctx) Text { return debugVisits(ctx, scene) },
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

// ----- items ------------------------------------------------------------

// debugItemRows is every registered item, and each one opens the handful of
// things that can be done to it from outside the fiction.
func debugItemRows(ctx *Ctx) []MenuItem {
	rows := make([]MenuItem, 0, ctx.Game.ItemManager.Len())
	for _, item := range ctx.Game.ItemManager.All() {
		rows = append(rows, MenuItem{
			Label:  plain(item.Name),
			Value:  debugHeldMark(item.Name),
			Expand: debugItemPage(item.Name),
		})
	}
	return rows
}

func debugItemPage(name string) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		rows := []MenuItem{
			{
				Label:    debugGive,
				Do:       Give(name),
				Enabled:  Not(HasItem(name)),
				KeepOpen: true,
			},
			{
				Label:    debugTakeAway,
				Do:       TakeAway(name),
				Enabled:  HasItem(name),
				KeepOpen: true,
			},
			{
				Label:   debugSelect,
				Do:      debugSelectItem(name),
				Enabled: HasItem(name),
			},
		}
		item, ok := ctx.Game.ItemManager.Get(name)
		if !ok {
			return rows
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
				Expand: debugItemPairRows(name),
			})
		}
		return rows
	}
}

// debugItemPairRows is the targets an item is paired with, straight out of its
// own OnUseWith. Entering one runs the pairing without walking anybody up to
// anything.
func debugItemPairRows(name string) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		item, ok := ctx.Game.ItemManager.Get(name)
		if !ok {
			return nil
		}
		rows := make([]MenuItem, 0, len(item.OnUseWith))
		for _, target := range debugSortedKeys(item.OnUseWith) {
			rows = append(rows, MenuItem{
				Label: plain(target),
				Do:    item.OnUseWith[target],
			})
		}
		return rows
	}
}

func debugSelectItem(name string) Action {
	return Custom(func(ctx *Ctx) Status {
		ctx.Game.Inventory.Select(name)
		return StatusDone
	})
}

func debugHeldMark(name string) func(*Ctx) Text {
	return func(ctx *Ctx) Text {
		if !ctx.Game.Inventory.Has(name) {
			return Text{}
		}
		if ctx.Game.Inventory.Selected() == name {
			return debugSelected
		}
		return debugHeld
	}
}

// ----- action points ----------------------------------------------------

// debugActionPointRows is everything clickable in the scene the player is
// standing in — the authored hotspots and the ones generated from the exits,
// in the order the engine hit-tests them. Each one opens the verbs it actually
// has an answer for.
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

// debugVerbRows is the verbs bound on one hotspot, plus the items it is paired
// with. Every one of them goes through invokeHotspotVerb, the same path a
// click takes, so the chat log and the status line read as if the player had
// done it.
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

// debugPairRows is the items a hotspot is paired with. Entering one hands the
// item over first — the point is to see the pairing run, not to be reminded
// that the player never picked the thing up.
func debugPairRows(scene, hotspot string) func(*Ctx) []MenuItem {
	return func(ctx *Ctx) []MenuItem {
		spot := debugHotspot(ctx.Game, scene, hotspot)
		if spot == nil {
			return nil
		}
		rows := make([]MenuItem, 0, len(spot.OnUseWith))
		for _, item := range debugSortedKeys(spot.OnUseWith) {
			rows = append(rows, MenuItem{
				Label: plain(item),
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
			ctx.Game.invokeHotspotVerb(spot, verb)
		}
		return StatusDone
	})
}

func debugInvokeWith(scene, hotspot, item string) Action {
	return Custom(func(ctx *Ctx) Status {
		g := ctx.Game
		g.Inventory.Add(item)
		g.Inventory.Select(item)
		if spot := debugHotspot(g, scene, hotspot); spot != nil {
			g.invokeHotspotVerb(spot, "use")
		}
		return StatusDone
	})
}

// debugHotspot finds a hotspot by name in a scene's expanded list. The rows
// hold names rather than pointers, so a page that outlived a scene change
// answers nil instead of poking at a stale hotspot.
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

// ----- exits ------------------------------------------------------------

// debugExitRows is the current scene's exits as the graph sees them: where each
// one leads, which edge it sits on, and — when it is locked — the flag that
// opens it, with the two rows that set and clear it.
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
		if exit.To == "" {
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

// ----- shared -----------------------------------------------------------

// debugTally is a value column that is just a number.
func debugTally(count func(*Ctx) int) func(*Ctx) Text {
	return func(ctx *Ctx) Text { return plain(strconv.Itoa(count(ctx))) }
}

// debugSortedKeys keeps a map-backed page in the same order from one frame to
// the next, which a Go map on its own does not.
func debugSortedKeys(m map[string]Action) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// The debug menu's own wording. It stays here rather than in i18n.strings.go
// with the rest of the engine's voice: this is a tool for whoever builds the
// game, not a line the player is ever shown, so it is not part of what a
// translator is handed.
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
