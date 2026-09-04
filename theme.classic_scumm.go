package inkwell

import "image/color"

// classic-scumm is the default: the dark blue-grey plate and amber
// highlight of the verb-bar era.

func classicScummTheme() *Theme {
	return &Theme{
		Name:                    "classic-scumm",
		PanelBG:                 color.RGBA{20, 20, 30, 255},
		StatusText:              color.White,
		FlashText:               color.RGBA{255, 220, 120, 255},
		VerbButtonBG:            color.RGBA{50, 50, 70, 255},
		VerbButtonSelectedBG:    color.RGBA{120, 90, 40, 255},
		VerbButtonText:          color.White,
		InventorySlotBG:         color.RGBA{40, 40, 50, 255},
		InventorySlotSelectedBG: color.RGBA{120, 90, 40, 255},
		SpeechBubbleBG:          color.RGBA{0, 0, 0, 200},
		SpeechDefaultText:       color.White,
		DialogBG:                color.RGBA{15, 15, 25, 240},
		DialogBorder:            color.RGBA{80, 80, 110, 255},
		DialogChoiceBG:          color.RGBA{30, 30, 45, 255},
		DialogChoiceHover:       color.RGBA{70, 70, 100, 255},
		DialogSpeaker:           color.RGBA{220, 220, 120, 255},
		DialogText:              color.White,
		EndCardBG:               color.RGBA{0, 0, 0, 230},
		EndCardText:             color.White,
		CursorColor:             color.RGBA{255, 255, 255, 255},
		HotspotOutline:          color.RGBA{255, 255, 0, 200},
		SceneBackdrop:           color.RGBA{0, 0, 0, 255},

		TopBarBG:             color.RGBA{15, 12, 22, 255},
		TopBarText:           color.RGBA{220, 220, 230, 255},
		TopBarAccent:         color.RGBA{255, 210, 110, 255},
		ChatLogBG:            color.RGBA{12, 10, 18, 230},
		ChatLogPrompt:        color.RGBA{255, 210, 110, 255},
		ChatLogResponse:      color.RGBA{220, 220, 230, 255},
		ChatLogSystem:        color.RGBA{160, 160, 180, 255},
		CharacterPanelBG:     color.RGBA{10, 8, 14, 170},
		CharacterPanelBorder: color.RGBA{120, 100, 60, 255},
		CharacterPanelTitle:  color.RGBA{255, 210, 110, 255},
	}
}
