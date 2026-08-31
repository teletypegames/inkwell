package inkwell

// LogKind classifies a chat-log entry — UI widgets style them differently.
type LogKind int

const (
	LogAction   LogKind = iota // player command ("> look at door")
	LogResponse                // in-world reply / Say output
	LogSystem                  // game/system message
)

// LogMessage is a single line in the chat-log buffer.
type LogMessage struct {
	Speaker string // empty for actions / system
	Text    string
	Kind    LogKind
}

// LogAction appends a player-command line (rendered with the prompt color
// by ChatLog widgets).
func (g *Game) LogAction(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogAction})
}

// LogResponse appends a reply line (rendered with the response color).
func (g *Game) LogResponse(speaker, text string) {
	g.appendLog(LogMessage{Speaker: speaker, Text: text, Kind: LogResponse})
}

// LogSystem appends a system/meta line (rendered with the response color).
func (g *Game) LogSystem(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogSystem})
}

func (g *Game) appendLog(m LogMessage) {
	g.messages = append(g.messages, m)
	if g.MaxLogLines > 0 && len(g.messages) > g.MaxLogLines {
		g.messages = g.messages[len(g.messages)-g.MaxLogLines:]
	}
}

// Messages returns the current log buffer (read-only).
func (g *Game) Messages() []LogMessage { return g.messages }
