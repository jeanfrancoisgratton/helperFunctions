package helperFunctions

import "testing"

func TestSI(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"string small", "42", "42"},
		{"string grouped", "1234567", "1,234,567"},
		{"int grouped", 1234567, "1,234,567"},
		{"int negative", -1234567, "-1,234,567"},
		{"int8", int8(120), "120"},
		{"int16", int16(12345), "12,345"},
		{"int32", int32(1234567), "1,234,567"},
		{"int64", int64(1234567890), "1,234,567,890"},
		{"uint", uint(1234567), "1,234,567"},
		{"uint8", uint8(250), "250"},
		{"uint16", uint16(12345), "12,345"},
		{"uint32", uint32(1234567), "1,234,567"},
		{"uint64", uint64(1234567890), "1,234,567,890"},
		{"float32", float32(1234), "1,234"},
		{"float64 with fraction", 1234567.5, "1,234,567.5"},
		{"below grouping threshold", 999, "999"},
		{"exact grouping boundary", 1000, "1,000"},
		{"zero", 0, "0"},
		{"unsupported type falls back to %v", true, "true"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SI(c.in)
			if got != c.want {
				t.Errorf("SI(%v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"digits", "12345", "54321"},
		{"empty", "", ""},
		{"single rune", "a", "a"},
		{"palindrome", "abcba", "abcba"},
		{"multi-byte runes", "héllo", "olléh"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ReverseString(c.in)
			if got != c.want {
				t.Errorf("ReverseString(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestBytesToUnit(t *testing.T) {
	cases := []struct {
		name     string
		bytes    uint64
		unit     rune
		decimals int
		want     string
		wantErr  bool
	}{
		{"gigabytes 3 decimals", 1234567890, 'g', 3, "1.150", false},
		{"gigabytes 2 decimals", 1234567890, 'g', 2, "1.15", false},
		{"megabytes uppercase", 1 << 20, 'M', 0, "1", false},
		{"terabytes", 1 << 40, 't', 1, "1.0", false},
		{"zero bytes", 0, 'g', 2, "0.00", false},
		{"invalid unit", 100, 'x', 2, "", true},
		{"negative decimals", 100, 'g', -1, "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := BytesToUnit(c.bytes, c.unit, c.decimals)
			if c.wantErr {
				if err == nil {
					t.Fatalf("BytesToUnit(%d, %q, %d) expected error, got nil", c.bytes, c.unit, c.decimals)
				}
				return
			}
			if err != nil {
				t.Fatalf("BytesToUnit(%d, %q, %d) unexpected error: %v", c.bytes, c.unit, c.decimals, err)
			}
			if got != c.want {
				t.Errorf("BytesToUnit(%d, %q, %d) = %q, want %q", c.bytes, c.unit, c.decimals, got, c.want)
			}
		})
	}
}
