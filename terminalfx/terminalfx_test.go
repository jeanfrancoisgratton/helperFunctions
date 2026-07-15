package terminalfx

import (
	"strings"
	"testing"
)

func TestGetTerminalSize(t *testing.T) {
	cols, rows := GetTerminalSize()
	if cols < 0 || rows < 0 {
		t.Errorf("GetTerminalSize() = (%d, %d), want both >= 0", cols, rows)
	}
}

func TestCenter(t *testing.T) {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80
	}

	t.Run("short input uses computed padding", func(t *testing.T) {
		input := "hi"
		want := strings.Repeat(" ", (width-len(input))/2) + input
		if got := Center(input); got != want {
			t.Errorf("Center(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("input longer than any terminal has no padding", func(t *testing.T) {
		input := strings.Repeat("x", 100000)
		if got := Center(input); got != input {
			t.Errorf("Center() with an oversized input modified it: got len %d, want len %d unchanged", len(got), len(input))
		}
	})

	t.Run("empty input", func(t *testing.T) {
		want := strings.Repeat(" ", width/2)
		if got := Center(""); got != want {
			t.Errorf("Center(\"\") = %q, want %q", got, want)
		}
	})
}

func TestRight(t *testing.T) {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80
	}

	t.Run("short input uses computed padding", func(t *testing.T) {
		input := "hi"
		want := strings.Repeat(" ", width-len(input)) + input
		if got := Right(input); got != want {
			t.Errorf("Right(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("input longer than any terminal has no padding", func(t *testing.T) {
		input := strings.Repeat("x", 100000)
		if got := Right(input); got != input {
			t.Errorf("Right() with an oversized input modified it: got len %d, want len %d unchanged", len(got), len(input))
		}
	})
}

func TestGlyphSigns(t *testing.T) {
	signs := map[string]func(string) string{
		"EuropeanStopSign":      EuropeanStopSign,
		"AmericanStopSign":      AmericanStopSign,
		"SkullBonesSign":        SkullBonesSign,
		"EnabledSign":           EnabledSign,
		"ErrorSign":             ErrorSign,
		"GreenGoSign":           GreenGoSign,
		"InProgressSign":        InProgressSign,
		"WarningSign":           WarningSign,
		"InfoSign":              InfoSign,
		"NoteSign":              NoteSign,
		"ScrollSign":            ScrollSign,
		"TipSign":               TipSign,
		"LightbulbSign":         LightbulbSign,
		"ThumbsUpSign":          ThumbsUpSign,
		"ThumbsDownSign":        ThumbsDownSign,
		"NotExistMathSign":      NotExistMathSign,
		"ExistMathSign":         ExistMathSign,
		"NotIncludedInMathSign": NotIncludedInMathSign,
		"IsIncludedInMathSign":  IsIncludedInMathSign,
		"DeltaSymbolMathSign":   DeltaSymbolMathSign,
	}

	for name, fn := range signs {
		t.Run(name, func(t *testing.T) {
			got := fn("hello")
			if !strings.Contains(got, "hello") {
				t.Errorf("%s(%q) = %q, want it to contain the original sentence", name, "hello", got)
			}
		})
	}
}

func TestBombSign(t *testing.T) {
	for _, coloured := range []bool{true, false} {
		got := BombSign("boom", coloured)
		if !strings.Contains(got, "boom") {
			t.Errorf("BombSign(%q, %v) = %q, want it to contain the original sentence", "boom", coloured, got)
		}
	}
}

func TestColorFunctions(t *testing.T) {
	colors := map[string]func(string) string{
		"Red":    Red,
		"Green":  Green,
		"White":  White,
		"Yellow": Yellow,
		"Blue":   Blue,
	}

	for name, fn := range colors {
		t.Run(name, func(t *testing.T) {
			got := fn("hello")
			if !strings.Contains(got, "hello") {
				t.Errorf("%s(%q) = %q, want it to contain the original text", name, "hello", got)
			}
		})
	}
}
