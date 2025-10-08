// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/23 03:16
// Original filename: /terminal-glyphs.go

package terminalfx

import (
	"fmt"
)

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
func EuropeanStopPanelGlyph(sentence string) string { return fmt.Sprintf("⛔ %s%s", sentence, reset) } // European stop
func AmericanStopPanelGlyph(sentence string) string { return fmt.Sprintf("🛑 %s%s", sentence, reset) } // American stop

// Fatal symbols
func FatalCollisionGlyph(sentence string) string  { return fmt.Sprintf("💥 %s%s", sentence, reset) }
func FatalSkullBonesGlyph(sentence string) string { return fmt.Sprintf("☠ %s%s", sentence, reset) }

// Go
func GreenGoGlyph(sentence string) string { return fmt.Sprintf("🟢 %s%s", sentence, reset) }

// Status / utility// InProgressChar returns a single-glyph indicator for "task in progress".
func InProgressGlyph(sentence string) string { return fmt.Sprintf("⏳ %s%s", sentence, reset) } // U+23F3 HOURGLASS NOT DONE
func EnabledGlyph(sentence string) string    { return fmt.Sprintf("✅ %s%s", sentence, reset) }
func ErrorGlyph(sentence string) string      { return fmt.Sprintf("❌ %s%s", sentence, reset) }
func WarningGlyph(sentence string) string    { return fmt.Sprintf("⚠ %s%s", sentence, reset) } // U+26A0 WARNING SIGN
func InfoGlyph(sentence string) string       { return fmt.Sprintf("🛈 %s%s", sentence, reset) } // circled info (U+1F6C8)
func NoteGlyph(sentence string) string       { return fmt.Sprintf("💬 %s%s", sentence, reset) } // speech bubble
func ScrollGlyph(sentence string) string     { return fmt.Sprintf("📜 %s%s", sentence, reset) } // scroll/document
func TipGlyph(sentence string) string        { return fmt.Sprintf("💡 %s%s", sentence, reset) }
func LightbulbGlyph(sentence string) string  { return fmt.Sprintf("💡 %s%s", sentence, reset) }
func ThumbsUpGlyph(sentence string) string   { return fmt.Sprintf("👍 %s%s", sentence, reset) }
func ThumbsDownGlyph(sentence string) string { return fmt.Sprintf("👎 %s%s", sentence, reset) }

// ===== Colored variants =====

func RedErrorGlyph(sentence string) string {
	return fmt.Sprintf("%s%s", red, ErrorGlyph(sentence))
}

func YellowWarningGlyph(sentence string) string {
	return fmt.Sprintf("%s⚠ %s%s", yellow, sentence, reset)
}

func GreenOkGlyph(sentence string) string {
	return fmt.Sprintf("%s%s", green, GreenGoGlyph(sentence))
}

func BlueInfoGlyph(sentence string) string {
	return fmt.Sprintf("%s%s%s", blue, InfoGlyph(sentence))
}

func BlueInProgressGlyph(sentence string) string {
	return fmt.Sprintf("%s%s%s", blue, InProgressGlyph(sentence))
}

func YellowTipGlyph(sentence string) string {
	return fmt.Sprintf("%s%s%s", yellow, TipGlyph(sentence))
}
