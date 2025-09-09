package helperFunctions

import (
	"bufio"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError/v2"
	"io"
	"os"
	"strings"
	"syscall"
	"unicode"

	"golang.org/x/term"
)

// Pager displays a slice of lines in a paginated manner.
func Pager(lines []string, linesPerPage int) {
	totalLines := len(lines)

	if linesPerPage <= 0 {
		if h, _, err := term.GetSize(int(syscall.Stdin)); err == nil {
			linesPerPage = h - 3
		} else {
			linesPerPage = 20
		}
	}

	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Failed to set raw terminal mode:", err)
		return
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	tty := os.NewFile(uintptr(syscall.Stdin), "/dev/tty")
	defer tty.Close()

	currentLine := 0
	for currentLine < totalLines {
		end := currentLine + linesPerPage
		if end > totalLines {
			end = totalLines
		}

		for _, line := range lines[currentLine:end] {
			fmt.Println(line)
		}
		currentLine = end

		if currentLine >= totalLines {
			break
		}

		fmt.Print("-- More -- (SPACE: page, ENTER: line, Q: quit) ")
		var buf [1]byte
		os.Stdin.Read(buf[:])
		fmt.Print("\r\033[K") // clear line

		switch unicode.ToLower(rune(buf[0])) {
		case 'q':
			return
		case '\r', '\n':
			if currentLine < totalLines {
				fmt.Println(lines[currentLine])
				currentLine++
			}
		case ' ':
			// Full page already printed
		default:
			// Treat other keys as ENTER
			if currentLine < totalLines {
				fmt.Println(lines[currentLine])
				currentLine++
			}
		}
	}
}

// PagerFromString splits a full string into lines, then pages it.
func PagerFromString(text string, linesPerPage int) {
	Pager(strings.Split(text, "\n"), linesPerPage)
}

// PagerFromFile reads a file and pages its content.
func PagerFromFile(filePath string, linesPerPage int) *cerr.CustomError {
	f, err := os.Open(filePath)
	if err != nil {
		return &cerr.CustomError{Title: "Error opening file", Message: err.Error()}
	}
	defer f.Close()

	return PagerFromReader(f, linesPerPage)
}

// PagerFromReader reads from an io.Reader and pages its content.
func PagerFromReader(r io.Reader, linesPerPage int) *cerr.CustomError {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return &cerr.CustomError{Title: "Error reading from reader", Message: err.Error()}
	}
	Pager(lines, linesPerPage)
	return nil
}

// AutoPager dispatches to the appropriate PagerFrom* function based on input type.
func AutoPager(input interface{}, linesPerPage int) *cerr.CustomError {
	if linesPerPage == 0 {
		height, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return &cerr.CustomError{Title: "Unable to get the terminal size", Message: err.Error()}
		}
		linesPerPage = height
	}
	switch v := input.(type) {
	case string:
		// Try to treat string as a file path first
		if fileInfo, err := os.Stat(v); err == nil && !fileInfo.IsDir() {
			return PagerFromFile(v, linesPerPage)
		}
		// Fallback: treat as raw text
		PagerFromString(v, linesPerPage)
		return nil

	case []string:
		Pager(v, linesPerPage)
		return nil

	case io.Reader:
		return PagerFromReader(v, linesPerPage)

	default:
		return &cerr.CustomError{Title: "Unsupported input type", Message: fmt.Sprintf("%T", input)}
	}
}
