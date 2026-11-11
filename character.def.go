package inkwell

import (
	"image/color"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

// Delivery is where a Character's or a Speaker's lines reach the player;
// see [inkbase.Delivery].
type Delivery = inkbase.Delivery

// The ways a line can be delivered.
const (
	DeliveryBubble = inkbase.DeliveryBubble // a speech bubble over the scene, and the log
	DeliveryLog    = inkbase.DeliveryLog    // the message log alone
)

// Character is someone with a body: drawn in a scene with its Poses,
// walking through the walkboxes at Speed, and speaking like a Speaker.
// Game.Player names the one the player steers.
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
	Poses map[Pose]*Sprite
	// Speed is how fast the character walks, in pixels per second.
	Speed float64
	// SpeechColor is the colour of the character's lines; nil, or a
	// fully transparent color.RGBA, leaves the theme's.
	SpeechColor color.Color
	// Delivery decides where Say puts this character's lines.
	Delivery Delivery
	// Start is where the character stands in a scene that lists no actors
	// of its own, and Game.Player names the one placed there.
	Start Point
	// Size hints used when the sprite is a placeholder rectangle.
	W, H float64
}

// GetName returns the character's Name, its key in the CharacterManager.
func (c Character) GetName() string { return c.Name }

// DisplayName is the Label if there is one, the Name otherwise.
func (c Character) DisplayName() Text {
	if !c.Label.Empty() {
		return c.Label
	}
	return plain(c.Name)
}

// TypeLabel returns "character", the kind a Character is registered as.
func (c Character) TypeLabel() string { return "character" }

// Hide is the Action that takes the character out of the picture until
// Show; see [Hide].
func (c *Character) Hide() Action { return Hide(c) }

// Show is the Action that brings the character back after Hide; see [Show].
func (c *Character) Show() Action { return Show(c) }
