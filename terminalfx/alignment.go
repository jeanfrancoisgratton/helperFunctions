// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: terminalfx/alignment.go
// Original timestamp: 2026.09.10 22:28:21

package terminalfx

import "charm.land/lipgloss/v2"

// ALIGNMENT FUNCTIONS
// ====================
//
// All three functions measure width via lipgloss (ANSI-escape and
// grapheme aware), so already-colourized or non-ASCII input aligns
// correctly instead of counting escape-sequence bytes as columns.

// Center centers a string within the terminal's width.
func Center(input string) string {
	return lipgloss.PlaceHorizontal(terminalWidthOrDefault(), lipgloss.Center, input)
}

// Left left-aligns a string within the terminal's width.
func Left(input string) string {
	return lipgloss.PlaceHorizontal(terminalWidthOrDefault(), lipgloss.Left, input)
}

// Right right-aligns a string within the terminal's width.
func Right(input string) string {
	return lipgloss.PlaceHorizontal(terminalWidthOrDefault(), lipgloss.Right, input)
}

func terminalWidthOrDefault() int {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80 // Default to 80 columns if size retrieval fails
	}
	return width
}
