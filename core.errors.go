package inkwell

import "errors"

var (
	ErrUnknownAsset     = errors.New("inkwell: unknown asset")
	ErrUnknownScene     = errors.New("inkwell: unknown scene")
	ErrUnknownCharacter = errors.New("inkwell: unknown character")
	ErrUnknownScript    = errors.New("inkwell: unknown script")
	ErrUnknownStage     = errors.New("inkwell: unknown stage")
	ErrNoStartScene     = errors.New("inkwell: StartScene not set or unknown scene")
	ErrNoStartStage     = errors.New("inkwell: StartStage not set or unknown stage")
	ErrNotStored        = errors.New("inkwell: nothing stored under this key")
)
