package inkwell

import (
	"image/color"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

type Delivery = inkbase.Delivery

const (
	DeliveryBubble = inkbase.DeliveryBubble
	DeliveryLog    = inkbase.DeliveryLog
)

type Character struct {
	Name string
	// Label is what the UI shows instead of Name — a code name, a title,
	// a shout. Empty means the Name is presentable as it stands.
	Label Text
	// Lore is what the player's record says about the character once they
	// are known — a bio for a dossier or a codex screen. Nil keeps the
	// character out of such a screen.
	Lore Text
	// Poses are the sprites the character is drawn with, one per Pose the
	// artwork covers. Nil, or a Pose with no entry, falls back to Idle; no
	// Idle draws the placeholder figure sized W by H.
	Poses       map[Pose]*Sprite
	Speed       float64
	SpeechColor color.Color
	// Delivery decides where Say puts this character's lines.
	Delivery Delivery
	// Start is where the character stands in a scene that lists no actors
	// of its own, and Game.Player names the one placed there.
	Start Point
	// Size hints used when the sprite is a placeholder rectangle.
	W, H float64
}

func (c Character) GetName() string { return c.Name }

// DisplayName is the Label if there is one, the Name otherwise.
func (c Character) DisplayName() Text {
	if !c.Label.Empty() {
		return c.Label
	}
	return plain(c.Name)
}
func (c Character) TypeLabel() string { return "character" }
