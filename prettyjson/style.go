// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>

package prettyjson

import "github.com/jeanfrancoisgratton/helperFunctions/v4/terminalfx"

// DefaultStyle returns a terminal-friendly style based on terminalfx.
//
// If you don't want any ANSI escape sequences, use PlainStyle().
func DefaultStyle() Style {
	return Style{
		Key:    terminalfx.Blue,
		String: terminalfx.Green,
		Number: terminalfx.Yellow,
		Bool:   terminalfx.Yellow,
		Null:   terminalfx.Red,
		Punct:  terminalfx.White,
	}
}

// PlainStyle disables colorization.
func PlainStyle() Style { return Style{} }
