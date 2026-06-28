// helperFunctions
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: /encode_decode_file.go
// Original timestamp: 2026/06/27 23:10:50

package helperFunctions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

// EncodeFile encrypts the contents of infile using AES-256-CFB (same key
// derivation as EncodeString) and writes the base64-encoded result to outfile.
func EncodeFile(infile, outfile, privateKey string) error {
	plaintext, err := os.ReadFile(infile)
	if err != nil {
		return err
	}

	key := sha256sum(privateKey)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return os.WriteFile(outfile, []byte(encoded), 0600)
}

// DecodeFile decrypts a file that was encrypted by EncodeFile and writes the
// recovered plaintext to outfile.
func DecodeFile(infile, outfile, privateKey string) error {
	encoded, err := os.ReadFile(infile)
	if err != nil {
		return err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return err
	}

	key := sha256sum(privateKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return os.ErrInvalid
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return os.WriteFile(outfile, ciphertext, 0600)
}
