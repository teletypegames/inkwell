package inkwell

import "image/color"

type Character struct {
	Name string
	// Label is what the UI shows instead of Name — a code name, a title,
	// a shout. Empty means the Name is presentable as it stands.
	Label       Text
	Sprite      string // Asset.Name (placeholder if missing)
	Animations  map[string]AnimationClip
	Speed       float64
	SpeechColor color.Color
	// Voice decides where Say puts this character's lines.
	Voice Voice
	// Start is where the character stands in a scene that lists no actors
	// of its own, and Game.Player names the one placed there.
	Start Point
	// Size hints used when the sprite is a placeholder rectangle.
	W, H float64
}

// Voice is how a character reaches the player: over the scene as a speech
// bubble, or through the message log alone — a radio, a tape, a narrator
// that has no body to speak from.
type Voice int

const (
	VoiceBubble Voice = iota // speech bubble over the scene, plus the log
	VoiceLog                 // the log alone
)

func (c Character) GetName() string { return c.Name }

// DisplayName is the Label if there is one, the Name otherwise.
func (c Character) DisplayName() Text {
	if !c.Label.Empty() {
		return c.Label
	}
	return plain(c.Name)
}
func (c Character) TypeLabel() string { return "character" }
