package inkwell

import (
	"image/color"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

// Speaker is a voice with no body: something that talks but never stands in
// a scene — a tape in a slot, a room-scale system, a narrator. Characters
// speak too, through the same fields, so anything that takes a Talker takes
// either.
type Speaker struct {
	Name string
	// Label is what a log or a dialog box calls the voice; empty falls back
	// to Name.
	Label    Text
	Delivery Delivery
	Color    color.Color
	// Lore is what the player's record says about the voice once it is
	// known — see Character.Lore.
	Lore Text
}

func (s Speaker) GetName() string   { return s.Name }
func (s Speaker) TypeLabel() string { return "speaker" }

func (s Speaker) DisplayName() Text {
	if !s.Label.Empty() {
		return s.Label
	}
	return plain(s.Name)
}

type Talker = inkbase.Talker

func (s *Speaker) TalkerName() string   { return s.Name }
func (c *Character) TalkerName() string { return c.Name }

// Say and Reply on the talker itself read the way a script does — who, then
// what: paul.Say(line). They are the package-level Say and Reply with the
// speaker filled in.
func (s *Speaker) Say(text Text) Action                          { return Say(s, text) }
func (c *Character) Say(text Text) Action                        { return Say(c, text) }
func (s *Speaker) Reply(answer func(state *State) Text) Action   { return Reply(s, answer) }
func (c *Character) Reply(answer func(state *State) Text) Action { return Reply(c, answer) }

func talkerName(t Talker) string {
	if t == nil {
		return ""
	}
	return t.TalkerName()
}

func (c Character) speaker() Speaker {
	return Speaker{
		Name:     c.Name,
		Label:    c.Label,
		Delivery: c.Delivery,
		Color:    c.SpeechColor,
	}
}

// noteSpoken records that the voice behind a name has been heard, which is
// how a Character or a Speaker becomes Known without a script saying so.
func (g *Game) noteSpoken(name string) {
	if c, ok := g.CharacterManager.Get(name); ok {
		g.State.NoteKnown(c)
		return
	}
	if s, ok := g.SpeakerManager.Get(name); ok {
		g.State.NoteKnown(s)
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

func (g *Game) DeliveryOf(name string) Delivery {
	if s, ok := g.Speaker(name); ok {
		return s.Delivery
	}
	return DeliveryBubble
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
