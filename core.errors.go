package inkwell

import "errors"

var (
	ErrUnknownAsset         = errors.New("inkwell: unknown asset")
	ErrUnknownScene         = errors.New("inkwell: unknown scene")
	ErrUnknownItem          = errors.New("inkwell: unknown item")
	ErrUnknownCharacter     = errors.New("inkwell: unknown character")
	ErrUnknownDialogue      = errors.New("inkwell: unknown dialogue")
	ErrUnknownDialogueNode  = errors.New("inkwell: unknown dialogue node")
	ErrUnknownScript        = errors.New("inkwell: unknown script")
	ErrUnknownVerb          = errors.New("inkwell: unknown verb")
	ErrDuplicateName        = errors.New("inkwell: duplicate name")
	ErrNoStartScene         = errors.New("inkwell: StartAt not set or unknown scene")
	ErrSceneMissingBackground = errors.New("inkwell: scene has no background")
)
