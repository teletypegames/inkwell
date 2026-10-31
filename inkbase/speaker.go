package inkbase

// Delivery is how a speaker reaches the player: over the scene as a speech
// bubble, or through the message log alone.
type Delivery int

// The ways a voice can be delivered.
const (
	DeliveryBubble Delivery = iota // speech bubble over the scene, plus the log
	DeliveryLog                    // the log alone
)

// Talker is whoever a line can be given to — in inkwell a *Character or a
// *Speaker. Say takes one, and so does a DialogLine, which is what lets a
// dialog file name the person or the voice itself rather than spell its
// name.
type Talker interface {
	TalkerName() string
}
