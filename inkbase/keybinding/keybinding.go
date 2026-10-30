package keybinding

import (
	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/inkbase/input"
	"git.teletypegames.org/engines/inkwell/inkbase/widget"
)

type host interface {
	WidgetInput() *input.Input
	DoAction(inkbase.Action)
}

// KeyBinding is a key and the action it runs, as a widget: it draws nothing
// and occupies nothing, so binding a key is a registration like any other.
// The key is claimed, so nothing underneath sees it, and When gates it the
// way it gates any widget.
type KeyBinding struct {
	// Name is the key the widget is registered under.
	Name string
	// When gates the binding; nil keeps it live always.
	When inkbase.Condition
	// Key is the key that triggers it; zero binds nothing.
	Key ebiten.Key
	// Do is the action run on a press; nil does nothing.
	Do inkbase.Action
	// DuringScript lets the binding answer while a script holds the game.
	DuringScript bool
}

// GetName returns the binding's Name, the key it is registered under.
func (k *KeyBinding) GetName() string { return k.Name }

// VisibleWhen returns When, the condition that keeps the binding live; nil
// keeps it live always.
func (k *KeyBinding) VisibleWhen() inkbase.Condition { return k.When }

// Layer puts the binding on [widget.LayerMenu], so it is ticked among the
// menus.
func (k *KeyBinding) Layer() widget.Layer { return widget.LayerMenu }

// TicksDuringScript reports DuringScript.
func (k *KeyBinding) TicksDuringScript() bool { return k.DuringScript }

// Draw does nothing: a binding has no picture.
func (k *KeyBinding) Draw(dst *ebiten.Image, ctx *widget.UICtx) {}

// Tick claims Key when it was pressed this frame and hands Do to the host
// to run. It does nothing on a host without input and action dispatch.
func (k *KeyBinding) Tick(ctx *widget.UICtx) {
	if k.Key == 0 || k.Do == nil {
		return
	}
	h, ok := ctx.Game.(host)
	if !ok {
		return
	}
	if h.WidgetInput().ClaimKey(k.Key) {
		h.DoAction(k.Do)
	}
}
