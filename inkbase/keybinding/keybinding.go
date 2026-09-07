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

type KeyBinding struct {
	Name         string
	When         inkbase.Condition
	Key          ebiten.Key
	Do           inkbase.Action
	DuringScript bool
}

func (k *KeyBinding) GetName() string                { return k.Name }
func (k *KeyBinding) VisibleWhen() inkbase.Condition { return k.When }
func (k *KeyBinding) Layer() widget.Layer            { return widget.LayerMenu }
func (k *KeyBinding) TicksDuringScript() bool        { return k.DuringScript }

func (k *KeyBinding) Draw(dst *ebiten.Image, ctx *widget.UICtx) {}

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
