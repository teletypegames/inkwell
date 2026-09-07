// Package cond holds adventure-specific conditions that need *Game access
// through narrow host interfaces. Each condition type-asserts ctx.Game to the
// slice it requires.
package cond

// InventoryChecker is the slice of *Game that HasItem and SelectedItem need.
type InventoryChecker interface {
	InventoryHas(name string) bool
	InventorySelected() string
}

// SceneChecker is the slice of *Game that InScene needs.
type SceneChecker interface {
	CurrentScene() string
}

// StageChecker is the slice of *Game that OnStage, OnWorld and Started need.
type StageChecker interface {
	CurrentStagePtr() any
	OnWorldStage() bool
	Started() bool
}
