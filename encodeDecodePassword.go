// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /encodeDecodePassword.go
// Original timestamp: 2024/04/10 15:03

package helperFunctions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Breaking change: if DebugMode is true, we catch the passwd in cleartext
func GetPassword(prompt string, debugmode bool) string {
	if debugmode {
		return GetStringValFromPrompt(prompt)
	}
	// Get the initial state of the terminal.
	initialTermState, e1 := terminal.GetState(syscall.Stdin)
	if e1 != nil {
		panic(e1)
	}

	// Restore it in the event of an interrupt.
	// CITATION: Konstantin Shaposhnikov - https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		_ = terminal.Restore(syscall.Stdin, initialTermState)
		os.Exit(1)
	}()

	// Now get the password.
	fmt.Print(prompt)
	p, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		panic(err)
	}

	// Stop looking for ^C on the channel.
	signal.Stop(c)

	// Return the password as a string.
	return string(p)
}

// Quick functions to encode and decode strings
// This is based on my encryption-decryption tool : https://github.com/jeanfrancoisgratton/encdec
func EncodeString(string2encode string, privateKey string) string {
	//func EncodeString(string2encode string) string {
	// privateKey is optional here
	if len(privateKey) != 32 {
		// yeah, I say "*crypt" instead of "*code", but I needed 32bits...
		privateKey = "secret key 2 encrypt and decrypt"
	}

	key := []byte(privateKey)
	plaintext := []byte(string2encode)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func DecodeString(encodedstring string, privateKey string) string {
	// privateKey is optional here
	if len(privateKey) != 32 {
		// yeah, I say "*crypt" instead of "*code", but I needed 32bits...
		privateKey = "secret key 2 encrypt and decrypt"
	}

	key := []byte(privateKey)
	ciphertext, _ := base64.URLEncoding.DecodeString(encodedstring)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
