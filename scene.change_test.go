package inkwell

import (
	"reflect"
	"testing"
)

func TestGoToHoldsTheScriptAcrossScenes(t *testing.T) {
	g := NewGame("t", 320, 200)
	var log []string
	note := func(s string) Action {
		return &immediateAction{fn: func(*Ctx) Status {
			log = append(log, s+"@"+g.currentScene)
			return StatusDone
		}}
	}
	a := &Scene{Name: "a", OnLeave: note("leaveA")}
	b := &Scene{Name: "b", OnEnter: note("enterB"), OnLeave: note("leaveB")}
	c := &Scene{Name: "c", OnEnter: note("enterC")}
	for _, s := range []*Scene{a, b, c} {
		g.SceneManager.Register(s)
	}
	g.currentScene = "a"

	r := Seq(GoTo(b), note("line1"), GoTo(c), note("line2"), Back(), note("line3")).Start()
	ctx := g.makeCtx()
	done := false
	for i := 0; i < 1000 && !done; i++ {
		g.transition.update(1.0 / 60)
		if g.transition.active && g.transition.out {
			continue
		}
		done = r.Tick(ctx) != StatusRunning
	}

	want := []string{"leaveA@b", "enterB@b", "line1@b", "leaveB@c", "enterC@c", "line2@c", "enterB@b", "line3@b"}
	if !done || !reflect.DeepEqual(log, want) {
		t.Fatalf("done=%v log=%v, want %v", done, log, want)
	}
	if len(g.queue) != 0 {
		t.Fatalf("hooks were queued: %d", len(g.queue))
	}
}

func TestHideLastsAcrossScenesAndDrawsOnlyTheCast(t *testing.T) {
	g := NewGame("t", 320, 200)
	paul := &Character{Name: "paul"}
	guard := &Character{Name: "guard"}
	g.CharacterManager.Register(paul)
	g.CharacterManager.Register(guard)
	g.Player = paul
	street := &Scene{Name: "street", Actors: []SceneActor{{Character: paul}, {Character: guard}}}
	roof := &Scene{Name: "roof"}
	mapScreen := &Scene{Name: "map", Actors: []SceneActor{}}
	for _, s := range []*Scene{street, roof, mapScreen} {
		g.SceneManager.Register(s)
	}
	g.currentScene = "street"
	g.placeActors("street")
	names := func(s *Scene) []string {
		var out []string
		for _, c := range sortedChars(g, s) {
			out = append(out, c.def.Name)
		}
		return out
	}

	g.placeActors("map")
	if got := names(mapScreen); len(got) != 0 {
		t.Fatalf("map draws %v, want nobody", got)
	}

	r := Seq(Hide(paul), GoTo(roof)).Start()
	ctx := g.makeCtx()
	for i := 0; i < 1000 && r.Tick(ctx) == StatusRunning; i++ {
		g.transition.update(1.0 / 60)
	}
	if got := names(roof); len(got) != 0 || g.CharacterInScene(paul) {
		t.Fatalf("hidden paul drawn on roof: %v", got)
	}

	g.SetCharacterHidden("paul", false)
	if got := names(roof); !reflect.DeepEqual(got, []string{"paul"}) {
		t.Fatalf("roof draws %v after Show, want [paul]", got)
	}

	g.SetCharacterHidden("paul", true)
	sf := g.buildSave()
	g.hidden = map[string]bool{}
	if err := g.applySave(&sf); err != nil {
		t.Fatal(err)
	}
	if !g.CharacterHidden("paul") {
		t.Fatal("hidden lost through a save")
	}
}
