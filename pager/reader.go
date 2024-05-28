package pager

import cerr "github.com/jeanfrancoisgratton/customError"

import (
	"os"
)

// This is just a wrapper that reads a text file and then calls up and executes the Pager
// The idea here is that you call the pager with the text to be paged, and a terminal height value
// This function here reads a text file and converts it in a single string (text),
// and then call upon Pager. If you already have a text block to paginate, all you need then is to call
// NewPager with the text block and the desired terminal height (a value of 0 means: default term height)
func PaginateFromFile(fileFullPath string, pageHeight int) *cerr.CustomError {
	data, err := os.ReadFile(fileFullPath)
	if err != nil {
		return &cerr.CustomError{Title: "Error reading the file:", Message: err.Error()}
	}

	// Initialize the pager with the text and desired page height
	text := string(data)
	p, err := NewPager(text, pageHeight) // 0 uses the default terminal height - 1
	if err != nil {
		return err
	}

	// If p is nil, it means the text was printed directly without paging
	if p != nil {
		// Run the pager
		p.Run()
	}
	return nil
}
