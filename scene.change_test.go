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
