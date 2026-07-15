package prettyjson

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintDefaultOptions(t *testing.T) {
	input := []byte(`{"b":1,"a":"hello","c":true,"d":null,"e":[1,2,3],"f":{}}`)

	want := "{\n" +
		"  \"a\": \"hello\",\n" +
		"  \"b\": 1,\n" +
		"  \"c\": true,\n" +
		"  \"d\": null,\n" +
		"  \"e\": [\n" +
		"    1,\n" +
		"    2,\n" +
		"    3\n" +
		"  ],\n" +
		"  \"f\": {}\n" +
		"}\n"

	got, err := SPrint(input, WithColorMode(ColorNever))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}
	if got != want {
		t.Errorf("SPrint() =\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatMatchesSPrint(t *testing.T) {
	input := []byte(`{"x":1}`)

	s, err := SPrint(input, WithColorMode(ColorNever))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}
	b, err := Format(input, WithColorMode(ColorNever))
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if string(b) != s {
		t.Errorf("Format() = %q, want it to match SPrint() = %q", b, s)
	}
}

func TestPrintWritesToConfiguredWriter(t *testing.T) {
	var buf bytes.Buffer
	if err := Print([]byte(`{"a":1}`), WithWriter(&buf), WithColorMode(ColorNever)); err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("Print() did not write anything to the configured writer")
	}
	if !strings.Contains(buf.String(), `"a": 1`) {
		t.Errorf("Print() output = %q, want it to contain %q", buf.String(), `"a": 1`)
	}
}

func TestPrintCustomIndent(t *testing.T) {
	input := []byte(`{"a":{"b":1}}`)
	want := "{\n" +
		"    \"a\": {\n" +
		"        \"b\": 1\n" +
		"    }\n" +
		"}\n"

	got, err := SPrint(input, WithColorMode(ColorNever), WithIndent("    "))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}
	if got != want {
		t.Errorf("SPrint() with custom indent =\n%s\nwant:\n%s", got, want)
	}
}

func TestPrintSortKeysDisabled(t *testing.T) {
	input := []byte(`{"z":1,"a":2}`)
	got, err := SPrint(input, WithColorMode(ColorNever), WithSortKeys(false))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}
	// Map iteration order is unspecified when sorting is disabled, so only
	// assert that both entries made it through, not their relative order.
	if !strings.Contains(got, `"z": 1`) || !strings.Contains(got, `"a": 2`) {
		t.Errorf("SPrint() with SortKeys(false) = %q, want it to contain both keys", got)
	}
}

func TestPrintColorAlwaysAppliesStyle(t *testing.T) {
	wrap := func(tag string) func(string) string {
		return func(s string) string { return tag + "<" + s + ">" }
	}
	style := Style{
		Key:    wrap("K"),
		String: wrap("S"),
		Number: wrap("N"),
		Bool:   wrap("B"),
		Null:   wrap("Z"),
		Punct:  wrap("P"),
	}

	input := []byte(`{"a":1,"b":true,"c":null,"d":"x"}`)
	got, err := SPrint(input, WithColorMode(ColorAlways), WithStyle(style), WithSortKeys(true))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}

	wantSubstrings := []string{
		`K<"a">`, `N<1>`,
		`K<"b">`, `B<true>`,
		`K<"c">`, `Z<null>`,
		`K<"d">`, `S<"x">`,
		`P<{>`, `P<}>`,
	}
	for _, s := range wantSubstrings {
		if !strings.Contains(got, s) {
			t.Errorf("SPrint() with ColorAlways = %q, want it to contain %q", got, s)
		}
	}
}

func TestPrintColorNeverIgnoresStyle(t *testing.T) {
	style := Style{Number: func(s string) string { return "SHOULD-NOT-APPEAR<" + s + ">" }}
	got, err := SPrint([]byte(`{"a":1}`), WithColorMode(ColorNever), WithStyle(style))
	if err != nil {
		t.Fatalf("SPrint() error = %v", err)
	}
	if strings.Contains(got, "SHOULD-NOT-APPEAR") {
		t.Errorf("SPrint() with ColorNever applied the style anyway: %q", got)
	}
}

func TestPrintTrailingWhitespaceIsOK(t *testing.T) {
	if _, err := SPrint([]byte("{\"a\":1}   \n\t"), WithColorMode(ColorNever)); err != nil {
		t.Errorf("SPrint() with trailing whitespace unexpected error: %v", err)
	}
}

func TestPrintErrors(t *testing.T) {
	cases := []struct {
		name    string
		payload string
	}{
		{"malformed JSON", `{"a":}`},
		{"empty payload", ``},
		{"trailing garbage after value", `{"a":1} garbage`},
		{"trailing second JSON value", `{"a":1}{"b":2}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := SPrint([]byte(c.payload), WithColorMode(ColorNever)); err == nil {
				t.Errorf("SPrint(%q) expected an error, got nil", c.payload)
			}
		})
	}
}

func TestIsTerminalWriter(t *testing.T) {
	var buf bytes.Buffer
	if isTerminalWriter(&buf) {
		t.Error("isTerminalWriter(*bytes.Buffer) = true, want false")
	}

	f, err := os.CreateTemp(t.TempDir(), "prettyjson-*")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()
	if isTerminalWriter(f) {
		t.Error("isTerminalWriter(regular *os.File) = true, want false")
	}
}

func TestStyles(t *testing.T) {
	plain := PlainStyle()
	if plain.Key != nil || plain.String != nil || plain.Number != nil ||
		plain.Bool != nil || plain.Null != nil || plain.Punct != nil {
		t.Errorf("PlainStyle() = %+v, want all-nil funcs", plain)
	}

	def := DefaultStyle()
	if def.Key == nil || def.String == nil || def.Number == nil ||
		def.Bool == nil || def.Null == nil || def.Punct == nil {
		t.Errorf("DefaultStyle() = %+v, want all funcs set", def)
	}
	// Sanity: none of the default colorizers should panic and each should
	// preserve the original text somewhere in its output.
	if !strings.Contains(def.Key("hi"), "hi") {
		t.Error("DefaultStyle().Key(\"hi\") lost the original text")
	}
}
