package models_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"perScoreCal/models"
	"testing"
)

func TestGetEmail(t *testing.T) {
	email := "test@mail.com"
	role := "Administrator"
	sessionInMinutes := "10"
	text := email + "," + role + "," + sessionInMinutes
	token := Encrypt(text)
	if result := models.GetEmail(token); email != result {
		t.Errorf("Expected email to be %s, but it was %s", email, result)
	}
}

func Encrypt(text string) string {
	key := []byte(models.Key)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}
