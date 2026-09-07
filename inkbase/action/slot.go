package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// SlotHost provides save and load.
type SlotHost interface {
	LoadLater(slot int)
	Save(slot int) error
}

type loadSlotAction struct{ slot int }

// LoadSlot restores a saved session.
func LoadSlot(slot int) inkbase.Action { return &loadSlotAction{slot: slot} }

func (a *loadSlotAction) Start() inkbase.Runner { return a }

func (a *loadSlotAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(SlotHost); ok {
		h.LoadLater(a.slot)
	}
	return inkbase.StatusDone
}

type saveSlotAction struct{ slot int }

// SaveSlot writes the session to a slot.
func SaveSlot(slot int) inkbase.Action { return &saveSlotAction{slot: slot} }

func (a *saveSlotAction) Start() inkbase.Runner { return a }

func (a *saveSlotAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	if h, ok := ctx.Game.(SlotHost); ok {
		if err := h.Save(a.slot); err != nil {
			return inkbase.StatusFailed
		}
	}
	return inkbase.StatusDone
}
