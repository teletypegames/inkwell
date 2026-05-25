package pncdsl

import "errors"

var (
	ErrUnknownAsset         = errors.New("pncdsl: unknown asset")
	ErrUnknownScene         = errors.New("pncdsl: unknown scene")
	ErrUnknownItem          = errors.New("pncdsl: unknown item")
	ErrUnknownCharacter     = errors.New("pncdsl: unknown character")
	ErrUnknownDialogue      = errors.New("pncdsl: unknown dialogue")
	ErrUnknownDialogueNode  = errors.New("pncdsl: unknown dialogue node")
	ErrUnknownScript        = errors.New("pncdsl: unknown script")
	ErrUnknownVerb          = errors.New("pncdsl: unknown verb")
	ErrDuplicateName        = errors.New("pncdsl: duplicate name")
	ErrNoStartScene         = errors.New("pncdsl: StartAt not set or unknown scene")
	ErrSceneMissingBackground = errors.New("pncdsl: scene has no background")
)
