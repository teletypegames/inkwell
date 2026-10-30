package inkwell

import "errors"

// The errors the engine returns; compare with [errors.Is], since most are
// wrapped with the name that was not found.
var (
	// ErrUnknownAsset is a name no registered Asset has.
	ErrUnknownAsset = errors.New("inkwell: unknown asset")
	// ErrUnknownScene is a name no registered Scene has.
	ErrUnknownScene = errors.New("inkwell: unknown scene")
	// ErrUnknownCharacter is a name no registered Character has.
	ErrUnknownCharacter = errors.New("inkwell: unknown character")
	// ErrUnknownScript is a name no registered Script has.
	ErrUnknownScript = errors.New("inkwell: unknown script")
	// ErrUnknownStage is a name no registered Stage has.
	ErrUnknownStage = errors.New("inkwell: unknown stage")
	// ErrNoStartScene is a game whose StartScene is nil or not registered.
	ErrNoStartScene = errors.New("inkwell: StartScene not set or unknown scene")
	// ErrNoStartStage is a game whose StartStage is nil or not registered.
	ErrNoStartStage = errors.New("inkwell: StartStage not set or unknown stage")
	// ErrNotStored is what a Store's Read returns for a key nothing has
	// written yet.
	ErrNotStored = errors.New("inkwell: nothing stored under this key")
)
