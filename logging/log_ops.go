// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/11/09 10:58
// Original filename: logging/log_ops.go

package logging

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
