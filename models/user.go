package models

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// User is a gorm model
type User struct {
	gorm.Model
	Email       string
	Questions   []Question
	Answers     []Answer
	UsersAnswer UsersAnswer
	Score       float32
}

// GetEmail from authToken
func GetEmail(authToken string, key []byte) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(authToken)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	dataByte := ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	data := fmt.Sprintf("%s", dataByte)

	dataArray := strings.Split(fmt.Sprintf("%s", data), ",")

	email := dataArray[0]
	// sessionTime := dataArray[1]

	fmt.Println("email: ", email)

	return email
}

// GetPersonalityScore ...
func GetPersonalityScore(user User, answer Answer, option int32, db *gorm.DB) (float32, error) {
	var weight Weight
	db.Find(&weight, uint(answer.Weights[option]))
	user.Score += float32(weight.Value)
	err := db.Save(&user).Error
	return user.Score, err
}
