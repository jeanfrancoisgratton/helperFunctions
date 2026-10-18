package helperFunctions

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptFileRoundtrip(t *testing.T) {
	dir := t.TempDir()
	infile := filepath.Join(dir, "plain.txt")
	encfile := filepath.Join(dir, "plain.txt.enc")
	outfile := filepath.Join(dir, "plain.decrypted.txt")

	content := []byte("the quick brown fox jumps over the lazy dog\nwith a second line")
	if err := os.WriteFile(infile, content, 0644); err != nil {
		t.Fatalf("WriteFile(infile): %v", err)
	}

	passphrase := "correct horse battery staple"

	if err := EncryptFile(infile, encfile, passphrase); err != nil {
		t.Fatalf("EncryptFile() error = %v", err)
	}

	encrypted, err := os.ReadFile(encfile)
	if err != nil {
		t.Fatalf("ReadFile(encfile): %v", err)
	}
	if _, err := base64.StdEncoding.DecodeString(string(encrypted)); err != nil {
		t.Fatalf("EncryptFile output is not valid base64: %v", err)
	}

	if err := DecryptFile(encfile, outfile, passphrase); err != nil {
		t.Fatalf("DecryptFile() error = %v", err)
	}

	got, err := os.ReadFile(outfile)
	if err != nil {
		t.Fatalf("ReadFile(outfile): %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("decrypted content = %q, want %q", got, content)
	}

	info, err := os.Stat(outfile)
	if err != nil {
		t.Fatalf("Stat(outfile): %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("decrypted file permissions = %o, want %o", perm, 0600)
	}
}

func TestDecryptFileWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	infile := filepath.Join(dir, "plain.txt")
	encfile := filepath.Join(dir, "plain.txt.enc")
	outfile := filepath.Join(dir, "plain.decrypted.txt")

	content := []byte("top secret content")
	if err := os.WriteFile(infile, content, 0644); err != nil {
		t.Fatalf("WriteFile(infile): %v", err)
	}

	if err := EncryptFile(infile, encfile, "correct passphrase"); err != nil {
		t.Fatalf("EncryptFile() error = %v", err)
	}
	if err := DecryptFile(encfile, outfile, "wrong passphrase"); err == nil {
		t.Fatal("DecryptFile() with the wrong passphrase expected an error, got nil")
	}
}

func TestDecryptFileTooShortCiphertext(t *testing.T) {
	dir := t.TempDir()
	encfile := filepath.Join(dir, "short.enc")
	outfile := filepath.Join(dir, "short.out")

	shortPayload := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	if err := os.WriteFile(encfile, []byte(shortPayload), 0644); err != nil {
		t.Fatalf("WriteFile(encfile): %v", err)
	}

	if err := DecryptFile(encfile, outfile, "any passphrase"); err == nil {
		t.Fatal("DecryptFile() with too-short ciphertext expected an error, got nil")
	}
}

func TestDecryptFileInvalidBase64(t *testing.T) {
	dir := t.TempDir()
	encfile := filepath.Join(dir, "invalid.enc")
	outfile := filepath.Join(dir, "invalid.out")

	if err := os.WriteFile(encfile, []byte("!!! not valid base64 !!!"), 0644); err != nil {
		t.Fatalf("WriteFile(encfile): %v", err)
	}

	if err := DecryptFile(encfile, outfile, "any passphrase"); err == nil {
		t.Fatal("DecryptFile() with invalid base64 content expected an error, got nil")
	}
}

func TestEncryptFileMissingInfile(t *testing.T) {
	dir := t.TempDir()
	err := EncryptFile(filepath.Join(dir, "does-not-exist.txt"), filepath.Join(dir, "out.enc"), "passphrase")
	if err == nil {
		t.Fatal("EncryptFile() with a missing input file expected an error, got nil")
	}
}

func TestEncodeFileAndEncryptFileAreNotInterchangeable(t *testing.T) {
	dir := t.TempDir()
	infile := filepath.Join(dir, "plain.txt")
	encfile := filepath.Join(dir, "plain.txt.enc")
	outfile := filepath.Join(dir, "plain.out.txt")

	if err := os.WriteFile(infile, []byte("some content"), 0644); err != nil {
		t.Fatalf("WriteFile(infile): %v", err)
	}

	passphrase := "shared passphrase"
	if err := EncodeFile(infile, encfile, passphrase); err != nil {
		t.Fatalf("EncodeFile() error = %v", err)
	}

	// Legacy (unauthenticated CFB) output fed into the new, authenticated
	// decryptor must not be silently accepted.
	if err := DecryptFile(encfile, outfile, passphrase); err == nil {
		t.Error("DecryptFile() accepted EncodeFile() output; expected an error since the two schemes are incompatible")
	}
}
