package inkwell

type loadSlotAction struct{ slot int }

// LoadSlot restores a saved session. It lands at the top of the next frame,
// for the same reason Restart does: replacing State, the inventory and every
// character underneath a running action is not something the running action
// should have to survive.
func LoadSlot(slot int) Action { return &loadSlotAction{slot: slot} }

func (a *loadSlotAction) Start() Runner { return a }

func (a *loadSlotAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	game.LoadLater(a.slot)
	return StatusDone
}

type saveSlotAction struct{ slot int }

// SaveSlot writes the session to a slot. Saving touches nothing, so unlike
// LoadSlot it happens where it is called.
func SaveSlot(slot int) Action { return &saveSlotAction{slot: slot} }

func (a *saveSlotAction) Start() Runner { return a }

func (a *saveSlotAction) Tick(ctx *Ctx) Status {
	game := gameOf(ctx)
	if err := game.Save(a.slot); err != nil {
		logf("saveSlot %d: %v", a.slot, err)
		return StatusFailed
	}
	return StatusDone
}
