package inkbase

import "image/color"

// Delivery is how a speaker reaches the player: over the scene as a speech
// bubble, or through the message log alone.
type Delivery int

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

func (s Speaker) GetName() string   { return s.Name }
func (s Speaker) TypeLabel() string { return "speaker" }

// DisplayName is the Label if there is one, the Name otherwise.
func (s Speaker) DisplayName() Text {
	if !s.Label.Empty() {
		return s.Label
	}
	return Plain(s.Name)
}
