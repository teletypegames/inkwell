package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// LogKind says what a message-log entry is, so a log widget can style it;
// see [inkbase.LogKind].
type LogKind = inkbase.LogKind

// LogMessage is one line of the message log; see [inkbase.LogMessage].
type LogMessage = inkbase.LogMessage

// The kinds of message-log entry; see [inkbase.LogAction].
const (
	LogAction   = inkbase.LogAction   // a player command, "> look at door"
	LogResponse = inkbase.LogResponse // an in-world reply, a line someone said
	LogSystem   = inkbase.LogSystem   // a message from the game itself
)

// LogAction appends a player command to the message log, trimming the
// oldest lines past MaxLogLines.
func (g *Game) LogAction(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogAction})
}

// LogResponse appends a line said by speaker to the message log, trimming
// the oldest lines past MaxLogLines.
func (g *Game) LogResponse(speaker, text string) {
	g.appendLog(LogMessage{Speaker: speaker, Text: text, Kind: LogResponse})
}

// LogSystem appends a message from the game itself to the message log,
// trimming the oldest lines past MaxLogLines.
func (g *Game) LogSystem(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogSystem})
}

func (g *Game) appendLog(m LogMessage) {
	g.messages = append(g.messages, m)
	if g.MaxLogLines > 0 && len(g.messages) > g.MaxLogLines {
		g.messages = g.messages[len(g.messages)-g.MaxLogLines:]
	}
}

// Messages returns the message log, oldest first. The slice is the game's
// own buffer, not a copy, so a caller must not change it.
func (g *Game) Messages() []LogMessage { return g.messages }
