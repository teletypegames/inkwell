package inkwell

import "github.com/hajimehoshi/ebiten/v2"

// KeyBinding is a key and what it does — the whole widget. It draws nothing
// and occupies nothing; it exists so that "Escape opens the main menu" is a
// registration like everything else rather than a special case somewhere in
// the loop.
//
// The key is claimed, so nothing underneath sees it. When gates it the way it
// gates any widget, which is how a binding belongs to one Stage.
type KeyBinding struct {
	Name string
	When Condition
	Key  ebiten.Key
	Do   Action

	// DuringScript lets the binding answer on the frames a running script
	// holds — what a pause key needs, since a cutscene is exactly when the
	// player reaches for one.
	DuringScript bool
}

// GetName returns the name the binding is registered under.
func (k *KeyBinding) GetName() string { return k.Name }

// VisibleWhen returns the condition that switches the binding on; nil means
// always.
func (k *KeyBinding) VisibleWhen() Condition { return k.When }

// Layer puts the binding on [LayerMenu], so it claims its key before the
// widgets under the menus see it.
func (k *KeyBinding) Layer() Layer { return LayerMenu }

// TicksDuringScript reports DuringScript.
func (k *KeyBinding) TicksDuringScript() bool { return k.DuringScript }

// Draw does nothing; a binding has no picture.
func (k *KeyBinding) Draw(dst *ebiten.Image, ctx *UICtx) {}

// Tick claims Key when it was pressed this frame and hands Do to the game.
// A zero Key or a nil Do makes the binding inert.
func (k *KeyBinding) Tick(ctx *UICtx) {
	if k.Key == 0 || k.Do == nil {
		return
	}
	g := gameOfUI(ctx)
	if g.Input.ClaimKey(k.Key) {
		g.Do(k.Do)
	}
}
