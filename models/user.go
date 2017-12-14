package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	upb "perScoreCal/perScoreProto/user"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	validator "gopkg.in/go-playground/validator.v9"
)

const Key = "fkzfgk0FY2CaYJhyXbshnPJaRrFtCwfj"

// User is a gorm model
type User struct {
	gorm.Model
	Email       string `validate:"required"`
	Questions   []Question
	Answers     []Answer
	UsersAnswer UsersAnswer
	Score       float32
}

// |
// |
// |

// GetEmail from authToken
func GetEmail(authToken string) string {
	mappedResult := Decrypt(authToken)
	return mappedResult["email"]
}

// |
// |
// |

// Decrypt ...
func Decrypt(cryptoText string) map[string]string {
	mappedResult := make(map[string]string)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	key := []byte(Key)
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

	dataArray := strings.Split(fmt.Sprintf("%s", dataByte), ",")

	mappedResult["email"] = dataArray[0]
	mappedResult["role"] = dataArray[1]
	mappedResult["sessionTime"] = dataArray[2]

	return mappedResult
}

// |
// |
// |

// CreateUser ...
func CreateUser(email string, db *gorm.DB) (User, error) {
	validate := validator.New()
	user := User{Email: email, Score: 0.0}
	err := validate.Struct(user)
	if err != nil {
		for _, errV := range err.(validator.ValidationErrors) {
			fmt.Println("*** Validation Errors ***")
			fmt.Println("NAMESPACE:", errV.Namespace())
			fmt.Println("FIELD:", errV.Field())
			fmt.Println("TAG:", errV.Tag())
			fmt.Println("TYPE:", errV.Type())
			fmt.Println("VALUE:", errV.Value())
			fmt.Println("PARAM:", errV.Param())
			fmt.Println()
		}

		return user, err
	}
	err = db.Create(&user).Error
	return user, err
}

// |
// |
// |

// GetPersonalityScore ...
func GetPersonalityScore(user User, answer Answer, option int32, db *gorm.DB) (float32, error) {
	var weight Weight
	db.Find(&weight, uint(answer.Weights[option]))
	user.Score += float32(weight.Value)
	err := db.Save(&user).Error
	return user.Score, err
}

// |
// |
// |

// GetEntries ...
func (user User) GetEntries(ctx context.Context, in *upb.GetEntriesRequest, db *gorm.DB) (*upb.GetEntriesResponse, error) {
	var response = new(upb.GetEntriesResponse)
	var categories []Category
	var err error

	if in.AuthToken == "" {
		response.Status = "FAILURE"
		response.Message = "Invalid request"
		log.Errorf("No AuthToken received")
		return response, errors.New(response.Message)
	}

	response.Status = "SUCCESS"
	response.Message = "You are in!"

	mappedResult := Decrypt(in.AuthToken)
	response.Role = mappedResult["role"]

	if mappedResult["role"] == "Administrator" {
		err = db.Where("approved = ?", false).Find(&categories).Error
		var questions []Question
		err = db.Where("approved = ?", false).Find(&questions).Error
		fmt.Println("admin ques", questions)
		fmt.Println("admin categories", categories)
		if err != nil {
			response.Status = "FAILURE"
			response.Message = "Failed to retrieve questions"
			log.Errorf("failed to retrieve questions: %v", err)
		} else {
			for index, question := range questions {
				response.Questions = append(response.Questions, new(upb.GetEntriesResponse_Question))
				var answer Answer
				result := db.Where("question_id = ?", question.ID).Find(&answer).RecordNotFound()
				if result == false {
					responseAnswer := new(upb.GetEntriesResponse_Question_Answer)
					responseAnswer.Option1 = answer.Option1
					responseAnswer.Option2 = answer.Option2
					responseAnswer.Option3 = answer.Option3
					responseAnswer.Option4 = answer.Option4
					responseAnswer.Option5 = answer.Option5
					response.Questions[index].Id = int32(question.ID)
					response.Questions[index].Title = question.Title
					response.Questions[index].Body = question.Body
					response.Questions[index].Answer = responseAnswer
				} else {
					log.Errorf("failed to retrieve answer: %v", err)
				}
			}
		}
	} else {
		// For Questioner and Respondent
		err = db.Where("approved = ?", true).Find(&categories).Error
	}

	if err != nil {
		response.Status = "FAILURE"
		response.Message = "Failed to retrieve categories"
		log.Errorf("failed to retrieve categories: %v", err)
	} else {
		for index, category := range categories {
			response.Categories = append(response.Categories, new(upb.GetEntriesResponse_Category))
			response.Categories[index].Id = int32(category.ID)
			response.Categories[index].Name = category.Name
			response.Categories[index].Parent = int32(category.Parent)
			response.Categories[index].Level = category.Level
			response.Categories[index].WeightRange = GetWeightRange(&category, db)
		}
	}
	fmt.Println("response get entry", response)
	return response, err
}

// |
// |
// |

// ApproveEntries ...
func (user User) ApproveEntries(ctx context.Context, in *upb.ApproveEntriesRequest, db *gorm.DB) (*upb.ApproveEntriesResponse, error) {
	var err error
	var response = new(upb.ApproveEntriesResponse)
	response.Status = "SUCCESS"
	response.Message = "Successfully retrived entries"

	for _, category := range in.Categories {
		if category.Approved == true {
			var dbCategory Category
			err = db.Find(&dbCategory, uint(category.Id)).Error
			if err != nil {
				response.Status = "FAILURE"
				response.Message = "Failed to retrieve category with ID: " + strconv.Itoa(int(category.Id))
				log.Errorf("failed to retrieve category: %v", err)
			} else {
				dbCategory.Approved = true
				err = db.Save(&dbCategory).Error
				if err != nil {
					response.Status = "FAILURE"
					response.Message = "Failed to save category with ID: " + strconv.Itoa(int(category.Id))
					log.Errorf("failed to save category: %v", err)
				}
			}
		}
	}

	for _, question := range in.Questions {
		if question.Approved == true {
			var dbQuestion Question
			err = db.Find(&dbQuestion, uint(question.Id)).Error
			if err != nil {
				response.Status = "FAILURE"
				response.Message = "Failed to retrieve question with ID: " + strconv.Itoa(int(question.Id))
				log.Errorf("failed to retrieve question: %v", err)
			} else {
				dbQuestion.Approved = true
				err = db.Save(&dbQuestion).Error
				if err != nil {
					response.Status = "FAILURE"
					response.Message = "Failed to save question with ID: " + strconv.Itoa(int(question.Id))
					log.Errorf("failed to save question: %v", err)
				}
			}
		}
	}

	var categories []Category
	err = db.Where("approved = ?", false).Find(&categories).Error

	if err != nil {
		response.Status = "FAILURE"
		response.Message = "Failed to retrieve categories"
		log.Errorf("failed to retrieve categories: %v", err)
	} else {
		for index, category := range categories {
			response.Categories = append(response.Categories, new(upb.ApproveEntriesResponse_Category))
			response.Categories[index].Id = int32(category.ID)
			response.Categories[index].Name = category.Name
			response.Categories[index].Parent = int32(category.Parent)
			response.Categories[index].Level = category.Level
			response.Categories[index].WeightRange = GetWeightRange(&category, db)
		}
	}

	var questions []Question
	err = db.Where("approved = ?", false).Find(&questions).Error

	if err != nil {
		response.Status = "FAILURE"
		response.Message = "Failed to retrieve questions"
		log.Errorf("failed to retrieve questions: %v", err)
	} else {
		for index, question := range questions {
			response.Questions = append(response.Questions, new(upb.ApproveEntriesResponse_Question))
			response.Questions[index].Id = int32(question.ID)
			response.Questions[index].Title = question.Title
			response.Questions[index].Body = question.Body
		}
	}

	return response, err
}
