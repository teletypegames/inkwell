package keybinding

import (
	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/input"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

// Host is what a [KeyBinding] needs from the game it is ticked for, the
// UICtx's Game. On a Game that does not implement it the binding is inert.
type Host interface {
	// WidgetInput is the frame's consume-on-use input; the binding claims
	// its key.
	WidgetInput() *input.Input
	// Do queues the binding's action.
	Do(a inkbase.Action)
}

// KeyBinding is a key and what it does — the whole widget. It draws nothing
// and occupies nothing; it exists so that "Escape opens the main menu" is a
// registration like everything else rather than a special case somewhere in
// the loop.
//
// The key is claimed, so nothing underneath sees it. When gates it the way it
// gates any widget, which is how a binding belongs to one Stage.
type KeyBinding struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the binding; nil keeps it live always.
	When inkbase.Condition
	// Key is the key that triggers it; zero binds nothing.
	Key ebiten.Key
	// Do is the action run on a press; nil does nothing.
	Do inkbase.Action
	// DuringScript lets the binding answer on the frames a running script
	// holds — what a pause key needs, since a cutscene is exactly when the
	// player reaches for one.
	DuringScript bool
}

// GetName returns the name the binding is registered under.
func (k *KeyBinding) GetName() string { return k.Name }

// VisibleWhen returns the condition that switches the binding on; nil means
// always.
func (k *KeyBinding) VisibleWhen() inkbase.Condition { return k.When }

// Layer puts the binding on [widget.LayerMenu], so it claims its key
// before the widgets under the menus see it.
func (k *KeyBinding) Layer() widget.Layer { return widget.LayerMenu }

// TicksDuringScript reports DuringScript.
func (k *KeyBinding) TicksDuringScript() bool { return k.DuringScript }

// Draw does nothing; a binding has no picture.
func (k *KeyBinding) Draw(dst *ebiten.Image, ctx *widget.UICtx) {}

// Tick claims Key when it was pressed this frame and hands Do to the game.
// A zero Key or a nil Do makes the binding inert.
func (k *KeyBinding) Tick(ctx *widget.UICtx) {
	if k.Key == 0 || k.Do == nil {
		return
	}
	h, ok := ctx.Game.(Host)
	if !ok {
		return
	}
	if h.WidgetInput().ClaimKey(k.Key) {
		h.Do(k.Do)
	}
}
