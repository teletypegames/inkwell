package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase/chatlog"

// SpeakerStyle says how a ChatLog names who is talking; see
// [chatlog.SpeakerStyle].
type SpeakerStyle = chatlog.SpeakerStyle

const (
	// SpeakerInline puts "DEX: " in front of every reply.
	SpeakerInline = chatlog.SpeakerInline
	// SpeakerHeader names only the latest speaker, once, above a rule.
	SpeakerHeader = chatlog.SpeakerHeader
	// SpeakerNone names nobody.
	SpeakerNone = chatlog.SpeakerNone
)

// ChatLog renders the in-game message log (player commands and in-world
// replies) as a scrolling box, in the active theme's chat log colours; see
// [chatlog.ChatLog] and [Game.ChatLogColors].
type ChatLog = chatlog.ChatLog

// ChatLogColors is the palette a [ChatLog] draws with, from the active
// theme: the chat log background, prompt, response and system colours, and
// the character panel border for the frame and the header rule. It
// satisfies [chatlog.Host].
func (g *Game) ChatLogColors() chatlog.Colors {
	th := g.Theme()
	return chatlog.Colors{
		BG:       th.ChatLogBG,
		Border:   th.CharacterPanelBorder,
		Prompt:   th.ChatLogPrompt,
		Response: th.ChatLogResponse,
		System:   th.ChatLogSystem,
	}
}
