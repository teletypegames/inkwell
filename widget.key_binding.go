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

func (k *KeyBinding) GetName() string         { return k.Name }
func (k *KeyBinding) VisibleWhen() Condition  { return k.When }
func (k *KeyBinding) Layer() Layer            { return LayerMenu }
func (k *KeyBinding) TicksDuringScript() bool { return k.DuringScript }

func (k *KeyBinding) Draw(dst *ebiten.Image, ctx *UICtx) {}

func (k *KeyBinding) Tick(ctx *UICtx) {
	if k.Key == 0 || k.Do == nil {
		return
	}
	g := gameOfUI(ctx)
	if g.Input.ClaimKey(k.Key) {
		g.Do(k.Do)
	}
}
