// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: terminalfx/terminal.go
// Original timestamp: 2024/04/10 15:21

package terminalfx

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// TERMINAL FUNCTIONS
func GetTerminalSize() (int, int) {
	cols, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0
	}
	return cols, rows
}

// IsTerminal reports whether f is connected to a terminal.
func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// Interactive reports whether the process is attached to a terminal on
// both stdin and stdout. Callers use this to decide whether interactive
// behavior (colourized output, paging, prompts) is appropriate, as
// opposed to running with redirected/piped I/O.
func Interactive() bool {
	return IsTerminal(os.Stdin) && IsTerminal(os.Stdout)
}

// Yeah... I know... nobody should clear a TTY in-tool... :p

func ClearTTY() {
	fmt.Print("\x1b[2J") // Clears screen
	fmt.Print("\x1b[H")  // Moves cursor to top-left corner
}

// HorizontalLine returns a horizontal rule made of Unicode box-drawing
// characters. A width of 0 uses the current terminal's width.
func HorizontalLine(width int) string {
	if width == 0 {
		width, _ = GetTerminalSize()
		if width == 0 {
			width = 80
		}
	}
	return strings.Repeat("─", width)
}
