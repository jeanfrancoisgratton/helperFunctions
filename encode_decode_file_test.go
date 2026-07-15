package helperFunctions

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeDecodeFileRoundtrip(t *testing.T) {
	dir := t.TempDir()
	infile := filepath.Join(dir, "plain.txt")
	encfile := filepath.Join(dir, "plain.txt.enc")
	outfile := filepath.Join(dir, "plain.decoded.txt")

	content := []byte("the quick brown fox jumps over the lazy dog\nwith a second line")
	if err := os.WriteFile(infile, content, 0644); err != nil {
		t.Fatalf("WriteFile(infile): %v", err)
	}

	passphrase := "correct horse battery staple"

	if err := EncodeFile(infile, encfile, passphrase); err != nil {
		t.Fatalf("EncodeFile() error = %v", err)
	}

	encoded, err := os.ReadFile(encfile)
	if err != nil {
		t.Fatalf("ReadFile(encfile): %v", err)
	}
	if _, err := base64.StdEncoding.DecodeString(string(encoded)); err != nil {
		t.Fatalf("EncodeFile output is not valid base64: %v", err)
	}

	if err := DecodeFile(encfile, outfile, passphrase); err != nil {
		t.Fatalf("DecodeFile() error = %v", err)
	}

	got, err := os.ReadFile(outfile)
	if err != nil {
		t.Fatalf("ReadFile(outfile): %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("decoded content = %q, want %q", got, content)
	}

	info, err := os.Stat(outfile)
	if err != nil {
		t.Fatalf("Stat(outfile): %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("decoded file permissions = %o, want %o", perm, 0600)
	}
}

func TestDecodeFileWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	infile := filepath.Join(dir, "plain.txt")
	encfile := filepath.Join(dir, "plain.txt.enc")
	outfile := filepath.Join(dir, "plain.decoded.txt")

	content := []byte("top secret content")
	if err := os.WriteFile(infile, content, 0644); err != nil {
		t.Fatalf("WriteFile(infile): %v", err)
	}

	if err := EncodeFile(infile, encfile, "correct passphrase"); err != nil {
		t.Fatalf("EncodeFile() error = %v", err)
	}
	if err := DecodeFile(encfile, outfile, "wrong passphrase"); err != nil {
		t.Fatalf("DecodeFile() error = %v", err)
	}

	got, err := os.ReadFile(outfile)
	if err != nil {
		t.Fatalf("ReadFile(outfile): %v", err)
	}
	if string(got) == string(content) {
		t.Error("DecodeFile with the wrong passphrase unexpectedly recovered the original content")
	}
}

func TestDecodeFileTooShortCiphertext(t *testing.T) {
	dir := t.TempDir()
	encfile := filepath.Join(dir, "short.enc")
	outfile := filepath.Join(dir, "short.out")

	// Fewer than aes.BlockSize (16) raw bytes once base64-decoded.
	shortPayload := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	if err := os.WriteFile(encfile, []byte(shortPayload), 0644); err != nil {
		t.Fatalf("WriteFile(encfile): %v", err)
	}

	err := DecodeFile(encfile, outfile, "any passphrase")
	if err == nil {
		t.Fatal("DecodeFile() with too-short ciphertext expected an error, got nil")
	}
}

func TestDecodeFileInvalidBase64(t *testing.T) {
	dir := t.TempDir()
	encfile := filepath.Join(dir, "invalid.enc")
	outfile := filepath.Join(dir, "invalid.out")

	if err := os.WriteFile(encfile, []byte("!!! not valid base64 !!!"), 0644); err != nil {
		t.Fatalf("WriteFile(encfile): %v", err)
	}

	err := DecodeFile(encfile, outfile, "any passphrase")
	if err == nil {
		t.Fatal("DecodeFile() with invalid base64 content expected an error, got nil")
	}
}

func TestEncodeFileMissingInfile(t *testing.T) {
	dir := t.TempDir()
	err := EncodeFile(filepath.Join(dir, "does-not-exist.txt"), filepath.Join(dir, "out.enc"), "passphrase")
	if err == nil {
		t.Fatal("EncodeFile() with a missing input file expected an error, got nil")
	}
}
