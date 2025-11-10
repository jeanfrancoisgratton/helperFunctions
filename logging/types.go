// Package logging
// logging/types.go
// helperFunctions - logging subpackage
// Refactor: Aug 13, 2025 (USER outside severity ladder) + ParseLevel()
// File: logging/types.go
//
// - No dependency on customError: callers pass plain strings.
// - Fixed levels owned by the package: None, Error, Info, Debug.
// - USER messages are *outside* the ladder: they log whenever level != None.
// - Format: TIMESTAMP [HEADER] MESSAGE:
// - HEADER = [DEBUG] | [INFO] | [ERROR] | [<user header>]
// Optionally, if the "LinePrefix" variable is set, the prefix will be pre-pended right before the message

package logging

import (
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// Timestamp layout is used by emit() in state.go
const timeLayout = "2006-01-02 15:04:05"

// LogLevel ordering: None < Error < Info < Debug
type LogLevel int32

const (
	None  LogLevel = iota // no logging at all
	Error                 // only Errorf
	Info                  // Infof + Errorf
	Debug                 // Debugf + Infof + Errorf
)

func (l LogLevel) String() string {
	switch l {
	case None:
		return "None"
	case Error:
		return "Error"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	default:
		return "Unknown"
	}
}

// ParseLevel converts a CLI/user string into a LogLevel.
// Unrecognized values default to None (strict).
func ParseLevel(s string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none", "":
		return None
	case "error":
		return Error
	case "info":
		return Info
	case "debug":
		return Debug
	default:
		return None
	}
}

var (
	logger            *log.Logger
	logFile           *os.File
	initOnce          sync.Once
	globalLevel       atomic.Int32 // default 0 (None)
	defaultUserHeader atomic.Value // holds string, defaults to "[USER]"
	DisplayPID        bool
	DisplayExecName   bool
	LogEntryPrefix    atomic.Value
	EffectiveUser     atomic.Value
)
