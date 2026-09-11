// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: terminalfx/type_styles.go
// Original timestamp: 2026.09.10 22:28:21

package terminalfx

import "charm.land/lipgloss/v2"

// TEXT STYLE FUNCTIONS
// =====================
//
// Note on Blink: many modern terminal emulators (iTerm2, Windows Terminal,
// and others) ignore or intentionally disable the blink attribute. It's
// implemented here for spec-completeness but may render as no-op.

func Bold(sentence string) string {
	return lipgloss.NewStyle().Bold(true).Render(sentence)
}

func Italic(sentence string) string {
	return lipgloss.NewStyle().Italic(true).Render(sentence)
}

func Underline(sentence string) string {
	return lipgloss.NewStyle().Underline(true).Render(sentence)
}

func Reverse(sentence string) string {
	return lipgloss.NewStyle().Reverse(true).Render(sentence)
}

func Blink(sentence string) string {
	return lipgloss.NewStyle().Blink(true).Render(sentence)
}

func Strikethrough(sentence string) string {
	return lipgloss.NewStyle().Strikethrough(true).Render(sentence)
}
