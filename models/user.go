package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	upb "perScoreCal/perScoreProto/user"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

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

// |
// |
// |

// CreateUser ...
func CreateUser(email string, db *gorm.DB) (User, error) {
	validate := validator.New()
	user := User{Email: email, Score: 0.0}
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
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

// GetInterests ...
func (user User) GetInterests(ctx context.Context, in *upb.GetInterestRequest, db *gorm.DB) (*upb.GetInterestResponse, error) {
	var response = new(upb.GetInterestResponse)
	var categories []Category
	err := db.Find(&categories).Error

	if err != nil {
		response.Success = false
		response.Message = "Failed to retrieve categories"
		log.Errorf("failed to retrieve categories: %v", err)
	} else {
		response.Title = "Please select your interest to proceed"
		response.Body = ""

		for index, category := range categories {
			response.Categories = append(response.Categories, new(upb.GetInterestResponse_Category))
			response.Categories[index].Id = int32(category.ID)
			response.Categories[index].Name = category.Name
			response.Categories[index].Parent = int32(category.Parent)
			response.Categories[index].Level = category.Level
			response.Categories[index].WeightRange = GetWeightRange(&category, db)
		}
	}

	return response, err
}

// |
// |
// |

// GetEntries ...
func (user User) GetEntries(ctx context.Context, in *upb.GetEntriesRequest, db *gorm.DB) (*upb.GetEntriesResponse, error) {
	var response = new(upb.GetEntriesResponse)
	var categories []Category
	err := db.Where("approved = ?", false).Find(&categories).Error

	if err != nil {
		response.Success = false
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

	var questions []Question
	err = db.Where("approved = ?", false).Find(&questions).Error

	if err != nil {
		response.Success = false
		response.Message = "Failed to retrieve questions"
		log.Errorf("failed to retrieve questions: %v", err)
	} else {
		for index, question := range questions {
			response.Questions = append(response.Questions, new(upb.GetEntriesResponse_Question))
			response.Questions[index].Id = int32(question.ID)
			response.Questions[index].Title = question.Title
			response.Questions[index].Body = question.Body
		}
	}

	return response, err
}

// |
// |
// |

// ApproveEntries ...
func (user User) ApproveEntries(ctx context.Context, in *upb.ApproveEntriesRequest, db *gorm.DB) (*upb.ApproveEntriesResponse, error) {
	var err error
	var response = new(upb.ApproveEntriesResponse)

	for _, category := range in.Categories {
		if category.Approved == true {
			var dbCategory Category
			err = db.Find(&dbCategory, uint(category.Id)).Error
			if err != nil {
				response.Success = false
				response.Message = "Failed to retrieve category with ID: " + strconv.Itoa(int(category.Id))
				log.Errorf("failed to retrieve category: %v", err)
			} else {
				dbCategory.Approved = true
				err = db.Save(&dbCategory).Error
				if err != nil {
					response.Success = false
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
				response.Success = false
				response.Message = "Failed to retrieve question with ID: " + strconv.Itoa(int(question.Id))
				log.Errorf("failed to retrieve question: %v", err)
			} else {
				dbQuestion.Approved = true
				err = db.Save(&dbQuestion).Error
				if err != nil {
					response.Success = false
					response.Message = "Failed to save question with ID: " + strconv.Itoa(int(question.Id))
					log.Errorf("failed to save question: %v", err)
				}
			}
		}
	}

	return response, err
}
