package inkwell

import "errors"

var (
	ErrUnknownAsset     = errors.New("inkwell: unknown asset")
	ErrUnknownScene     = errors.New("inkwell: unknown scene")
	ErrUnknownCharacter = errors.New("inkwell: unknown character")
	ErrUnknownScript    = errors.New("inkwell: unknown script")
	ErrNoStartScene     = errors.New("inkwell: StartScene not set or unknown scene")
)
