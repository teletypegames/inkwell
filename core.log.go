package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

type LogKind = inkbase.LogKind
type LogMessage = inkbase.LogMessage

const (
	LogAction   = inkbase.LogAction
	LogResponse = inkbase.LogResponse
	LogSystem   = inkbase.LogSystem
)

func (g *Game) LogAction(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogAction})
}

func (g *Game) LogResponse(speaker, text string) {
	g.appendLog(LogMessage{Speaker: speaker, Text: text, Kind: LogResponse})
}

func (g *Game) LogSystem(text string) {
	g.appendLog(LogMessage{Text: text, Kind: LogSystem})
}

func (g *Game) appendLog(m LogMessage) {
	g.messages = append(g.messages, m)
	if g.MaxLogLines > 0 && len(g.messages) > g.MaxLogLines {
		g.messages = g.messages[len(g.messages)-g.MaxLogLines:]
	}
}

func (g *Game) Messages() []LogMessage { return g.messages }
