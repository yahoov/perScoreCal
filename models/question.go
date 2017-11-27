package models

import (
	"context"
	"errors"
	"fmt"
	pb "perScoreCal/perScoreProto/question"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect for gorm
)

// Question is a gorm model
type Question struct {
	gorm.Model
	Title      string
	Body       string
	Answer     Answer
	Weight     Weight
	Approved   bool
	CategoryID uint
	UserID     uint
}

const key = "fkzfgk0FY2CaYJhyXbshnPJaRrFtCwfj"

var categoriesFailed []string

// |
// |
// |

// CreateInDB question and return response
func (question Question) CreateInDB(ctx context.Context, in *pb.CreateQuestionRequest, db *gorm.DB) (*pb.CreateQuestionResponse, error) {
	var err error
	var user User
	var response = new(pb.CreateQuestionResponse)
	email := GetEmail(in.AuthToken, []byte(key))
	if email == "" {
		response.Success = false
		response.Message = "Failed to retrieve email"
		return response, errors.New("Failed to retrieve email")
	}
	err = db.Where("Email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("No user found with email: ", email)
		fmt.Println("Creating new user with email: ", email, " ...")
		user = User{Email: email, Score: 0.0}
		err = db.Create(&user).Error
		if err != nil {
			response.Success = false
			response.Message = "Failed to create new user"
			log.Errorf("failed to create answer: %v", err)
		}

		if in.Title != "" {
			question, err = createQuestion(ctx, in, db, user, question)
			if err != nil {
				response.Success = false
				response.Message = "Failed to create question"
				log.Errorf("failed to create question: %v", err)
			} else {
				_, err = createAnswer(ctx, in, db, question, user)
				if err != nil {
					response.Success = false
					response.Message = "Failed to create answer"
					log.Errorf("failed to create answer: %v", err)
				} else {
					categoriesFailed = categoriesFailed[:0]
					createQuestionCategory(ctx, in.Categories, db, question)
					if len(categoriesFailed) > 0 {
						// fmt.Println("RESPONSE:", response)
						// fmt.Printf("RESPONSE TYPE: %T", response)
						response.Success = false
						response.Message = "Failed to create following categories: " + strings.Join(categoriesFailed[:], ", ")
					}
				}
			}

			if err == nil {
				qID := strconv.FormatUint(uint64(question.ID), 10)
				response.Success = true
				response.Message = "Successfully created question with ID: " + qID
			}
		}
	}

	var categories []Category
	db.Find(&categories)

	for index, category := range categories {
		// var responseCategory *pb.CreateQuestionResponse_Category
		// response.Categories = make([]*pb.CreateQuestionResponse_Category, index)
		response.Categories = append(response.Categories, new(pb.CreateQuestionResponse_Category))
		// fmt.Println("Index:", index)
		// fmt.Println(response.Categories)
		// response.Categories[index] = responseCategory
		response.Categories[index].Name = category.Name
		response.Categories[index].Parent = int32(category.Parent)
		response.Categories[index].Level = category.Level
		response.Categories[index].WeightRange = GetWeightRange(&category, db)
	}

	return response, err
}

// |
// |
// |

// GetFromDB question in response to the previous question
func (question Question) GetFromDB(ctx context.Context, in *pb.GetQuestionRequest, db *gorm.DB) (*pb.GetQuestionResponse, error) {
	var response *pb.GetQuestionResponse
	email := GetEmail(in.AuthToken, []byte(key))
	db.First(&question, in.QuestionId)
	var category Category
	db.First(&category, question.CategoryID)
	var user User
	db.Where("email = ?", email).First(&user)
	var answer Answer
	db.Where("QuestionID = ?", question.ID).First(&answer)
	option, err := RegisterAnswer(answer, user, in, db)
	if err != nil {
		response.Success = false
		response.Message = "Failed to register answer of question: " + question.Title
	} else {
		nextQuestion, err := getNextQuestion(question, category, db)
		if err != nil {
			response.Success = false
			response.Message = "Failed to get next question for: " + question.Title
		} else {
			response.Success = true
			response.Message = "Successfully retreived next question"
			response.Title = nextQuestion.Title
			response.Body = nextQuestion.Body
			response.Answer.Option1 = answer.Option1
			response.Answer.Option2 = answer.Option2
			response.Answer.Option3 = answer.Option3
			response.Answer.Option4 = answer.Option4
			response.Answer.Option5 = answer.Option5
			score, err := GetPersonalityScore(user, answer, option, db)
			if err != nil {
				response.Success = false
				response.Message = "Failed to get personality score for: " + user.Email
			} else {
				response.Score = score
			}
		}
	}

	return response, nil
}

// |
// |
// |

func getNextQuestion(question Question, category Category, db *gorm.DB) (Question, error) {
	var nextQuestion Question
	weight := question.Weight.Value
	var questions []Question
	err := db.Model(category).Order("Weight asc").Association("Questions").Find(&questions).Error
	if err != nil {
		log.Errorf("failed to get question: %v", err)
	} else {
		for _, dbQuestion := range questions {
			if question.Title == dbQuestion.Title {
				continue
			}
			if dbQuestion.Weight.Value >= weight {
				nextQuestion = dbQuestion
				break
			}
		}
	}
	return nextQuestion, err
}

// |
// |
// |

func createWeight(ctx context.Context, in *pb.CreateQuestionRequest, db *gorm.DB, answer Answer) ([5]byte, error) {
	var answerWeights [5]byte
	var err error
	for _, weight := range in.Answer.Weights {
		var createdWeight Weight
		createdWeight.AnswerID = answer.ID
		createdWeight.Value = weight.Value
		createdWeight.Option = weight.Option
		err = db.Create(&createdWeight).Error
		if err != nil {
			log.Errorf("failed to create weightage: %v", err)
		}
		fmt.Println("answer weight option:", weight.Option)
		answerWeights[weight.Option-1] = byte(createdWeight.ID)
	}

	return answerWeights, err
}

// |
// |
// |

func createAnswer(ctx context.Context, in *pb.CreateQuestionRequest, db *gorm.DB, question Question, user User) (Answer, error) {
	var err error
	answer := Answer{
		UserID:     user.ID,
		QuestionID: question.ID,
		Option1:    in.Answer.Option1,
		Option2:    in.Answer.Option2,
		Option3:    in.Answer.Option3,
		Option4:    in.Answer.Option4,
		Option5:    in.Answer.Option5,
	}

	categoriesFailed = categoriesFailed[:0]
	createAnswerCategory(ctx, in.Answer.Categories, db, answer)

	if len(categoriesFailed) > 0 {
		err = errors.New("Failed to create answer categories")
		log.Errorf("failed to create following answer categories: %v", categoriesFailed)
	} else {
		err = db.Create(&answer).Error
		if err != nil {
			log.Errorf("failed to create answer: %v", err)
		} else {
			var answerWeights [5]byte
			answerWeights, err = createWeight(ctx, in, db, answer)
			if err != nil {
				log.Errorf("failed to create answer weightage: %v", err)
			} else {
				for _, aw := range answerWeights {
					answer.Weights = append(answer.Weights, byte(aw))
				}
			}
		}
	}

	return answer, err
}

// |
// |
// |

func createQuestion(ctx context.Context, in *pb.CreateQuestionRequest, db *gorm.DB, user User, question Question) (Question, error) {
	question.UserID = user.ID
	question.Title = in.Title
	question.Body = in.Body
	question.Approved = false
	err := db.Create(&question).Error

	if err != nil {
		log.Errorf("failed to create question: %v", err)
	} else {
		var questionWeight = Weight{
			QuestionID: question.ID,
			Value:      in.Weight.Value,
		}
		err = db.Create(&questionWeight).Error
		if err != nil {
			log.Errorf("failed to create question weightage: %v", err)
		}
	}
	return question, err
}

// |
// |
// |

func createQuestionCategory(ctx context.Context, requestCategories []*pb.CreateQuestionRequest_Category, db *gorm.DB, question Question) Category {
	var category Category
	for _, requestCategory := range requestCategories {
		if len(requestCategory.Categories) > 0 {
			createQuestionCategory(ctx, requestCategory.Categories, db, question)
		}
		category.Name = requestCategory.Name
		category.Parent = uint(requestCategory.Parent)
		category.Level = GetLevel(&category, db)
		category.Approved = false
		err := db.Create(&category).Error
		if err != nil {
			log.Errorf("failed to create category: %v", err)
			categoriesFailed = append(categoriesFailed, requestCategory.Name)
		} else {
			question.CategoryID = uint(category.ID)
			err = db.Save(&question).Error
			if err != nil {
				log.Errorf("failed to save question: %v", err)
			}
		}
	}
	return category
}

// |
// |
// |

func createAnswerCategory(ctx context.Context, requestCategories []*pb.CreateQuestionRequest_Answer_Category, db *gorm.DB, answer Answer) Category {
	var category Category
	for _, requestCategory := range requestCategories {
		if len(requestCategory.Categories) > 0 {
			createAnswerCategory(ctx, requestCategory.Categories, db, answer)
		}
		category.Name = requestCategory.Name
		category.Parent = uint(requestCategory.Parent)
		category.Level = GetLevel(&category, db)
		category.Approved = false
		err := db.Create(&category).Error
		if err != nil {
			log.Errorf("failed to create category: %v", err)
			categoriesFailed = append(categoriesFailed, requestCategory.Name)
		} else {
			if requestCategory.Option != 0 {
				answer.Categories = make([]byte, requestCategory.Option)
				fmt.Println("Option:", requestCategory.Option-1)
				fmt.Println(answer.Categories)
				answer.Categories[requestCategory.Option-1] = byte(category.ID)
				err = db.Save(&answer).Error
				if err != nil {
					log.Errorf("failed to save answer: %v", err)
				}
			}
		}
	}
	return category
}
