package logging

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var timestampRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} `)

func readLines(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}
	trimmed := strings.TrimRight(string(data), "\n")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		in   string
		want LogLevel
	}{
		{"none", None},
		{"", None},
		{"error", Error},
		{"Error", Error},
		{"  ERROR  ", Error},
		{"info", Info},
		{"Info", Info},
		{"debug", Debug},
		{"DEBUG", Debug},
		{"garbage", None},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if got := ParseLevel(c.in); got != c.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestLogLevelString(t *testing.T) {
	cases := []struct {
		in   LogLevel
		want string
	}{
		{None, "None"},
		{Error, "Error"},
		{Info, "Info"},
		{Debug, "Debug"},
		{LogLevel(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.in.String(); got != c.want {
			t.Errorf("LogLevel(%d).String() = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSetLevelGetLevelEnabled(t *testing.T) {
	t.Cleanup(func() { SetLevel(None) })

	SetLevel(None)
	if GetLevel() != None {
		t.Fatalf("GetLevel() = %v, want None", GetLevel())
	}
	for _, l := range []LogLevel{Debug, Info, Error} {
		if Enabled(l) {
			t.Errorf("Enabled(%v) = true while global level is None, want false", l)
		}
	}

	SetLevel(Info)
	if !Enabled(Info) {
		t.Error("Enabled(Info) = false while global level is Info, want true")
	}
	if !Enabled(Error) {
		t.Error("Enabled(Error) = false while global level is Info, want true")
	}
	if Enabled(Debug) {
		t.Error("Enabled(Debug) = true while global level is Info, want false")
	}

	SetLevel(Debug)
	for _, l := range []LogLevel{Debug, Info, Error} {
		if !Enabled(l) {
			t.Errorf("Enabled(%v) = false while global level is Debug, want true", l)
		}
	}
}

// TestInitAndEmit exercises Init, the severity-gated emitters, and Userf.
//
// logging keeps its configuration (EntryPrefix, UserHeader, EffectiveUser)
// behind a sync.Once, so only the *first* Init call in the whole test binary
// can set them. This test must therefore be the first one to call Init, and
// every other Init-dependent test must build on top of what it establishes.
func TestInitAndEmit(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "first.log")

	err := Init(logPath, Debug, LogInitOptions{
		EntryPrefix:        "myapp",
		UserHeader:         "[CUSTOM]",
		DisplayCurrentUser: true,
		DisplayExecName:    true,
		DisplayPID:         true,
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	t.Cleanup(Close)

	if GetLevel() != Debug {
		t.Fatalf("GetLevel() after Init(..., Debug, ...) = %v, want Debug", GetLevel())
	}

	wantUser := ""
	if cu, err := user.Current(); err == nil {
		wantUser = cu.Username
	}
	wantExec := filepath.Base(os.Args[0])
	wantPID := strconv.Itoa(os.Getpid())

	Debugf("debug value: %s", "x")
	Infof("info message")
	Errorf("error message")
	Userf("user message", "")

	lines := readLines(t, logPath)
	if len(lines) != 4 {
		t.Fatalf("got %d log lines, want 4: %#v", len(lines), lines)
	}

	checks := []struct {
		line       string
		wantHeader string
		wantMsg    string
	}{
		{lines[0], "[DEBUG]", "debug value: x"},
		{lines[1], "[INFO]", "info message"},
		{lines[2], "[ERROR]", "error message"},
		{lines[3], "[CUSTOM]", "user message"},
	}

	for _, c := range checks {
		if !timestampRE.MatchString(c.line) {
			t.Errorf("line %q does not start with a timestamp", c.line)
		}
		if !strings.Contains(c.line, c.wantHeader) {
			t.Errorf("line %q does not contain header %q", c.line, c.wantHeader)
		}
		if !strings.Contains(c.line, "<myapp>") {
			t.Errorf("line %q does not contain entry prefix <myapp>", c.line)
		}
		if wantUser != "" && !strings.Contains(c.line, "("+wantUser+")") {
			t.Errorf("line %q does not contain effective user (%s)", c.line, wantUser)
		}
		if !strings.Contains(c.line, wantExec) {
			t.Errorf("line %q does not contain exec name %q", c.line, wantExec)
		}
		if !strings.Contains(c.line, "PID "+wantPID) {
			t.Errorf("line %q does not contain PID %s", c.line, wantPID)
		}
		if !strings.Contains(c.line, c.wantMsg) {
			t.Errorf("line %q does not contain message %q", c.line, c.wantMsg)
		}
		if !strings.HasSuffix(c.line, ">>") {
			t.Errorf("line %q does not end with '>>'", c.line)
		}
	}

	// Gating: dropping the threshold below Debug must suppress further Debugf calls.
	SetLevel(Error)
	Debugf("this should not be written")
	if got := readLines(t, logPath); len(got) != 4 {
		t.Fatalf("after SetLevel(Error), Debugf() added a line: got %d lines, want 4: %#v", len(got), got)
	}
}

// TestInitRotation depends on TestInitAndEmit having already run (and thus
// consumed the package's one-shot sync.Once) in the same test binary.
func TestInitRotation(t *testing.T) {
	dir := t.TempDir()
	secondPath := filepath.Join(dir, "second.log")

	err := Init(secondPath, Info, LogInitOptions{
		EntryPrefix:        "ignored-because-once-already-fired",
		DisplayCurrentUser: true,
		DisplayExecName:    true,
		DisplayPID:         true,
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	t.Cleanup(Close)

	if GetLevel() != Info {
		t.Fatalf("GetLevel() after rotation Init(..., Info, ...) = %v, want Info", GetLevel())
	}

	Infof("second file message")

	lines := readLines(t, secondPath)
	if len(lines) != 1 {
		t.Fatalf("got %d lines in rotated log, want 1: %#v", len(lines), lines)
	}
	// EntryPrefix is only applied on the very first Init() call in the
	// process (it's set inside a sync.Once); later Init() calls cannot
	// change it, so the original "myapp" prefix must still be in effect.
	if !strings.Contains(lines[0], "<myapp>") {
		t.Errorf("rotated log line %q should still carry the original <myapp> prefix", lines[0])
	}
	if strings.Contains(lines[0], "ignored-because-once-already-fired") {
		t.Errorf("rotated log line %q unexpectedly picked up the new EntryPrefix", lines[0])
	}
}
