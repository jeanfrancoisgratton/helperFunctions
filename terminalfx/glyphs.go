// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/23 03:16
// Original filename: /terminal-glyphs.go

package terminalfx

import "fmt"

// ANSI color codes
const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

// ===== Base glyph functions =====

// Stop signs
func EuropeanStopPanelGlyph() string { return "⛔" } // European stop
func AmericanStopPanelGlyph() string { return "🛑" } // American stop

// Fatal symbols
func FatalCollisionGlyph() string  { return "💥" }
func FatalSkullBonesGlyph() string { return "☠" }

// Go
func GreenGoGlyph() string { return "🟢" }

// Status / utility
func EnabledGlyph() string    { return "✅" }
func ErrorGlyph() string      { return "❌" }
func WarningGlyph() string    { return "⚠" } // U+26A0 WARNING SIGN
func InfoGlyph() string       { return "🛈" } // circled info (U+1F6C8)
func NoteGlyph() string       { return "💬" } // speech bubble
func ScrollGlyph() string     { return "📜" } // scroll/document
func TipGlyph() string        { return "💡" }
func LightbulbGlyph() string  { return "💡" }
func ThumbsUpGlyph() string   { return "👍" }
func ThumbsDownGlyph() string { return "👎" }

// ===== Colored variants =====

func RedErrorGlyph() string {
	return fmt.Sprintf("%s%s%s", red, ErrorGlyph(), reset)
}

func YellowWarningGlyph() string {
	return fmt.Sprintf("%s⚠%s", yellow, reset)
}

func GreenOkGlyph() string {
	return fmt.Sprintf("%s%s%s", green, GreenGoGlyph(), reset)
}

func BlueInfoGlyph() string {
	return fmt.Sprintf("%s%s%s", blue, InfoGlyph(), reset)
}

func YellowTipGlyph() string {
	return fmt.Sprintf("%s%s%s", yellow, TipGlyph(), reset)
}

// ===== Map-based helper =====

// Glyph returns the glyph string for a given kind.
// Supported kinds: "stop-eu", "stop-us", "fatal-collision", "fatal-skull",
// "go", "error", "info", "note", "scroll", "tip", "lightbulb",
// "thumbs-up", "thumbs-down".
func Glyph(kind string) string {
	glyphs := map[string]string{
		"enabled":         EnabledGlyph(),
		"warning":         WarningGlyph(),
		"stop-eu":         EuropeanStopPanelGlyph(),
		"stop-us":         AmericanStopPanelGlyph(),
		"fatal-collision": FatalCollisionGlyph(),
		"fatal-skull":     FatalSkullBonesGlyph(),
		"go":              GreenGoGlyph(),
		"error":           ErrorGlyph(),
		"info":            InfoGlyph(),
		"note":            NoteGlyph(),
		"scroll":          ScrollGlyph(),
		"tip":             TipGlyph(),
		"lightbulb":       LightbulbGlyph(),
		"thumbs-up":       ThumbsUpGlyph(),
		"thumbs-down":     ThumbsDownGlyph(),
	}

	if g, ok := glyphs[kind]; ok {
		return g
	}
	return "?" // fallback if unknown
}
