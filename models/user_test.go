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

const Key = "fkzfgk0FY2CaYJhyXbshnPJaRrFtCwfj"

func TestGetEmail(t *testing.T) {
	testEmail := "test@mail.com"
	token := Encrypt(testEmail)
	if email := models.GetEmail(token); email != testEmail {
		t.Errorf("Expected email to be %s, but it was %s", testEmail, email)
	}
}

func Encrypt(text string) string {
	key := []byte(Key)
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
