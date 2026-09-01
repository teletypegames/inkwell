package inkwell

import "image/color"

// Speaker is a voice with no body: a radio, a cassette, a narrator —
// anything Say can put a line in the mouth of that never stands in a
// scene. A Character is a speaker with a body, and the engine resolves a
// speaker's name through both: the characters first, then the speakers
// registered here. A name found in neither speaks in a bubble under its
// own name, which is what a bare Say always did.
type Speaker struct {
	Name string
	// Label is what the UI shows instead of Name. Empty means the Name is
	// presentable as it stands.
	Label Text
	// Voice decides where Say puts this speaker's lines.
	Voice Voice
	// Color tints the speaker's lines in the log and the bubble. Nil takes
	// the theme's default.
	Color color.Color
}

func (s Speaker) GetName() string { return s.Name }

// DisplayName is the Label if there is one, the Name otherwise.
func (s Speaker) DisplayName() Text {
	if !s.Label.Empty() {
		return s.Label
	}
	return plain(s.Name)
}

// speaker is a Character's voice, or a registered Speaker's.
func (c Character) speaker() Speaker {
	return Speaker{
		Name:  c.Name,
		Label: c.Label,
		Voice: c.Voice,
		Color: c.SpeechColor,
	}
}

// Speaker resolves a name to whoever speaks under it: a character, or a
// speaker registered without a body. The second result is false when the
// name is neither.
func (g *Game) Speaker(name string) (Speaker, bool) {
	if c, ok := g.CharacterManager.Get(name); ok {
		return c.speaker(), true
	}
	if s, ok := g.SpeakerManager.Get(name); ok {
		return s, true
	}
	return Speaker{}, false
}

// VoiceOf reports how a speaker reaches the player. An unregistered name
// speaks in a bubble, which is what a bare Say always did.
func (g *Game) VoiceOf(name string) Voice {
	if s, ok := g.Speaker(name); ok {
		return s.Voice
	}
	return VoiceBubble
}

// SpeakerLabel is what the UI should print for a speaker: the Label if the
// entity carries one, the name otherwise.
func (g *Game) SpeakerLabel(name string) Text {
	if s, ok := g.Speaker(name); ok {
		return s.DisplayName()
	}
	return plain(name)
}

// SpeechColor is the color a speaker's lines are drawn in, or fallback when
// the speaker has none of its own.
func (g *Game) SpeechColor(name string, fallback color.Color) color.Color {
	if s, ok := g.Speaker(name); ok && s.Color != nil {
		if rgba, ok := s.Color.(color.RGBA); ok && rgba.A == 0 {
			return fallback
		}
		return s.Color
	}
	return fallback
}
