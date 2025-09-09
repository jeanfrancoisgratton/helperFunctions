// logging/logger.go
// helperFunctions - logging subpackage
// Refactor: Aug 13, 2025 (USER outside severity ladder)
// File: logging/logger.go
//
// Exported helpers used by callers.
// Format: TIMESTAMP [HEADER] MESSAGE:
package logging

import (
	"fmt"
	"os"
	"path/filepath"
)

// Debugf emits when level >= Debug.
func Debugf(msg string, args ...any) {
	if !Enabled(Debug) {
		return
	}
	emit("[DEBUG]", formatMessage(msg, args...))
}

// Infof emits when level >= Info.
func Infof(msg string, args ...any) {
	if !Enabled(Info) {
		return
	}
	emit("[INFO]", formatMessage(msg, args...))
}

// Errorf emits when level >= Error.
func Errorf(msg string, args ...any) {
	if !Enabled(Error) {
		return
	}
	emit("[ERROR]", formatMessage(msg, args...))
}

// Userf emits regardless of the severity ladder; it only respects None.
// If header is empty, uses the default set in Init (default "[USER]").
// Header is bracketed if missing brackets.
func Userf(msg string, header string, args ...any) {
	if !userEnabled() {
		return
	}
	h := header
	if h == "" {
		h = currentUserHeader()
	}
	// ensure bracketed
	if h[0] != '[' || h[len(h)-1] != ']' {
		h = "[" + h + "]"
	}
	emit(h, formatMessage(msg, args...))
}

func formatMessage(msg string, args ...any) string {
	procInfo := ""
	if DisplayExecName {
		procInfo = filepath.Base(os.Args[0])
		if DisplayPID {
			procInfo = fmt.Sprintf("%s (PID %d)", procInfo, os.Getpid())
		}
	} else {
		if DisplayPID {
			procInfo = fmt.Sprintf("PID %d", os.Getpid())
		}
	}
	if procInfo != "" {
		msg = fmt.Sprintf("%s %s >", procInfo, msg)
	}

	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}
