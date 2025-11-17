// logging/message_ops.go
// helperFunctions - logging subpackage
// Refactor: Aug 13, 2025 (USER outside severity ladder)
// File: logging/message_ops.go
//
// Exported helpers used by callers.
// Format: TIMESTAMP [HEADER] MESSAGE:

package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// formatMessage : this is where we format the log entry line
// The format is <PREFIX> (USERNAME) PROCNAME (PID) MESSAGE ARGS
// PREFIX, USERNAME, PROCNAME and PID are optional
func formatMessage(msg string, args ...any) string {
	procInfo := ""
	prefixline := ""
	eUser := ""

	if prefixline, _ = LogEntryPrefix.Load().(string); prefixline != "" {
		prefixline = "<" + prefixline + ">"
	}

	if eUser, _ = EffectiveUser.Load().(string); eUser != "" {
		eUser = "(" + eUser + ")"
	}

	if displayExecName {
		procInfo = filepath.Base(os.Args[0])
		if displayPID {
			procInfo = fmt.Sprintf("%s (PID %d)", procInfo, os.Getpid())
		}
	} else {
		if displayPID {
			procInfo = fmt.Sprintf("PID %d", os.Getpid())
		}
	}
	if procInfo != "" {
		msg = fmt.Sprintf("%s %s %s %s >>", prefixline, eUser, procInfo, msg)
	}

	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

// userEnabled is the gating rule for Userf:
// log if global level != None (outside severity ladder).
func userEnabled() bool {
	return GetLevel() != None
}

// emit writes a single, already-gated line.
// header must include brackets (e.g., "[INFO]"); message is the final text.
func emit(header string, message string) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	ts := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("%s %s %s", ts, header, message)
}

func currentUserHeader() string {
	v := defaultUserHeader.Load()
	if v == nil {
		return "[USER]"
	}
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	return "[USER]"
}
