// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /misc.go
// Original timestamp: 2024/04/10 15:20

package helperFunctions

import (
	"fmt"
	"strconv"
	"strings"
)

// NUMBER FORMATTING FUNCTIONS
// ===========================

// This function was originally written in 1993, in C, by my friend Jean-François Gauthier
// I've ported it in C# in 2011. It is still loosely based on J.F.Gauthier's version, somehow; credit is given where credit is due
// This function transforms a multi-digit number in International Notation; 1234567 thus becomes 1,234,567

func SI(n any) string {
	var str string

	switch v := n.(type) {
	case string:
		str = v

	case int:
		str = strconv.Itoa(v)

	case int8:
		str = strconv.FormatInt(int64(v), 10)

	case int16:
		str = strconv.FormatInt(int64(v), 10)

	case int32:
		str = strconv.FormatInt(int64(v), 10)

	case int64:
		str = strconv.FormatInt(v, 10)

	case uint:
		str = strconv.FormatUint(uint64(v), 10)

	case uint8:
		str = strconv.FormatUint(uint64(v), 10)

	case uint16:
		str = strconv.FormatUint(uint64(v), 10)

	case uint32:
		str = strconv.FormatUint(uint64(v), 10)

	case uint64:
		str = strconv.FormatUint(v, 10)

	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)

	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)

	default:
		return fmt.Sprintf("%v", n)
	}

	negative := false
	if strings.HasPrefix(str, "-") {
		negative = true
		str = str[1:]
	}

	intPart := str
	fracPart := ""

	if idx := strings.IndexByte(str, '.'); idx != -1 {
		intPart = str[:idx]
		fracPart = str[idx:]
	}

	var b strings.Builder

	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}

	if fracPart != "" {
		b.WriteString(fracPart)
	}

	result := b.String()

	if negative {
		result = "-" + result
	}

	return result
}

// This function takes a string and returns its reverse
// Thus, "12345" becomes "54321"
func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// BytesToUnit converts a number in bytes to its equivalent in MB, GB or TB.
//
// The returned string preserves exactly the requested number of decimal places.
//
// Examples:
//
//	BytesToUnit(1234567890, 'g', 3) => "1.150"
//	BytesToUnit(1234567890, 'g', 2) => "1.15"
func BytesToUnit(bytes uint64, unit rune, decimals int) (string, error) {
	var divisor float64

	switch unit {
	case 'm', 'M':
		divisor = 1 << 20
	case 'g', 'G':
		divisor = 1 << 30
	case 't', 'T':
		divisor = 1 << 40
	default:
		return "", fmt.Errorf("invalid unit '%c'", unit)
	}

	if decimals < 0 {
		return "", fmt.Errorf("decimals cannot be negative")
	}

	value := float64(bytes) / divisor

	return fmt.Sprintf("%.*f", decimals, value), nil
}
