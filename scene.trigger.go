package pncdsl

// Trigger fires Do when its When condition first becomes true after the
// trigger arms (on scene enter). One-shot per scene visit by default.
type Trigger struct {
	Name string
	When Condition
	Do   Action
	Once bool
}
