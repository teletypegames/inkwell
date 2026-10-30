package inkbase

// LogKind classifies a chat-log entry — UI widgets style them differently.
type LogKind int

// The kinds of chat-log entry.
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
