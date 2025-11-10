// logging/state.go
// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/08/13 22:51
// Original filename: logging/state.go
//
// Init/Close and global state helpers.
package logging

import (
	"log"
	"os"
	"os/user"
)

// Init : *DEPRECATION NOTICE*
// This function signature will change in v4.
// For backward compatibility, Init() right now uses the legacy signature and calls the next-gen
// Version with empty/null/false values for the extended parameters.
// If you wish to use the extended parameters, use InitExtended()
func Init(path string, level LogLevel, userHeader string, displayExecName, displayPID bool) error {
	return InitExtended(path, level, "", userHeader, false, displayExecName, displayPID)
}

// InitWithPrefix :
// In v4, this will be the new implementation of Init, albeit with most of the options passed in a struct
// The function sets output, global threshold, default user header. log entry prefix, etc.
// path "-" or "" -> stdout; otherwise file (0640) is opened/created.
// Re-invocation rotates to the new target.

// To initialize the log facilities, you set the following variables
// path							:-> the path to the logfile
// level						:-> the loglevel (none, debug, info, error, user)
// entryPrefix					:-> a prefix to add before every log entry
// userHeader					:-> a user-defined prefix to add if the loglevel is set to USER
// displayCurrentUser (boolean)	:-> the user currently running the tool
// displayExecName (boolean)	:-> display the executable name in the log entry
// displayPID (boolean) 		:-> display the process PID

// displayExecName and displayPID might not be relevant for app-specific logfiles. In other words:
// If this package is called to log into, say, /var/log/myapp.log, we could safely assume that displayExecName here
// Would be set to "myapp", not really useful, right ?

func InitExtended(path string, level LogLevel, entryPrefix string, userHeader string,
	displayCurrentUser bool, displayExecName, displayPID bool) error {
	var err error
	initOnce.Do(func() {
		globalLevel.Store(int32(None))
		defaultUserHeader.Store("[USER]")
		LogEntryPrefix.Store(entryPrefix)
		if displayCurrentUser {
			cUsr, err := user.Current()
			if err != nil {
				EffectiveUser.Store("")
			} else {
				EffectiveUser.Store(cUsr.Username)
			}
		}
	})

	// Close previously opened file if any (except stdout)
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}

	var out *os.File
	if path == "-" || path == "" {
		out = os.Stdout
	} else {
		out, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		logFile = out
	}

	logger = log.New(out, "", 0) // we format lines ourselves
	SetLevel(level)

	if userHeader == "" {
		defaultUserHeader.Store("[USER]")
	} else {
		defaultUserHeader.Store(userHeader)
	}

	DisplayPID = displayPID
	DisplayExecName = displayExecName

	return nil
}

// Close closes the underlying file if we opened one. Safe to call multiple times.
func Close() {
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}

// SetLevel sets the global threshold.
func SetLevel(l LogLevel) { globalLevel.Store(int32(l)) }

// GetLevel returns the current global threshold.
func GetLevel() LogLevel { return LogLevel(globalLevel.Load()) }

// Enabled reports whether messages at level l should be emitted under the current threshold.
// None disables everything; otherwise we emit when l <= current threshold.
func Enabled(l LogLevel) bool {
	cur := GetLevel()
	if cur == None {
		return false
	}
	return l <= cur
}
