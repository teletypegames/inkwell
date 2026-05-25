package domain

import "testing"

// TestBuildValidates is the headless smoke test described in PLAN.md:
// it composes the full game and asks the library to cross-check every
// name reference between managers. No window opens.
func TestBuildValidates(t *testing.T) {
	g := Build()
	if err := g.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	for _, name := range []string{"bedroom", "kitchen"} {
		if !g.SceneManager.Has(name) {
			t.Errorf("missing scene %q", name)
		}
	}
	for _, name := range []string{"key", "beans", "mug"} {
		if !g.ItemManager.Has(name) {
			t.Errorf("missing item %q", name)
		}
	}
	if !g.DialogueManager.Has("cat_morning") {
		t.Errorf("missing dialogue cat_morning")
	}
	if !g.ScriptManager.Has("intro") || !g.ScriptManager.Has("victory") {
		t.Errorf("missing intro/victory script")
	}
}
