package lang

import (
	"runtime"
	"strings"
)

// relevantCaller searches the call stack for the first function outside of pkg
// The purpose of this function is to provide more helpful error messages.
// how to use? example:
// caller := relevantCaller("net/http.")
// logx.Info("superfluous WriteHeader call from %s (%s:%d)", caller.Function, path.Base(caller.File), caller.Line)
func RelevantCaller(outsidePkg string) runtime.Frame {
	pc := make([]uintptr, 16)
	n := runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	var frame runtime.Frame
	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, outsidePkg) {
			return frame
		}
		if !more {
			break
		}
	}
	return frame
}
