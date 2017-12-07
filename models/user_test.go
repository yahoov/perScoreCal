package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	pb "perScoreCal/perScoreProto/user"
	"testing"

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
	var interestRequest pb.GetInterestRequest

	db, err := gorm.Open("postgres", "host=localhost user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")

	defer db.Close()
	if err != nil {
		t.Errorf("Error in setupdb: %+v", err)
	}
	response, err := user.GetInterests(context.Background(), &interestRequest, db)
	if err != nil {
		t.Errorf("Error", err)
	}
	fmt.Println("response,", response)
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
