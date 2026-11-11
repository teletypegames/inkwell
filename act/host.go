package act

import (
	"git.teletypegames.org/engines/inkwell/dialog"
	"git.teletypegames.org/engines/inkwell/inkbase"
	"git.teletypegames.org/engines/inkwell/sprite"
)

// SpeechHost is the slice of *Game a Say or Reply action needs.
type SpeechHost interface {
	NoteSpoken(name string)
	DeliveryOf(name string) inkbase.Delivery
	SetSpeech(speaker, text string)
	ClearSpeech()
	LogResponse(speaker, text string)
	Talking() string
	SetTalking(name string)
	VoiceOn() bool
	PlayVoice(v *inkbase.Voice) float64
	StopVoice()
	ConsumedClick() bool
}

// InventoryHost is the slice of *Game an inventory action needs.
type InventoryHost interface {
	InventoryAdd(name string)
	InventoryRemove(name string)
	InventoryHas(name string) bool
	FlashLine(text inkbase.Text)
	ItemLabel(name string) inkbase.Text
}

// DialogHost is the slice of *Game a dialog action needs.
type DialogHost interface {
	StartDialog(d *dialog.Dialog)
	DialogActive() bool
	EndDialog()
	GotoDialogNode(name string)
}

// SceneHost is the slice of *Game a scene-change action needs.
type SceneHost interface {
	EnterScene(name string) inkbase.Runner
	PreviousScene() string
}

// MovementHost is the slice of *Game a Walk action needs.
type MovementHost interface {
	WalkCharacter(name string, to inkbase.Point)
	CharacterMoving(name string) bool
}

// PoseHost is the slice of *Game a SetPose action needs.
type PoseHost interface {
	SetCharacterPose(name string, pose sprite.Pose)
}

// VisibilityHost is the slice of *Game a Hide or Show action needs.
type VisibilityHost interface {
	SetCharacterHidden(name string, hidden bool)
}

// VoiceHost is the slice of *Game a SetVoice action needs.
type VoiceHost interface {
	SetVoiceEnabled(on bool)
	StopVoice()
	SaveSettings() error
}

// EndCardHost is the slice of *Game a ShowEnd action needs.
type EndCardHost interface {
	ShowEndCard(text string)
}

// FlashHost is the slice of *Game a Flash action needs.
type FlashHost interface {
	FlashLine(text inkbase.Text)
}
