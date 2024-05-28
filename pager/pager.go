package pager

import (
	"bufio"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"golang.org/x/term"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Pager represents the pager functionality
type Pager struct {
	lines       []string
	pageHeight  int
	currentPage int
	currentLine int
}

// NewPager initializes a new Pager
// This is the entry point, unless you use the helper FilePager() in reader.go
// A height value less than 1 means that you'd use the terminal's current height
// If that value is less than 10 (arbitrary call, here), there is no paging the text block is displayed as is
func NewPager(text string, height int) (*Pager, *cerr.CustomError) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, &cerr.CustomError{Fatality: cerr.Warning, Message: "No lines in buffer"}
	}

	termHeight, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termHeight = 24 // default terminal height
	}
	if height <= 0 {
		height = termHeight - 1
	}

	if termHeight < 10 || height < 10 {
		//fmt.Println("Terminal height or user-supplied height is below 10. Printing without paging:")
		for _, line := range lines {
			fmt.Println(line)
		}
		return nil, nil
	}

	return &Pager{lines: lines, pageHeight: height - 1, currentPage: 0, currentLine: 0}, nil
}

// DisplayPage displays the current page
func (p *Pager) DisplayPage() {
	start := p.currentPage * p.pageHeight
	end := start + p.pageHeight
	if end > len(p.lines) {
		end = len(p.lines)
	}
	for i := start; i < end; i++ {
		fmt.Println(p.lines[i])
	}
}

// DisplayLine displays the current line
func (p *Pager) DisplayLine() {
	if p.currentLine < len(p.lines) {
		fmt.Println(p.lines[p.currentLine])
		p.currentLine++
	}
}

// Run starts the pager functionality
func (p *Pager) Run() {
	reader := bufio.NewReader(os.Stdin)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	for {
		p.DisplayPage()
		fmt.Print("--- More ---")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case " ":
			p.currentPage++
		case "\n":
			p.DisplayLine()
		case "q":
			return
		case "u":
			if p.currentPage > 0 {
				p.currentPage--
			}
		case "y":
			if p.currentLine > 0 {
				p.currentLine--
				p.DisplayLine()
			}
		case "n":
			fmt.Printf("Line: %d/%d (%.2f%%)\n", p.currentLine, len(p.lines), float64(p.currentLine)/float64(len(p.lines))*100)
		}

		select {
		case <-signalChan:
			fmt.Println("Program exited by CTRL+C")
			os.Exit(1)
		default:
		}
	}
}
