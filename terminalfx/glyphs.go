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

// Stop signs
func EuropeanStopSign(sentence string) string { return fmt.Sprintf("⛔ %s%s", reset, sentence) }

func AmericanStopSign(sentence string) string { return fmt.Sprintf("🛑 %s%s", reset, sentence) }

// Fatal error signs
func BombSign(sentence string, coloured bool) string {
	return fmt.Sprintf("💥 %s%s", reset, sentence)
}

func SkullBonesSign(sentence string) string { return fmt.Sprintf("☠ %s%s", reset, sentence) }

// Accepted/denied
func EnabledSign(sentence string) string { return fmt.Sprintf("✅ %s%s", reset, sentence) }
func ErrorSign(sentence string) string   { return fmt.Sprintf("❌ %s%s", reset, sentence) }

// Status / utility// InProgressChar returns a single-glyph indicator for "task in progress".
func GreenGoSign(sentence string) string { return fmt.Sprintf("🟢 %s%s", reset, sentence) }

func InProgressSign(sentence string) string { return fmt.Sprintf("⏳ %s%s", reset, sentence) } // U+23F3 HOURGLASS NOT DONE
func WarningSign(sentence string) string    { return fmt.Sprintf("⚠ %s%s", reset, sentence) } // U+26A0 WARNING SIGN
func InfoSign(sentence string) string       { return fmt.Sprintf("🛈 %s%s", reset, sentence) } // circled info (U+1F6C8)
func NoteSign(sentence string) string       { return fmt.Sprintf("💬 %s%s", reset, sentence) } // speech bubble
func ScrollSign(sentence string) string     { return fmt.Sprintf("📜 %s%s", reset, sentence) } // scroll/document
func TipSign(sentence string) string        { return fmt.Sprintf("💡 %s%s", reset, sentence) }
func LightbulbSign(sentence string) string  { return fmt.Sprintf("💡 %s%s", reset, sentence) }
func ThumbsUpSign(sentence string) string   { return fmt.Sprintf("👍 %s%s", reset, sentence) }
func ThumbsDownSign(sentence string) string { return fmt.Sprintf("👎 %s%s", reset, sentence) }
