package theme

import "image/color"

// paper-notebook is the light theme: ink on cream, for a game that reads
// as something written down rather than something lit up.

func paperNotebookTheme() *Theme {
	cream := color.RGBA{245, 235, 210, 255}
	ink := color.RGBA{50, 40, 30, 255}
	return &Theme{
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
