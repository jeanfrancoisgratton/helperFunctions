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
// Lines wider than the terminal are hard-wrapped so the pagination math
// (page size, status bar position) stays accurate on narrow terminals.
// The terminal is re-measured on resize (SIGWINCH) so an in-progress
// session adapts without needing to be restarted; on Windows, where no
// such signal exists, the new size is picked up on the next keystroke.
//
// Keybindings active while paging:
//
//	Space   – advance one full page
//	Enter   – advance one line
//	Ctrl+U  – go back one full page
//	Ctrl+G  – show current position overlay; the key that dismisses it
//	          is then handled normally (e.g. Ctrl+G, Q quits)
//	Q / q   – quit and clear the screen
//
// Page restores the terminal to its original state before returning,
// whether it exits normally (Q) or because of an error.
func Page(lines []string, bannerSize, footerSize int) error {
	if len(lines) == 0 {
		return nil
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("pager: cannot enter raw mode: %w", err)
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	termWidth, termHeight := measureTerminal()
	banner, content, footer, pageSize := layout(lines, bannerSize, footerSize, termWidth, termHeight)
	totalLines := len(content)
	currentLine := 0

	// keyCh/errCh are fed by a background reader so the render loop can
	// also select on resize events. The reader goroutine outlives Page
	// if it's parked on a blocking Read when we return (no further
	// keystroke ever arrives) — harmless, it just sits idle.
	keyCh := make(chan byte)
	errCh := make(chan error, 1)
	go readKeys(keyCh, errCh)

	resizeCh, stopResize := watchResize()
	defer stopResize()

	recompute := func() {
		termWidth, termHeight = measureTerminal()
		banner, content, footer, pageSize = layout(lines, bannerSize, footerSize, termWidth, termHeight)
		totalLines = len(content)
		if currentLine >= totalLines {
			currentLine = totalLines - 1
		}
		if currentLine < 0 {
			currentLine = 0
		}
	}

	var pendingKey byte
	havePending := false

outer:
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

		var key byte
		if havePending {
			key = pendingKey
			havePending = false
		} else {
			select {
			case key = <-keyCh:
			case <-resizeCh:
				recompute()
				continue outer
			case <-errCh:
				break outer
			}
		}

		switch key {

		case keySpace: // ─── next page ──────────────────────────────────
			if !atEnd {
				currentLine += pageSize
				if currentLine >= totalLines {
					currentLine = totalLines - 1
				}
			}

		case keyCR, keyLF: // ─ next line ────────────────────────────────
			if !atEnd {
				currentLine++
			}

		case keyCtrlU: // ── previous page ───────────────────────────────
			currentLine -= pageSize
			if currentLine < 0 {
				currentLine = 0
			}

		case keyCtrlG: // ── position overlay ────────────────────────────
			// Overwrite the status bar in place; whatever key dismisses
			// the overlay becomes the next command (so e.g. Ctrl+G, Q
			// quits immediately instead of requiring Q twice).
			rawStatusClear()
			fmt.Printf("[Line %d/%d]  -- press any key --\r", currentLine+1, totalLines)
			select {
			case pendingKey = <-keyCh:
				havePending = true
			case <-resizeCh:
				recompute()
			case <-errCh:
				break outer
			}

		case keyQUIT, keyQUITu: // ─ quit ────────────────────────────────
			rawClear()
			return nil
		}
	}

	return nil
}

// measureTerminal returns the current terminal width and height, falling
// back to sane defaults when the size can't be determined.
func measureTerminal() (width, height int) {
	width, height = terminalfx.GetTerminalSize()
	if width <= 0 {
		width = defaultTermWidth
	}
	if height <= 0 {
		height = defaultTermHeight
	}
	return width, height
}

// readKeys reads single bytes from stdin and forwards them on keyCh until
// an error (e.g. EOF) occurs, at which point it reports the error on errCh
// and returns.
func readKeys(keyCh chan<- byte, errCh chan<- error) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			errCh <- err
			return
		}
		if n == 0 {
			continue
		}
		keyCh <- buf[0]
	}
}

// layout carves lines into sticky banner/footer regions and scrollable
// content, wrapping every region to termWidth so the pagination math
// (pageSize, status bar counts) stays accurate regardless of how long
// individual lines are. If the banner and footer would together fill the
// screen once wrapped, both are dropped and all lines become content.
func layout(lines []string, bannerSize, footerSize, termWidth, termHeight int) (banner, content, footer []string, pageSize int) {
	effectiveBanner := bannerSize
	effectiveFooter := footerSize
	if effectiveBanner < 0 {
		effectiveBanner = 0
	}
	if effectiveFooter < 0 {
		effectiveFooter = 0
	}
	if effectiveBanner > len(lines) {
		effectiveBanner = len(lines)
	}
	if effectiveFooter > len(lines)-effectiveBanner {
		effectiveFooter = len(lines) - effectiveBanner
	}

	bannerLines := lines[:effectiveBanner]
	contentEnd := len(lines) - effectiveFooter
	contentLines := lines[effectiveBanner:contentEnd]
	var footerLines []string
	if effectiveFooter > 0 {
		footerLines = lines[contentEnd:]
	}

	wrappedBanner := wrapLines(bannerLines, termWidth)
	wrappedFooter := wrapLines(footerLines, termWidth)

	if len(wrappedBanner)+len(wrappedFooter) >= termHeight {
		wrappedBanner = nil
		wrappedFooter = nil
		contentLines = lines
	}

	pageSize = termHeight - len(wrappedBanner) - len(wrappedFooter) - 1
	if pageSize < 1 {
		pageSize = 1
	}

	return wrappedBanner, wrapLines(contentLines, termWidth), wrappedFooter, pageSize
}

// wrapLines hard-wraps each line to width runes. width <= 0 disables
// wrapping. This does not attempt word-boundary wrapping or account for
// ANSI escape sequences embedded in a line.
func wrapLines(lines []string, width int) []string {
	if width <= 0 {
		return lines
	}
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		r := []rune(l)
		if len(r) == 0 {
			out = append(out, "")
			continue
		}
		for len(r) > width {
			out = append(out, string(r[:width]))
			r = r[width:]
		}
		out = append(out, string(r))
	}
	return out
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
