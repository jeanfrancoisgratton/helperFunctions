// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /misc.go
// Original timestamp: 2024/04/10 15:20

package helperFunctions

import (
	"fmt"

	"github.com/jeanfrancoisgratton/helperFunctions/v2/terminalfx"
)

// NUMBER FORMATTING FUNCTIONS
// ===========================

// This function was originally written in 1993, in C, by my friend Jean-François Gauthier (jief@brebis.dyndns.org)
// I've ported it in C# in 2011. It is still loosely based on J.F.Gauthier's version, somehow; credit is given where credit is due
// This function transforms a multi-digit number in International Notation; 1234567 thus becomes 1,234,567
func SI(nombre interface{}) string {
	var str string
	switch n := nombre.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("%d", n)
	case float32, float64:
		str = fmt.Sprintf("%.f", n)
	default:
		return "Invalid input"
	}

	negative := false
	if str[0] == '-' {
		negative = true
		str = str[1:]
	}

	// Insert commas every three digits from the right
	var formatted string
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			formatted += ","
		}
		formatted += string(digit)
	}

	if negative {
		formatted = "-" + formatted
	}

	return formatted
}

// This function takes a string and returns its reverse
// Thus, "12345" becomes "54321"
func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// Prints the changelog of the calling tool
func ChangeLog(cl string, clear bool) {
	if clear {
		terminalfx.ClearTTY()
	}
	fmt.Print(cl)
}
