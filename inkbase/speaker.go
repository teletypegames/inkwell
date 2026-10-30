package inkbase

import "image/color"

// Delivery is how a speaker reaches the player: over the scene as a speech
// bubble, or through the message log alone.
type Delivery int

// The ways a Speaker can be delivered.
const (
	DeliveryBubble Delivery = iota // speech bubble over the scene, plus the log
	DeliveryLog                    // the log alone
)

// Speaker is a voice with no body: something that talks but never stands in
// a scene.
type Speaker struct {
	Name     string
	Label    Text
	Delivery Delivery
	Color    color.Color
	Lore     Text
}

// Talker is whoever a line can be given to: a *Character or a *Speaker.
// Say takes one, and so does a DialogLine, which is what lets a dialog file
// name the person or the voice itself rather than spell its name.
type Talker interface {
	TalkerName() string
}

// GetName returns the speaker's Name, its key in the registry.
func (s Speaker) GetName() string { return s.Name }

// TypeLabel returns "speaker", the speaker's kind in the known-entity
// key and in log and panic messages.
func (s Speaker) TypeLabel() string { return "speaker" }

// DisplayName is the Label if there is one, the Name otherwise.
func (s Speaker) DisplayName() Text {
	if !s.Label.Empty() {
		return s.Label
	}
	return Plain(s.Name)
}
