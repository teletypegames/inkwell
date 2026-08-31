package inkwell

import "image/color"

// sierra-coin is the same HUD in violet and copper.

func sierraCoinTheme() Theme {
	return Theme{
		Name:                    "sierra-coin",
		PanelBG:                 color.RGBA{15, 10, 25, 255},
		StatusText:              color.RGBA{220, 220, 240, 255},
		FlashText:               color.RGBA{255, 200, 100, 255},
		VerbButtonBG:            color.RGBA{40, 30, 70, 255},
		VerbButtonSelectedBG:    color.RGBA{170, 110, 60, 255},
		VerbButtonText:          color.RGBA{240, 240, 255, 255},
		InventorySlotBG:         color.RGBA{30, 20, 50, 255},
		InventorySlotSelectedBG: color.RGBA{170, 110, 60, 255},
		SpeechBubbleBG:          color.RGBA{10, 10, 20, 220},
		SpeechDefaultText:       color.RGBA{240, 240, 255, 255},
		DialogBG:                color.RGBA{10, 5, 20, 245},
		DialogBorder:            color.RGBA{120, 60, 180, 255},
		DialogChoiceBG:          color.RGBA{25, 15, 45, 255},
		DialogChoiceHover:       color.RGBA{60, 40, 100, 255},
		DialogSpeaker:           color.RGBA{200, 180, 240, 255},
		DialogText:              color.RGBA{240, 240, 255, 255},
		EndCardBG:               color.RGBA{0, 0, 0, 240},
		EndCardText:             color.RGBA{220, 200, 255, 255},
		CursorColor:             color.RGBA{255, 200, 140, 255},
		HotspotOutline:          color.RGBA{200, 140, 255, 180},
		SceneBackdrop:           color.RGBA{0, 0, 0, 255},

		TopBarBG:             color.RGBA{12, 8, 24, 255},
		TopBarText:           color.RGBA{210, 200, 240, 255},
		TopBarAccent:         color.RGBA{255, 180, 110, 255},
		ChatLogBG:            color.RGBA{10, 8, 20, 235},
		ChatLogPrompt:        color.RGBA{255, 180, 110, 255},
		ChatLogResponse:      color.RGBA{210, 200, 240, 255},
		ChatLogSystem:        color.RGBA{150, 130, 190, 255},
		CharacterPanelBG:     color.RGBA{8, 6, 18, 180},
		CharacterPanelBorder: color.RGBA{120, 60, 180, 255},
		CharacterPanelTitle:  color.RGBA{255, 180, 110, 255},
	}
}
