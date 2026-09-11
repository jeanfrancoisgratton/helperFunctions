package terminalfx

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
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

	t.Run("short input is centered within terminal width", func(t *testing.T) {
		input := "hi"
		got := Center(input)
		if len(got) != width {
			t.Fatalf("Center(%q) length = %d, want %d", input, len(got), width)
		}
		if strings.TrimSpace(got) != input {
			t.Errorf("Center(%q) = %q, content lost after trimming padding", input, got)
		}
		leading := len(got) - len(strings.TrimLeft(got, " "))
		trailing := len(got) - len(strings.TrimRight(got, " "))
		if d := leading - trailing; d < -1 || d > 1 {
			t.Errorf("Center(%q) padding is unbalanced: leading=%d trailing=%d", input, leading, trailing)
		}
	})

	t.Run("input longer than any terminal has no padding", func(t *testing.T) {
		input := strings.Repeat("x", 100000)
		if got := Center(input); got != input {
			t.Errorf("Center() with an oversized input modified it: got len %d, want len %d unchanged", len(got), len(input))
		}
	})

	t.Run("empty input produces only padding up to terminal width", func(t *testing.T) {
		got := Center("")
		if len(got) != width {
			t.Errorf("Center(\"\") length = %d, want %d", len(got), width)
		}
		if strings.TrimSpace(got) != "" {
			t.Errorf("Center(\"\") = %q, want only whitespace", got)
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
		"Red":        Red,
		"DimRed":     DimRed,
		"Green":      Green,
		"DimGreen":   DimGreen,
		"White":      White,
		"DimWhite":   DimWhite,
		"Yellow":     Yellow,
		"DimYellow":  DimYellow,
		"Blue":       Blue,
		"DimBlue":    DimBlue,
		"Magenta":    Magenta,
		"DimMagenta": DimMagenta,
		"Cyan":       Cyan,
		"DimCyan":    DimCyan,
		"Brown":      Brown,
		"DimBrown":   DimBrown,
		"Copper":     Copper,
		"DimCopper":  DimCopper,
		"Bronze":     Bronze,
		"DimBronze":  DimBronze,
		"Gold":       Gold,
		"DimGold":    DimGold,
		"Silver":     Silver,
		"DimSilver":  DimSilver,
		"Gray":       Gray,
		"DimGray":    DimGray,
		"Orange":     Orange,
		"DimOrange":  DimOrange,
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

func TestTypeStyleFunctions(t *testing.T) {
	styles := map[string]func(string) string{
		"Bold":          Bold,
		"Italic":        Italic,
		"Underline":     Underline,
		"Reverse":       Reverse,
		"Blink":         Blink,
		"Strikethrough": Strikethrough,
	}

	for name, fn := range styles {
		t.Run(name, func(t *testing.T) {
			got := fn("hello")
			// Underline/Strikethrough style non-space runes individually
			// (so spaces stay unstyled), so strip ANSI codes before checking.
			if !strings.Contains(ansi.Strip(got), "hello") {
				t.Errorf("%s(%q) = %q, want it to contain the original text", name, "hello", got)
			}
		})
	}
}

func TestLeft(t *testing.T) {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80
	}

	t.Run("short input is padded on the right up to terminal width", func(t *testing.T) {
		input := "hi"
		want := input + strings.Repeat(" ", width-len(input))
		if got := Left(input); got != want {
			t.Errorf("Left(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("input longer than any terminal has no padding", func(t *testing.T) {
		input := strings.Repeat("x", 100000)
		if got := Left(input); got != input {
			t.Errorf("Left() with an oversized input modified it: got len %d, want len %d unchanged", len(got), len(input))
		}
	})
}

func TestHorizontalLine(t *testing.T) {
	t.Run("explicit width", func(t *testing.T) {
		want := strings.Repeat("─", 10)
		if got := HorizontalLine(10); got != want {
			t.Errorf("HorizontalLine(10) = %q, want %q", got, want)
		}
	})

	t.Run("zero width uses terminal width", func(t *testing.T) {
		width, _ := GetTerminalSize()
		if width == 0 {
			width = 80
		}
		want := strings.Repeat("─", width)
		if got := HorizontalLine(0); got != want {
			t.Errorf("HorizontalLine(0) = %q, want %q", got, want)
		}
	})
}
