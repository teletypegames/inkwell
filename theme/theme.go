package theme

import "image/color"

// Theme is a registry-stored entity carrying every color the built-in
// widgets read. Multiple themes can be registered; the active one is
// addressed by name on the Game.
type Theme struct {
	Name string

	PanelBG    color.Color
	StatusText color.Color
	FlashText  color.Color

	VerbButtonBG         color.Color
	VerbButtonSelectedBG color.Color
	VerbButtonText       color.Color

	InventorySlotBG         color.Color
	InventorySlotSelectedBG color.Color

	SpeechBubbleBG    color.Color
	SpeechDefaultText color.Color

	DialogBG          color.Color
	DialogBorder      color.Color
	DialogChoiceBG    color.Color
	DialogChoiceHover color.Color
	DialogSpeaker     color.Color
	DialogText        color.Color

	EndCardBG   color.Color
	EndCardText color.Color

	CursorColor    color.Color
	HotspotOutline color.Color

	// SceneBackdrop is painted before the scene background. Useful for
	// letterboxing or filling areas not covered by a smaller backdrop.
	SceneBackdrop color.Color

	// Rich-UI widgets (TopBar, CharacterPanel, ChatLog).
	TopBarBG     color.Color
	TopBarText   color.Color
	TopBarAccent color.Color // score / highlighted number color

	ChatLogBG       color.Color
	ChatLogPrompt   color.Color // "> look at door" lines
	ChatLogResponse color.Color // in-world reply lines
	ChatLogSystem   color.Color // system / meta lines

	CharacterPanelBG     color.Color
	CharacterPanelBorder color.Color
	CharacterPanelTitle  color.Color // "PLAYER" / "NPC" badge color
}

func (t Theme) GetName() string   { return t.Name }
func (t Theme) TypeLabel() string { return "theme" }
