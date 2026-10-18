// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /encrypt_decrypt_string.go
// Original timestamp: 2026/10/18

package helperFunctions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters for key derivation, matching the parameters used in
// golang.org/x/crypto/argon2's own reference example: 1 pass, 64 MiB of
// memory, 4 parallel lanes, a 32-byte key (AES-256).
const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
	saltSize      = 16
)

// deriveKeyArgon2id derives an AES-256 key from passphrase and salt using
// Argon2id (RFC 9106). Unlike sha256sum (used by EncodeString/EncodeFile),
// this is a deliberately slow, memory-hard KDF: it resists GPU/ASIC
// brute-forcing of the passphrase, and a random per-call salt means the
// same passphrase never derives the same key twice.
func deriveKeyArgon2id(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
}

// EncryptString encrypts string2encrypt with AES-256-GCM, deriving the key
// from passphrase via Argon2id with a random salt. The result is
// base64-encoded salt || nonce || ciphertext (ciphertext includes the GCM
// authentication tag), safe to store or transmit as plain text.
//
// This is the authenticated counterpart to EncodeString: tampering with
// the returned string is detected on decrypt (DecryptString returns an
// error) rather than silently producing corrupted plaintext. The two are
// not interchangeable — data encrypted with one can't be decrypted with
// the other.
func EncryptString(string2encrypt string, passphrase string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key := deriveKeyArgon2id(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(string2encrypt), nil)

	out := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return base64.StdEncoding.EncodeToString(out), nil
}

// DecryptString decrypts a string produced by EncryptString. It returns an
// error — rather than panicking, unlike DecodeString — when the passphrase
// is wrong or the data has been tampered with; AES-GCM's authentication
// tag makes both cases detectable instead of yielding corrupted output.
func DecryptString(encrypted string, passphrase string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	if len(raw) < saltSize {
		return "", errors.New("helperFunctions: ciphertext too short")
	}
	salt, raw := raw[:saltSize], raw[saltSize:]

	key := deriveKeyArgon2id(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("helperFunctions: ciphertext too short")
	}
	nonce, raw := raw[:nonceSize], raw[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, raw, nil)
	if err != nil {
		return "", errors.New("helperFunctions: decryption failed (wrong passphrase or corrupted data)")
	}
	return string(plaintext), nil
}
