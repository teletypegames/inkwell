package inkwell

import (
	"fmt"
	"os"
)

var DebugLog = false

func logf(format string, args ...any) {
	if !DebugLog {
		return
	}
	fmt.Fprintf(os.Stderr, "[inkwell] "+format+"\n", args...)
}
