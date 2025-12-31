// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /json/types.go
// Original timestamp: 2025/12/31

package prettyjson

import "io"

// ColorMode controls when ANSI colorization is used.
type ColorMode uint8

const (
	// ColorAuto enables color only when the target writer is a terminal.
	ColorAuto ColorMode = iota
	// ColorAlways forces color even when the writer is not a terminal.
	ColorAlways
	// ColorNever disables color.
	ColorNever
)

// Style is a set of colorizer functions.
//
// Each function receives the raw token string and returns the potentially
// colorized version.
//
// Any nil function means "no color" for that token type.
type Style struct {
	Key    func(string) string // object keys (including quotes)
	String func(string) string // string values (including quotes)
	Number func(string) string // numbers
	Bool   func(string) string // true/false
	Null   func(string) string // null
	Punct  func(string) string // punctuation: { } [ ] , :
}

// Options controls formatting and output.
type Options struct {
	Writer   io.Writer
	Indent   string
	SortKeys bool
	Color    ColorMode
	Style    Style
}

// Option is a functional option for Print/SPrint.
type Option func(*Options)

func defaultOptions() Options {
	return Options{
		Writer:   nil, // resolved in Print
		Indent:   "  ",
		SortKeys: true,
		Color:    ColorAuto,
		Style:    DefaultStyle(),
	}
}

// WithWriter sets the output writer (defaults to os.Stdout for Print).
func WithWriter(w io.Writer) Option {
	return func(o *Options) { o.Writer = w }
}

// WithIndent sets the indentation string (defaults to two spaces).
func WithIndent(indent string) Option {
	return func(o *Options) { o.Indent = indent }
}

// WithSortKeys enables or disables object-key sorting.
func WithSortKeys(sort bool) Option {
	return func(o *Options) { o.SortKeys = sort }
}

// WithColorMode sets colorization mode.
func WithColorMode(m ColorMode) Option {
	return func(o *Options) { o.Color = m }
}

// WithStyle overrides the current style.
func WithStyle(s Style) Option {
	return func(o *Options) { o.Style = s }
}
