package inkwell

import (
	"fmt"
	"os"
)

// DebugLog turns on the engine's debug trace: queued actions, audio that
// could not play, failed saves, written to standard error with an
// "[inkwell]" prefix. Off by default.
var DebugLog = false

func logf(format string, args ...any) {
	if !DebugLog {
		return
	}
	fmt.Fprintf(os.Stderr, "[inkwell] "+format+"\n", args...)
}
