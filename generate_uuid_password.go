// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /generate_uuid_password.go
// Original timestamp : 2026.11.07 10:14:55

package helperFunctions

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	hexChars     = "0123456789abcdef"
	specialChars = "!@#$%^&*()-_=+[]{};:,.<>?/|~"
)

// GenerateUUIDorPassword returns a `length`-characters string. It takes 3 params:
// - the exact length of the string to be returned
// - hexOnly: only hexadecimal chars to be returned; perfect to generate UUIDs
// - allowSpecialChars: if true, special characters are include; perfect to generate passwords
// - both booleans cannot be true simultaneously; if it happens both are assumed to be false

func GenerateUUIDorPassword(length int, hexOnly, allowSpecialChars bool) (string, error) {
	if length <= 0 {
		return "", nil
	}

	charset := alphaNumeric

	// Both true means both options are disabled.
	if hexOnly && !allowSpecialChars {
		charset = hexChars
	} else if allowSpecialChars && !hexOnly {
		charset += specialChars
	}

	var result strings.Builder
	result.Grow(length)

	maxval := big.NewInt(int64(len(charset)))

	for range length {
		n, err := rand.Int(rand.Reader, maxval)
		if err != nil {
			return "", err
		}

		result.WriteByte(charset[n.Int64()])
	}

	return result.String(), nil
}
