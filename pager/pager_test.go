package pager

import (
	"os"
	"strings"
	"testing"
)

// Page() is heavily coupled to a real interactive terminal (raw mode,
// blocking single-keystroke reads), so it can't be driven end-to-end in an
// automated test without a pty. These tests cover the two paths that are
// safe and deterministic to exercise: the early-return for empty input, and
// the raw-mode failure when stdin isn't a terminal (guaranteed here because
// we point os.Stdin at a pipe, which is never a tty).

func TestPageEmptyInput(t *testing.T) {
	if err := Page(nil, 0, 0); err != nil {
		t.Errorf("Page(nil, 0, 0) error = %v, want nil", err)
	}
	if err := Page([]string{}, 2, 2); err != nil {
		t.Errorf("Page([]string{}, 2, 2) error = %v, want nil", err)
	}
}

func TestPageNonTTYStdin(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	orig := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = orig }()

	err = Page([]string{"line one", "line two"}, 0, 0)
	if err == nil {
		t.Fatal("Page() with non-tty stdin expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "raw mode") {
		t.Errorf("Page() error = %q, want it to mention raw mode", err.Error())
	}
}
