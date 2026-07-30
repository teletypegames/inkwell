package inkwell

import "image/color"

type Character struct {
	Name        string
	Sprite      string // Asset.Name (placeholder if missing)
	Animations  map[string]AnimationClip
	Speed       float64
	SpeechColor color.Color
	Start       Point
	// Size hints used when the sprite is a placeholder rectangle.
	W, H float64
}

func (c Character) GetName() string   { return c.Name }
func (c Character) TypeLabel() string { return "character" }
