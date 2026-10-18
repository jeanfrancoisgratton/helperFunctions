// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /encrypt_decrypt_file.go
// Original timestamp: 2026/10/18

package helperFunctions

import "os"

// EncryptFile encrypts the contents of infile using EncryptString
// (AES-256-GCM, Argon2id-derived key) and writes the result to outfile.
func EncryptFile(infile, outfile, passphrase string) error {
	plaintext, err := os.ReadFile(infile)
	if err != nil {
		return err
	}

	encrypted, err := EncryptString(string(plaintext), passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(outfile, []byte(encrypted), 0600)
}

// DecryptFile decrypts a file produced by EncryptFile and writes the
// recovered plaintext to outfile.
func DecryptFile(infile, outfile, passphrase string) error {
	encrypted, err := os.ReadFile(infile)
	if err != nil {
		return err
	}

	plaintext, err := DecryptString(string(encrypted), passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(outfile, []byte(plaintext), 0600)
}
