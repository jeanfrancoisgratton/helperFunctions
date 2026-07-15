package helperFunctions

import (
	"encoding/base64"
	"testing"
)

func TestEncodeDecodeStringRoundtrip(t *testing.T) {
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
			encoded := EncodeString(c.plaintext, c.passphrase)

			if _, err := base64.StdEncoding.DecodeString(encoded); err != nil {
				t.Fatalf("EncodeString output is not valid base64: %v", err)
			}

			decoded := DecodeString(encoded, c.passphrase)
			if decoded != c.plaintext {
				t.Errorf("DecodeString(EncodeString(%q)) = %q, want %q", c.plaintext, decoded, c.plaintext)
			}
		})
	}
}

func TestEncodeStringUsesRandomIV(t *testing.T) {
	first := EncodeString("same plaintext", "same passphrase")
	second := EncodeString("same plaintext", "same passphrase")

	if first == second {
		t.Fatal("EncodeString produced identical ciphertext for two calls; expected a random IV to vary the output")
	}

	if got := DecodeString(first, "same passphrase"); got != "same plaintext" {
		t.Errorf("first ciphertext decoded to %q, want %q", got, "same plaintext")
	}
	if got := DecodeString(second, "same passphrase"); got != "same plaintext" {
		t.Errorf("second ciphertext decoded to %q, want %q", got, "same plaintext")
	}
}

func TestDecodeStringWrongPassphrase(t *testing.T) {
	plaintext := "sensitive data"
	encoded := EncodeString(plaintext, "correct passphrase")

	got := DecodeString(encoded, "wrong passphrase")
	if got == plaintext {
		t.Error("DecodeString with the wrong passphrase unexpectedly recovered the original plaintext")
	}
}

func TestGetPasswordDebugMode(t *testing.T) {
	var got string
	withStdin(t, "s3cr3t\n", func() {
		got = GetPassword("Password: ", true)
	})
	if got != "s3cr3t" {
		t.Errorf("GetPassword(debugmode=true) = %q, want %q", got, "s3cr3t")
	}
}
