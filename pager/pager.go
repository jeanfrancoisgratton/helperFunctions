// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: pager/pager.go

package pager

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// Page paginates lines in a manner similar to the Unix `more` command.
//
// The first bannerSize lines and the last footerSize lines of lines are
// treated as sticky regions: they are reprinted on every page and are never
// part of the scrollable content window.  If bannerSize+footerSize is equal
// to or greater than the terminal height both values are ignored and all
// lines are treated as scrollable content.
//
// Keybindings active while paging:
//
//	Space   – advance one full page
//	Enter   – advance one line
//	Ctrl+U  – go back one full page
//	Ctrl+G  – show current position overlay (clears on next redraw)
//	Q / q   – quit and clear the screen
//
// Page restores the terminal to its original state before returning,
// whether it exits normally (Q) or because of an error.
func Page(lines []string, bannerSize, footerSize int) error {
	if len(lines) == 0 {
		return nil
	}

	// Determine usable terminal height.
	_, termHeight := terminalfx.GetTerminalSize()
	if termHeight <= 0 {
		termHeight = defaultTermHeight
	}

	// If sticky regions would crowd out all content rows, drop them.
	effectiveBanner := bannerSize
	effectiveFooter := footerSize
	if effectiveBanner < 0 {
		effectiveBanner = 0
	}
	if effectiveFooter < 0 {
		effectiveFooter = 0
	}
	if effectiveBanner+effectiveFooter >= termHeight {
		effectiveBanner = 0
		effectiveFooter = 0
	}

	// Clamp sticky sizes to the actual slice length.
	if effectiveBanner > len(lines) {
		effectiveBanner = len(lines)
	}
	if effectiveFooter > len(lines)-effectiveBanner {
		effectiveFooter = len(lines) - effectiveBanner
	}

	// Carve up the slice.
	banner := lines[:effectiveBanner]

	contentEnd := len(lines) - effectiveFooter
	content := lines[effectiveBanner:contentEnd]

	var footer []string
	if effectiveFooter > 0 {
		footer = lines[contentEnd:]
	}

	// One row is always reserved for the status bar.
	pageSize := termHeight - effectiveBanner - effectiveFooter - 1
	if pageSize < 1 {
		pageSize = 1
	}

	totalLines := len(content)

	// If everything fits on one screen we still run the loop so the user
	// sees the "(END)" prompt and can quit intentionally.

	// Switch stdin to raw mode for single-keystroke reads.
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("pager: cannot enter raw mode: %w", err)
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	currentLine := 0
	buf := make([]byte, 4)

	for {
		// ── Render frame ──────────────────────────────────────────────────────

		rawClear()

		for _, l := range banner {
			rawPrintln(l)
		}

		end := currentLine + pageSize
		if end > totalLines {
			end = totalLines
		}
		for i := currentLine; i < end; i++ {
			rawPrintln(content[i])
		}

		for _, l := range footer {
			rawPrintln(l)
		}

		// ── Status bar ────────────────────────────────────────────────────────

		atEnd := end >= totalLines
		if atEnd {
			rawStatus("(END) [%d/%d]  Q to quit", totalLines, totalLines)
		} else {
			rawStatus("-- More -- [%d/%d]  Space=▶page  Enter=▶line  ^U=◀page  ^G=pos  Q=quit",
				end, totalLines)
		}

		// ── Keystroke ─────────────────────────────────────────────────────────

		n, readErr := os.Stdin.Read(buf)
		if readErr != nil || n == 0 {
			break
		}

		switch buf[0] {

		case keySpace: // ─── next page ──────────────────────────────────────
			if !atEnd {
				currentLine += pageSize
				if currentLine >= totalLines {
					currentLine = totalLines - 1
				}
			}

		case keyCR, keyLF: // ─ next line ────────────────────────────────────
			if !atEnd {
				currentLine++
			}

		case keyCtrlU: // ── previous page ───────────────────────────────────
			currentLine -= pageSize
			if currentLine < 0 {
				currentLine = 0
			}

		case keyCtrlG: // ── position overlay ────────────────────────────────
			// Overwrite the status bar in place; a subsequent keypress will
			// cause a full redraw, restoring the normal status line.
			rawStatusClear()
			fmt.Printf("[Line %d/%d]  -- press any key --\r", currentLine+1, totalLines)
			os.Stdin.Read(buf) //nolint:errcheck — ignore the consumed key

		case keyQUIT, keyQUITu: // ─ quit ────────────────────────────────────
			term.Restore(fd, oldState) //nolint:errcheck
			rawClear()
			return nil
		}
	}

	return nil
}

// ── terminal helpers (raw-mode safe) ─────────────────────────────────────────

// rawClear clears the screen and homes the cursor using ANSI escapes.
// This mirrors terminalfx.ClearTTY but is kept local so the pager
// package has no hidden coupling to that package's unexported symbols.
func rawClear() {
	fmt.Print("\x1b[2J\x1b[H")
}

// rawPrintln writes s followed by CR+LF, which is required in raw mode
// because the terminal driver no longer translates \n into CR+LF.
func rawPrintln(s string) {
	fmt.Printf("%s\r\n", s)
}

// rawStatus writes a formatted status bar at the current cursor position
// without appending a newline, so it stays on the last row.
func rawStatus(format string, args ...any) {
	fmt.Printf(format, args...)
}

// rawStatusClear moves to the start of the current line and erases it.
func rawStatusClear() {
	fmt.Print("\r\x1b[K")
}
