// logging/state.go
// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/08/13 22:51
// Original filename: logging/state.go
//

package logging

import (
	"log"
	"os"
	"os/user"
)

// Init :
// The function sets output, global threshold, default user header. log entry prefix, etc.
// path "-" or "" -> stdout; otherwise file (0640) is opened/created.
// Re-invocation rotates to the new target.

// To initialize the log facilities, you set the following variables
// path							:-> the path to the logfile
// level						:-> the loglevel (none, debug, info, error, user)

// The other parameters are set with the LogInitOptions structure, which initializes the following members:
// EntryPrefix					:-> a prefix to add before every log entry
// UserHeader					:-> a user-defined prefix to add if the loglevel is set to USER
// DisplayCurrentUser (boolean)	:-> the user currently running the tool
// DisplayExecName (boolean)	:-> display the executable name in the log entry
// DisplayPID (boolean) 		:-> display the process PID

// displayExecName and displayPID might not be relevant for app-specific logfiles. In other words:
// If this package is called to log into, say, /var/log/myapp.log, we could safely assume that displayExecName here
// Would be set to "myapp", not really useful, right ?

func Init(path string, level LogLevel, logOptions LogInitOptions) error {
	var err error
	initOnce.Do(func() {
		globalLevel.Store(int32(None))
		defaultUserHeader.Store(logOptions.UserHeader)
		LogEntryPrefix.Store(logOptions.EntryPrefix)
		if logOptions.DisplayCurrentUser {
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

	displayPID = logOptions.DisplayPID
	displayExecName = logOptions.DisplayExecName

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
