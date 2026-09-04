package inkwell

import "image/color"

// Speaker is a voice with no body: something that talks but never stands in
// a scene — a tape in a slot, a room-scale system, a narrator. Characters
// speak too, through the same fields, so anything that takes a Talker takes
// either.
type Speaker struct {
	Name string
	// Label is what a log or a dialog box calls the voice; empty falls back
	// to Name.
	Label Text
	Voice Voice
	Color color.Color
}

func (s Speaker) GetName() string   { return s.Name }
func (s Speaker) TypeLabel() string { return "speaker" }

func (s Speaker) DisplayName() Text {
	if !s.Label.Empty() {
		return s.Label
	}
	return plain(s.Name)
}

// Talker is whoever a line can be given to: a *Character or a *Speaker. Say
// takes one, and so does a DialogLine, which is what lets a dialog file
// name the person or the voice itself rather than spell its name.
type Talker interface {
	talkerName() string
}

func (s *Speaker) talkerName() string   { return s.Name }
func (c *Character) talkerName() string { return c.Name }

func talkerName(t Talker) string {
	if t == nil {
		return ""
	}
	return t.talkerName()
}

func (c Character) speaker() Speaker {
	return Speaker{
		Name:  c.Name,
		Label: c.Label,
		Voice: c.Voice,
		Color: c.SpeechColor,
	}
}

// Speaker resolves a name that has spoken — a log entry's, a bubble's —
// to the voice behind it: a character first, then a registered Speaker.
func (g *Game) Speaker(name string) (Speaker, bool) {
	if c, ok := g.CharacterManager.Get(name); ok {
		return c.speaker(), true
	}
	if s, ok := g.SpeakerManager.Get(name); ok {
		return *s, true
	}
	return Speaker{}, false
}

func (g *Game) VoiceOf(name string) Voice {
	if s, ok := g.Speaker(name); ok {
		return s.Voice
	}
	return VoiceBubble
}

func (g *Game) SpeakerLabel(name string) Text {
	if s, ok := g.Speaker(name); ok {
		return s.DisplayName()
	}
	return plain(name)
}

func (g *Game) SpeechColor(name string, fallback color.Color) color.Color {
	if s, ok := g.Speaker(name); ok && s.Color != nil {
		if rgba, ok := s.Color.(color.RGBA); ok && rgba.A == 0 {
			return fallback
		}
		return s.Color
	}
	return fallback
}
