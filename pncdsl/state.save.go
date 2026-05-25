package pncdsl

// Save/Load are stubbed out in this milestone — the runtime state machine
// is in place, but JSON serialization will land alongside the polish pass.
func (g *Game) Save(slot int) error { _ = slot; return nil }
func (g *Game) Load(slot int) error { _ = slot; return nil }
