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
	Paging   bool
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
		Paging:   false,
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

// WithPaging tells Print to route its output through the pager subpackage
// when the rendered content doesn't fit on one screen.
//
// This is only ever a hint: it's honored by Print, and only when the
// output writer is the default (os.Stdout) and the process is attached to
// a terminal on both stdin and stdout (see terminalfx.Interactive). In any
// other situation — a custom writer, redirected/piped I/O, or content
// that already fits on one screen — it's silently ignored and Print
// behaves exactly as it would with Paging left false.
//
// SPrint and Format ignore this option entirely: paging is an interactive,
// side-effecting operation and doesn't fit their "return the formatted
// text" contract.
func WithPaging(enabled bool) Option {
	return func(o *Options) { o.Paging = enabled }
}
