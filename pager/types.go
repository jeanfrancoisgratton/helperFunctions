// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: pager/types.go

// Package pager provides a terminal pager similar to the Unix `more` command.
//
// The main entry point is [Page], which accepts a slice of lines and optional
// sticky banner/footer sizes.  Banner lines are taken from the top of the
// slice; footer lines are taken from the bottom.  If their combined size would
// leave no room for content the sticky regions are silently dropped.
//
// Supported keystrokes while paging:
//
//	Space   – advance one full page
//	Enter   – advance one line
//	Ctrl+U  – go back one full page
//	Ctrl+G  – display current line position without redrawing
//	Q / q   – quit the pager
package pager

// keyCode constants for the keystrokes handled by the pager.
const (
	keyCtrlU = 0x15 // Ctrl+U  – page back
	keyCtrlG = 0x07 // Ctrl+G  – show position
	keySpace = 0x20 // Space   – page forward
	keyCR    = 0x0D // Enter (carriage return)
	keyLF    = 0x0A // Enter (line feed)
	keyQUIT  = 'q'
	keyQUITu = 'Q'
)

// defaultTermHeight is used when the terminal size cannot be determined.
const defaultTermHeight = 24
