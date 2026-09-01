package inkwell

import "image/color"

// terminal-green is a phosphor screen: one green, one dimmer green, and
// black behind them.

func terminalGreenTheme() Theme {
	black := color.RGBA{8, 12, 8, 255}
	green := color.RGBA{120, 220, 120, 255}
	dim := color.RGBA{60, 120, 60, 255}
	return Theme{
		Name:                    "terminal-green",
		PanelBG:                 black,
		StatusText:              green,
		FlashText:               color.RGBA{220, 220, 100, 255},
		VerbButtonBG:            color.RGBA{12, 22, 12, 255},
		VerbButtonSelectedBG:    dim,
		VerbButtonText:          green,
		InventorySlotBG:         color.RGBA{12, 22, 12, 255},
		InventorySlotSelectedBG: dim,
		SpeechBubbleBG:          color.RGBA{0, 16, 0, 220},
		SpeechDefaultText:       green,
		DialogueBG:              color.RGBA{0, 16, 0, 240},
		DialogueBorder:          green,
		DialogueChoiceBG:        color.RGBA{12, 22, 12, 255},
		DialogueChoiceHover:     dim,
		DialogueSpeaker:         color.RGBA{180, 255, 180, 255},
		DialogueText:            green,
		EndCardBG:               color.RGBA{0, 0, 0, 250},
		EndCardText:             green,
		CursorColor:             green,
		HotspotOutline:          dim,
		SceneBackdrop:           black,

		TopBarBG:             color.RGBA{0, 18, 0, 255},
		TopBarText:           green,
		TopBarAccent:         color.RGBA{220, 220, 100, 255},
		ChatLogBG:            color.RGBA{0, 12, 0, 230},
		ChatLogPrompt:        color.RGBA{220, 220, 100, 255},
		ChatLogResponse:      green,
		ChatLogSystem:        dim,
		CharacterPanelBG:     color.RGBA{0, 10, 0, 200},
		CharacterPanelBorder: green,
		CharacterPanelTitle:  color.RGBA{220, 220, 100, 255},
	}
}
