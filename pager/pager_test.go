package pager

import (
	"os"
	"reflect"
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

func TestWrapLines(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		width int
		want  []string
	}{
		{"no wrap needed", []string{"abc", "de"}, 5, []string{"abc", "de"}},
		{"exact width", []string{"abcde"}, 5, []string{"abcde"}},
		{"wraps once", []string{"abcdef"}, 5, []string{"abcde", "f"}},
		{"wraps twice", []string{"abcdefghij"}, 4, []string{"abcd", "efgh", "ij"}},
		{"empty line preserved", []string{""}, 5, []string{""}},
		{"zero width disables wrapping", []string{"abcdef"}, 0, []string{"abcdef"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLines(tt.lines, tt.width)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapLines(%v, %d) = %v, want %v", tt.lines, tt.width, got, tt.want)
			}
		})
	}
}

func TestLayoutBasic(t *testing.T) {
	lines := []string{"B1", "B2", "c1", "c2", "c3", "c4", "F1"}
	banner, content, footer, pageSize := layout(lines, 2, 1, 80, 10)

	if !reflect.DeepEqual(banner, []string{"B1", "B2"}) {
		t.Errorf("banner = %v", banner)
	}
	if !reflect.DeepEqual(content, []string{"c1", "c2", "c3", "c4"}) {
		t.Errorf("content = %v", content)
	}
	if !reflect.DeepEqual(footer, []string{"F1"}) {
		t.Errorf("footer = %v", footer)
	}
	// termHeight(10) - banner(2) - footer(1) - status bar(1) = 6
	if pageSize != 6 {
		t.Errorf("pageSize = %d, want 6", pageSize)
	}
}

func TestLayoutDropsStickyRegionsWhenCrowded(t *testing.T) {
	lines := []string{"B1", "B2", "c1", "c2", "F1"}
	banner, content, footer, pageSize := layout(lines, 2, 1, 80, 3)

	if banner != nil {
		t.Errorf("banner = %v, want nil (dropped)", banner)
	}
	if footer != nil {
		t.Errorf("footer = %v, want nil (dropped)", footer)
	}
	if !reflect.DeepEqual(content, lines) {
		t.Errorf("content = %v, want all lines %v", content, lines)
	}
	if pageSize != 2 {
		t.Errorf("pageSize = %d, want 2", pageSize)
	}
}

func TestLayoutWrapsWideContent(t *testing.T) {
	lines := []string{"abcdefghij"}
	_, content, _, _ := layout(lines, 0, 0, 4, 10)

	want := []string{"abcd", "efgh", "ij"}
	if !reflect.DeepEqual(content, want) {
		t.Errorf("content = %v, want %v", content, want)
	}
}
