// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /terminal.go
// Original timestamp: 2024/04/10 15:21

package helperFunctions

import (
	"fmt"
	"github.com/jwalton/gchalk"
	"strings"
	"syscall"
	"unsafe"
)

const (
	terminalEscape = "\x1b"
)

// TERMINAL FUNCTIONS
func GetTerminalSize() (int, int) {
	var size struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&size)))
	if err != 0 {
		return 0, 0
	}
	return int(size.cols), int(size.rows)
}

// Yeah... I know.. nobody should clear a TTY in-tool... :p
func ClearTTY() {
	fmt.Print("\x1b[2J") // Clears screen
	fmt.Print("\x1b[H")  // Moves cursor to top-left corner
}

// Center string in the terminal
func Center(input string) string {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80 // Default to 80 columns if size retrieval fails
	}

	// Calculate the padding needed to center the text
	padding := (width - len(input)) / 2
	if padding < 0 {
		padding = 0 // Avoid negative padding if the string is longer than the terminal width
	}

	return strings.Repeat(" ", padding) + input
}

// Right aligns the input string to the right side of the terminal.
func Right(input string) string {
	width, _ := GetTerminalSize()
	if width == 0 {
		width = 80 // Default width if retrieval fails
	}

	// Calculate the padding needed to right-align the text
	padding := width - len(input)
	if padding < 0 {
		padding = 0 // Avoid negative padding if the string is longer than the terminal width
	}

	return strings.Repeat(" ", padding) + input
}

// COLOR FUNCTIONS
// ===============
func Red(sentence string) string {
	return fmt.Sprintf("%s", gchalk.WithBrightRed().Bold(sentence))
}

func Green(sentence string) string {
	return fmt.Sprintf("%s", gchalk.WithBrightGreen().Bold(sentence))
}

func White(sentence string) string {
	return fmt.Sprintf("%s", gchalk.WithBrightWhite().Bold(sentence))
}

func Yellow(sentence string) string {
	return fmt.Sprintf("%s", gchalk.WithBrightYellow().Bold(sentence))
}

func Blue(sentence string) string {
	return fmt.Sprintf("%s", gchalk.WithBlue().Bold(sentence))
}
