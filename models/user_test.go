package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"

	upb "perScoreCal/perScoreProto/user"

	"github.com/jinzhu/gorm"
)

func TestGetEmail(t *testing.T) {
	testEmail := "test@mail.com"
	token := Encrypt(testEmail)
	if email := GetEmail(token); email != testEmail {
		t.Errorf("Expected email to be %s, but it was %s", testEmail, email)
	}
}

func TestGetInterests(t *testing.T) {
	var user User
	var interestRequest upb.GetInterestRequest
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("TEST_HOST"), os.Getenv("TEST_DBNAME"), os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_SSLMODE"))
	db, err := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbString)

	defer db.Close()
	if err != nil {
		t.Errorf("Error in opening DB connection: %+v", err)
	}
	_, err = user.GetInterests(context.Background(), &interestRequest, db)
	if err != nil {
		t.Errorf("Error: %s", err)
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
