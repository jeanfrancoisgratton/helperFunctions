package helperFunctions

import (
	"encoding/base64"
	"testing"
)

func TestEncryptDecryptStringRoundtrip(t *testing.T) {
	cases := []struct {
		name       string
		plaintext  string
		passphrase string
	}{
		{"simple", "hello world", "correct horse battery staple"},
		{"empty plaintext", "", "some passphrase"},
		{"unicode", "héllo, 世界", "pässphräse"},
		{"long passphrase", "s3cr3t", "this is a much longer passphrase than 32 bytes for sure"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			encrypted, err := EncryptString(c.plaintext, c.passphrase)
			if err != nil {
				t.Fatalf("EncryptString() error = %v", err)
			}

			if _, err := base64.StdEncoding.DecodeString(encrypted); err != nil {
				t.Fatalf("EncryptString output is not valid base64: %v", err)
			}

			decrypted, err := DecryptString(encrypted, c.passphrase)
			if err != nil {
				t.Fatalf("DecryptString() error = %v", err)
			}
			if decrypted != c.plaintext {
				t.Errorf("DecryptString(EncryptString(%q)) = %q, want %q", c.plaintext, decrypted, c.plaintext)
			}
		})
	}
}

func TestEncryptStringUsesRandomSaltAndNonce(t *testing.T) {
	first, err := EncryptString("same plaintext", "same passphrase")
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	second, err := EncryptString("same plaintext", "same passphrase")
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	if first == second {
		t.Fatal("EncryptString produced identical ciphertext for two calls; expected a random salt/nonce to vary the output")
	}

	for _, ct := range []string{first, second} {
		got, err := DecryptString(ct, "same passphrase")
		if err != nil {
			t.Fatalf("DecryptString() error = %v", err)
		}
		if got != "same plaintext" {
			t.Errorf("DecryptString(%q) = %q, want %q", ct, got, "same plaintext")
		}
	}
}

func TestDecryptStringWrongPassphrase(t *testing.T) {
	encrypted, err := EncryptString("sensitive data", "correct passphrase")
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	if _, err := DecryptString(encrypted, "wrong passphrase"); err == nil {
		t.Error("DecryptString with the wrong passphrase expected an error, got nil")
	}
}

func TestDecryptStringTamperedCiphertextDetected(t *testing.T) {
	encrypted, err := EncryptString("sensitive data", "correct passphrase")
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	raw[len(raw)-1] ^= 0xFF // flip a bit in the ciphertext/tag
	tampered := base64.StdEncoding.EncodeToString(raw)

	if _, err := DecryptString(tampered, "correct passphrase"); err == nil {
		t.Error("DecryptString with tampered ciphertext expected an error (GCM auth failure), got nil")
	}
}

func TestDecryptStringTooShortCiphertext(t *testing.T) {
	shortPayload := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	if _, err := DecryptString(shortPayload, "any passphrase"); err == nil {
		t.Error("DecryptString() with too-short ciphertext expected an error, got nil")
	}
}

func TestDecryptStringInvalidBase64(t *testing.T) {
	if _, err := DecryptString("!!! not valid base64 !!!", "any passphrase"); err == nil {
		t.Error("DecryptString() with invalid base64 input expected an error, got nil")
	}
}
