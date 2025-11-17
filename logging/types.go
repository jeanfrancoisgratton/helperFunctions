// Package logging
// logging/types.go
// helperFunctions - logging subpackage
// Refactor: Aug 13, 2025 (USER outside severity ladder) + ParseLevel()
// File: logging/types.go
//
// - Fixed levels owned by the package: None, Error, Info, Debug.
// - USER messages are *outside* the ladder: they log whenever level != None.
// - Format: TIMESTAMP [HEADER] MESSAGE:
// - HEADER = [DEBUG] | [INFO] | [ERROR] | [<user header>]

package logging

import (
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	logger            *log.Logger
	logFile           *os.File
	initOnce          sync.Once
	globalLevel       atomic.Int32 // default 0 (None)
	defaultUserHeader atomic.Value // holds string, defaults to "[USER]"
	displayPID        bool
	displayExecName   bool
	LogEntryPrefix    atomic.Value
	EffectiveUser     atomic.Value
)

// LogInitOptions : an easy way to future-proof the logging facilities
// This struct will be passed as a param to Init(), so we can extend it at will without breaking the package
type LogInitOptions struct {
	EntryPrefix        string // Logfile entry prefix (after the timestamp and loglevel), default ""
	UserHeader         string // Header to use when loglevel set to USER, default ""
	DisplayCurrentUser bool   // Username executing the calling process, default false
	DisplayExecName    bool   // Display the calling executable name, default false
	DisplayPID         bool   // Display the calling ProcessID (PID), default false
}

// LogLevel ordering: None < Error < Info < Debug
type LogLevel int32

const (
	None  LogLevel = iota // no logging at all
	Error                 // only Errorf
	Info                  // Infof + Errorf
	Debug                 // Debugf + Infof + Errorf
)

func (ll LogLevel) String() string {
	switch ll {
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
