package inkwell

import "image/color"

// RegisterPresetThemes installs the four built-in themes. NewGame calls
// this automatically and selects "classic-scumm".
func RegisterPresetThemes(g *Game) {
	g.ThemeManager.Register(classicScummTheme())
	g.ThemeManager.Register(sierraCoinTheme())
	g.ThemeManager.Register(paperNotebookTheme())
	g.ThemeManager.Register(terminalGreenTheme())
}

func classicScummTheme() Theme {
	return Theme{
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

func paperNotebookTheme() Theme {
	cream := color.RGBA{245, 235, 210, 255}
	ink := color.RGBA{50, 40, 30, 255}
	return Theme{
		Name:                    "paper-notebook",
		PanelBG:                 cream,
		StatusText:              ink,
		FlashText:               color.RGBA{180, 50, 50, 255},
		VerbButtonBG:            color.RGBA{225, 215, 190, 255},
		VerbButtonSelectedBG:    color.RGBA{180, 150, 80, 255},
		VerbButtonText:          ink,
		InventorySlotBG:         color.RGBA{230, 220, 195, 255},
		InventorySlotSelectedBG: color.RGBA{180, 150, 80, 255},
		SpeechBubbleBG:          color.RGBA{255, 255, 245, 235},
		SpeechDefaultText:       ink,
		DialogBG:                color.RGBA{250, 240, 215, 250},
		DialogBorder:            ink,
		DialogChoiceBG:          color.RGBA{235, 225, 200, 255},
		DialogChoiceHover:       color.RGBA{220, 200, 160, 255},
		DialogSpeaker:           color.RGBA{120, 70, 30, 255},
		DialogText:              ink,
		EndCardBG:               color.RGBA{245, 235, 210, 240},
		EndCardText:             ink,
		CursorColor:             ink,
		HotspotOutline:          color.RGBA{120, 70, 30, 200},
		SceneBackdrop:           cream,

		TopBarBG:             color.RGBA{220, 205, 170, 255},
		TopBarText:           ink,
		TopBarAccent:         color.RGBA{120, 70, 30, 255},
		ChatLogBG:            color.RGBA{250, 240, 215, 230},
		ChatLogPrompt:        color.RGBA{120, 70, 30, 255},
		ChatLogResponse:      ink,
		ChatLogSystem:        color.RGBA{120, 110, 90, 255},
		CharacterPanelBG:     color.RGBA{250, 240, 215, 200},
		CharacterPanelBorder: ink,
		CharacterPanelTitle:  color.RGBA{120, 70, 30, 255},
	}
}

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
		DialogBG:                color.RGBA{0, 16, 0, 240},
		DialogBorder:            green,
		DialogChoiceBG:          color.RGBA{12, 22, 12, 255},
		DialogChoiceHover:       dim,
		DialogSpeaker:           color.RGBA{180, 255, 180, 255},
		DialogText:              green,
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
