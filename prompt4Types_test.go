package helperFunctions

import (
	"os"
	"reflect"
	"testing"
)

// withStdin temporarily replaces os.Stdin with a pipe fed with input, runs fn,
// then restores the original os.Stdin.
func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	orig := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = orig
	}()

	go func() {
		_, _ = w.WriteString(input)
		_ = w.Close()
	}()

	fn()
	_ = r.Close()
}

func TestGetStringValFromPrompt(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"normal line", "hello\n", "hello"},
		{"empty line", "\n", ""},
		{"line with spaces", "  spaced  \n", "  spaced  "},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got string
			withStdin(t, c.input, func() {
				got = GetStringValFromPrompt("prompt: ")
			})
			if got != c.want {
				t.Errorf("GetStringValFromPrompt() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestGetIntValFromPrompt(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int
	}{
		{"valid int", "42\n", 42},
		{"negative int", "-7\n", -7},
		{"empty line defaults to zero", "\n", 0},
		{"non-numeric falls back to one", "notanumber\n", 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			withStdin(t, c.input, func() {
				got = GetIntValFromPrompt("prompt: ")
			})
			if got != c.want {
				t.Errorf("GetIntValFromPrompt() = %d, want %d", got, c.want)
			}
		})
	}
}

func TestGetBoolValFromPrompt(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"true word", "true\n", true},
		{"True capitalized", "True\n", true},
		{"one", "1\n", true},
		{"false word", "false\n", false},
		{"zero", "0\n", false},
		{"arbitrary word", "yes\n", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got bool
			withStdin(t, c.input, func() {
				got = GetBoolValFromPrompt("prompt: ")
			})
			if got != c.want {
				t.Errorf("GetBoolValFromPrompt() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestGetStringSliceFromPrompt(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{"multiple entries", "a\nb\nc\n\n", []string{"a", "b", "c"}},
		{"immediate terminator", "\n", []string{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []string
			withStdin(t, c.input, func() {
				got = GetStringSliceFromPrompt("prompt")
			})
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetStringSliceFromPrompt() = %#v, want %#v", got, c.want)
			}
		})
	}
}

func TestGetValueFromPrompt(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  any
	}{
		{"unsigned integer", "123\n", uint(123)},
		{"signed negative integer", "-5\n", int(-5)},
		{"boolean", "true\n", true},
		{"plain string", "hello\n", "hello"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got any
			withStdin(t, c.input, func() {
				got = GetValueFromPrompt("prompt: ")
			})
			if got != c.want {
				t.Errorf("GetValueFromPrompt() = %#v (%T), want %#v (%T)", got, got, c.want, c.want)
			}
		})
	}
}
