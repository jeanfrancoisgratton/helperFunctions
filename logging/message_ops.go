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

func formatMessage(msg string, args ...any) string {
	procInfo := ""

	if prefixline, _ := LogEntryPrefix.Load().(string); prefixline != "" {
		msg = prefixline + " " + msg
	}

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
	ts := time.Now().Format(timeLayout)
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
