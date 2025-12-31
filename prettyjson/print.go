// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /json/print.go
// Original timestamp: 2025/12/31

package prettyjson

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"encoding/json"
)

// Print pretty-prints a JSON payload to the configured writer.
//
// By default:
//   - indentation is two spaces
//   - object keys are sorted for stable output
//   - colors are enabled only when printing to a terminal
func Print(payload []byte, opts ...Option) error {
	o := defaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}

	w := o.Writer
	if w == nil {
		w = os.Stdout
	}

	v, err := parsePayload(payload)
	if err != nil {
		return err
	}

	p := newPrinter(w, o)
	if err := p.printValue(v, 0); err != nil {
		return err
	}
	return p.writeString("\n")
}

// SPrint pretty-prints a JSON payload and returns it as a string.
//
// Note: With the default ColorAuto mode, the returned string will be
// uncolored (because the underlying writer is not a terminal). Use
// WithColorMode(ColorAlways) if you explicitly want ANSI codes.
func SPrint(payload []byte, opts ...Option) (string, error) {
	var buf bytes.Buffer
	opts = append([]Option{WithWriter(&buf)}, opts...)
	if err := Print(payload, opts...); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Format pretty-prints a JSON payload and returns the formatted bytes.
func Format(payload []byte, opts ...Option) ([]byte, error) {
	s, err := SPrint(payload, opts...)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func parsePayload(payload []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}

	// Reject trailing non-whitespace after the first JSON value.
	if err := dec.Decode(&struct{}{}); err != nil {
		if errors.Is(err, io.EOF) {
			return v, nil
		}
		return nil, fmt.Errorf("invalid JSON payload: trailing data after first JSON value")
	}
	return nil, fmt.Errorf("invalid JSON payload: trailing data after first JSON value")
}

type printer struct {
	w        io.Writer
	indent   string
	sortKeys bool
	color    bool
	style    Style
}

func newPrinter(w io.Writer, o Options) *printer {
	color := false
	switch o.Color {
	case ColorAlways:
		color = true
	case ColorNever:
		color = false
	default:
		color = isTerminalWriter(w)
	}

	return &printer{
		w:        w,
		indent:   o.Indent,
		sortKeys: o.SortKeys,
		color:    color,
		style:    o.Style,
	}
}

func (p *printer) writeString(s string) error {
	_, err := io.WriteString(p.w, s)
	return err
}

func (p *printer) ws(level int) string {
	if level <= 0 {
		return ""
	}
	out := ""
	for i := 0; i < level; i++ {
		out += p.indent
	}
	return out
}

func (p *printer) punct(s string) string {
	if p.color && p.style.Punct != nil {
		return p.style.Punct(s)
	}
	return s
}

func (p *printer) key(s string) string {
	if p.color && p.style.Key != nil {
		return p.style.Key(s)
	}
	return s
}

func (p *printer) str(s string) string {
	if p.color && p.style.String != nil {
		return p.style.String(s)
	}
	return s
}

func (p *printer) num(s string) string {
	if p.color && p.style.Number != nil {
		return p.style.Number(s)
	}
	return s
}

func (p *printer) boolean(s string) string {
	if p.color && p.style.Bool != nil {
		return p.style.Bool(s)
	}
	return s
}

func (p *printer) null(s string) string {
	if p.color && p.style.Null != nil {
		return p.style.Null(s)
	}
	return s
}

func (p *printer) printValue(v any, level int) error {
	switch t := v.(type) {
	case map[string]any:
		return p.printObject(t, level)
	case []any:
		return p.printArray(t, level)
	case string:
		b, _ := json.Marshal(t)
		return p.writeString(p.str(string(b)))
	case json.Number:
		// Use the raw representation when possible.
		return p.writeString(p.num(t.String()))
	case float64:
		// Shouldn't happen with UseNumber(), but keep it robust.
		return p.writeString(p.num(strconv.FormatFloat(t, 'f', -1, 64)))
	case bool:
		if t {
			return p.writeString(p.boolean("true"))
		}
		return p.writeString(p.boolean("false"))
	case nil:
		return p.writeString(p.null("null"))
	default:
		// Last resort: let encoding/json render the value.
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Errorf("unsupported JSON value type %T: %w", v, err)
		}
		return p.writeString(string(b))
	}
}

func (p *printer) printObject(m map[string]any, level int) error {
	if err := p.writeString(p.punct("{")); err != nil {
		return err
	}
	if len(m) == 0 {
		return p.writeString(p.punct("}"))
	}
	if err := p.writeString("\n"); err != nil {
		return err
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if p.sortKeys {
		sort.Strings(keys)
	}

	for i, k := range keys {
		if err := p.writeString(p.ws(level + 1)); err != nil {
			return err
		}

		kb, _ := json.Marshal(k)
		if err := p.writeString(p.key(string(kb))); err != nil {
			return err
		}
		if err := p.writeString(p.punct(": ")); err != nil {
			return err
		}
		if err := p.printValue(m[k], level+1); err != nil {
			return err
		}

		if i < len(keys)-1 {
			if err := p.writeString(p.punct(",")); err != nil {
				return err
			}
		}
		if err := p.writeString("\n"); err != nil {
			return err
		}
	}
	if err := p.writeString(p.ws(level)); err != nil {
		return err
	}
	return p.writeString(p.punct("}"))
}

func (p *printer) printArray(a []any, level int) error {
	if err := p.writeString(p.punct("[")); err != nil {
		return err
	}
	if len(a) == 0 {
		return p.writeString(p.punct("]"))
	}
	if err := p.writeString("\n"); err != nil {
		return err
	}

	for i := range a {
		if err := p.writeString(p.ws(level + 1)); err != nil {
			return err
		}
		if err := p.printValue(a[i], level+1); err != nil {
			return err
		}
		if i < len(a)-1 {
			if err := p.writeString(p.punct(",")); err != nil {
				return err
			}
		}
		if err := p.writeString("\n"); err != nil {
			return err
		}
	}
	if err := p.writeString(p.ws(level)); err != nil {
		return err
	}
	return p.writeString(p.punct("]"))
}
