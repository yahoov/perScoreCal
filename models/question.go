package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	qpb "perScoreCal/perScoreProto/question"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
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

var categories []Category
var categoriesFailed []string

// |
// |
// |

// CreateInDB question and return response
func (question Question) CreateInDB(ctx context.Context, in *qpb.CreateQuestionRequest, db *gorm.DB) (*qpb.CreateQuestionResponse, error) {
	var err error
	var user User
	var response = new(qpb.CreateQuestionResponse)

	if in.AuthToken == "" {
		response.Status = "FAILURE"
		response.Message = "Invalid request"
		log.Errorf("No AuthToken received")
		return response, errors.New(response.Message)
	}

	response.Status = "SUCCESS"
	response.Message = "Question was created successully!"

	email := GetEmail(in.AuthToken)
	if email == "" {
		response.Status = "FAILURE"
		response.Message = "Failed to retrieve email"
		log.Errorf("failed to retrieve email from token: %v", err)
		return response, errors.New("Failed to retrieve email")
	}
	result := db.Where("email = ?", email).First(&user).RecordNotFound()

	if result != false {
		fmt.Println("No user found with email: ", email)
		fmt.Println("Creating new user with email: ", email, " ...")
		user, err = CreateUser(email, db)
		if err != nil {
			response.Status = "FAILURE"
			response.Message = "Failed to create new user. " + fmt.Sprintf("Error: %s", err)
			log.Errorf("failed to create user: %v", err)
		}
	}

	if in.Title != "" {
		question, err = createQuestion(ctx, in, db, user, question)
		if err != nil {
			response.Status = "FAILURE"
			response.Message = "Failed to create question"
			log.Errorf("failed to create question: %v", err)
		} else {
			var answer Answer
			answer, err = createAnswer(ctx, in, db, question, user)
			if err == nil {
				err = db.Save(&answer).Error
			}
			if err != nil {
				response.Status = "FAILURE"
				response.Message = "Failed to create answer"
				log.Errorf("failed to create answer: %v", err)
			} else {
				categoriesFailed = categoriesFailed[:0]
				createQuestionCategories(ctx, in.Category, db, question)
				if len(categoriesFailed) > 0 {
					response.Status = "FAILURE"
					response.Message = "Failed to create following categories: " + strings.Join(categoriesFailed[:], ", ")
				}
			}
		}

		if err == nil {
			response.Status = "SUCCESS"
			response.Message = "Successfully created question"
		}
	} else {
		response.Status = "FAILURE"
		response.Message = "Failed to create question"
	}

	categories = categories[:0]
	db.Find(&categories)

	for index, category := range categories {
		response.Categories = append(response.Categories, new(qpb.CreateQuestionResponse_Category))
		response.Categories[index].Id = int32(category.ID)
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
func (question Question) GetFromDB(ctx context.Context, in *qpb.GetQuestionRequest, db *gorm.DB) (*qpb.GetQuestionResponse, error) {
	var err error
	var user User
	response := new(qpb.GetQuestionResponse)

	if in.AuthToken == "" {
		response.Status = "FAILURE"
		response.Message = "Invalid request"
		log.Errorf("No AuthToken received")
		return response, errors.New(response.Message)
	}

	email := GetEmail(in.AuthToken)

	response.Status = "SUCCESS"
	response.Message = "Next Question retrived successully!"

	if email == "" {
		response.Status = "FAILURE"
		response.Message = "Failed to retrieve email"
		log.Errorf("failed to retrieve email from token: %v", err)
		return response, errors.New("Failed to retrieve email")
	}

	result := db.Where("email = ?", email).First(&user).RecordNotFound()
	if result == true {
		fmt.Println("No user found with email: ", email)
		fmt.Println("Creating new user with email: ", email, " ...")
		user, err = CreateUser(email, db)
		if err != nil {
			response.Status = "FAILURE"
			response.Message = "Failed to create new user. " + fmt.Sprintf("Error: %s", err)
			log.Errorf("failed to create user: %v", err)
		}
	}

	fmt.Println("Request QuestionId:", in.QuestionId)

	if in.QuestionId == 0 {
		categoryID := in.GetCategoryId()
		var category Category
		result := db.Where("id = ?", categoryID).First(&category).RecordNotFound()
		if result == false {
			var question Question
			result = db.Where("category_id = ?", categoryID).First(&question).RecordNotFound()
			if result == false {
				var answer Answer
				err = db.Model(&question).Related(&answer, "Answer").Error
				if err == nil {
					fmt.Println("Answer1: ->", answer)
					question.Answer = answer
					var userAnswer UsersAnswer
					result = db.Where("user_id = ? AND answer_id = ?", user.ID, answer.ID).First(&userAnswer).RecordNotFound()
					if result == false {
						question = getQuestionFromSubCategory(category, question, user, db)
					}
				}
			}
			if question.ID != 0 {
				fmt.Println("Question: ->", question)
				response.Status = "SUCCESS"
				response.Message = "Successfully retreived next question 1"
				response.Id = int32(question.ID)
				response.Title = question.Title
				response.Body = question.Body
				responseAnswer := new(qpb.GetQuestionResponse_Answer)
				responseAnswer.Option1 = question.Answer.Option1
				responseAnswer.Option2 = question.Answer.Option2
				responseAnswer.Option3 = question.Answer.Option3
				responseAnswer.Option4 = question.Answer.Option4
				responseAnswer.Option5 = question.Answer.Option5
				response.Answer = responseAnswer
				response.Score = user.Score
			} else {
				response.Status = "FAILURE"
				response.Message = "No more open challenge available in this category"
				err = errors.New(response.Message)
			}
		} else {
			response.Status = "FAILURE"
			response.Message = "Could not find question of this category: " + string(categoryID)
			err = errors.New(response.Message)
		}
	} else {
		result = db.First(&question, in.QuestionId).RecordNotFound()
		if result == false {
			var category Category
			result = db.First(&category, question.CategoryID).RecordNotFound()
			if result == false {
				var answer Answer
				result = db.Model(&question).Related(&answer, "Answer").RecordNotFound()
				if result == false {
					var option int32
					option, err = RegisterAnswer(answer, user, in, db)
					if err != nil {
						response.Status = "FAILURE"
						response.Message = "Failed to register answer of question: " + question.Title
					} else {
						var nextQuestion Question
						nextQuestion, err = getNextQuestion(question, category, db)
						if err != nil {
							response.Status = "FAILURE"
							response.Message = "Failed to get next question for: " + question.Title
						} else {
							response.Status = "SUCCESS"
							response.Message = "Successfully retreived next question 2"
							response.Id = int32(question.ID)
							response.Title = nextQuestion.Title
							response.Body = nextQuestion.Body
							responseAnswer := new(qpb.GetQuestionResponse_Answer)
							responseAnswer.Option1 = nextQuestion.Answer.Option1
							responseAnswer.Option2 = nextQuestion.Answer.Option2
							responseAnswer.Option3 = nextQuestion.Answer.Option3
							responseAnswer.Option4 = nextQuestion.Answer.Option4
							responseAnswer.Option5 = nextQuestion.Answer.Option5
							response.Answer = responseAnswer

							var score float32
							score, err = GetPersonalityScore(user, answer, option, db)
							if err != nil {
								response.Status = "FAILURE"
								response.Message = "Failed to get personality score for: " + user.Email
							} else {
								response.Score = score
							}
						}
					}
				} else {
					response.Status = "FAILURE"
					response.Message = "Could find answer for question: " + question.Title
					err = errors.New(response.Message)
				}
			} else {
				response.Status = "FAILURE"
				response.Message = "Could find category with ID: " + strconv.Itoa(int(question.CategoryID))
				err = errors.New(response.Message)
			}
		} else {
			response.Status = "FAILURE"
			response.Message = "Could find question with ID: " + strconv.Itoa(int(in.QuestionId))
			err = errors.New(response.Message)
		}
	}

	return response, err
}

// |
// |
// |

func getNextQuestion(question Question, category Category, db *gorm.DB) (Question, error) {
	var err error
	var nextQuestion Question
	weight := new(Weight)
	err = db.Where("question_id = ?", question.ID).Find(&weight).Error
	if err != nil {
		log.Errorf("failed to get question: %v", err)
	}

	var questions []Question
	var sortedQuestions []Question
	err = db.Model(category).Association("Questions").Find(&questions).Error
	if err != nil {
		log.Errorf("failed to get question: %v", err)
	} else {
		var weightValues []int
		for _, question := range questions {
			questionWeight := new(Weight)
			err := db.Where("question_id = ?", question.ID).Find(&questionWeight).Error
			if err != nil {
				log.Errorf("failed to get weight: %v", err)
			}
			weightValues = append(weightValues, int(questionWeight.Value))
		}
		sort.Ints(weightValues)
		for _, weightValue := range weightValues {
			if len(questions) != len(sortedQuestions) {
				for _, dbQuestion := range questions {
					if weightValue == int(weight.Value) {
						sortedQuestions = append(sortedQuestions, dbQuestion)
					}
				}
			} else {
				break
			}
		}
		for _, dbQuestion := range sortedQuestions {
			sortedDbQuestionWeight := new(Weight)
			err = db.Where("question_id = ?", dbQuestion.ID).Find(&sortedDbQuestionWeight).Error
			if err != nil {
				log.Errorf("failed to get weight: %v", err)
				continue
			}
			if question.Title == dbQuestion.Title {
				continue
			}
			if sortedDbQuestionWeight.Value >= weight.Value {
				nextQuestion = dbQuestion
				var answer Answer
				err = db.Model(&nextQuestion).Related(&answer, "Answer").Error
				if err == nil {
					nextQuestion.Answer = answer
				}
				break
			}
		}
	}
	return nextQuestion, err
}

// |
// |
// |

func createAnswer(ctx context.Context, in *qpb.CreateQuestionRequest, db *gorm.DB, question Question, user User) (Answer, error) {
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
	createAnswerCategories(ctx, in.Answer.Categories, db, answer)

	if len(categoriesFailed) > 0 {
		err = errors.New("Failed to create answer categories")
		log.Errorf("failed to create following answer categories: %v", categoriesFailed)
	} else {
		err = db.Create(&answer).Error
		if err != nil {
			log.Errorf("failed to create answer: %v", err)
		} else {
			var answerWeights [5]byte
			answerWeights, err = CreateWeight(ctx, in, db, answer)
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

func createQuestion(ctx context.Context, in *qpb.CreateQuestionRequest, db *gorm.DB, user User, question Question) (Question, error) {
	question.UserID = user.ID
	question.Title = in.Title
	question.Body = in.Body
	question.Approved = false
	err := db.Create(&question).Error

	if err != nil {
		log.Errorf("failed to create question: %v", err)
	} else {
		var questionWeight = Weight{
			QuestionID: sql.NullInt64{Int64: int64(question.ID), Valid: true},
			AnswerID:   sql.NullInt64{Int64: 0, Valid: false},
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

func createQuestionCategories(ctx context.Context, requestCategory *qpb.CreateQuestionRequest_Category, db *gorm.DB, question Question) {
	var err error
	categories = categories[:0]
	assembleQuestionCategories(ctx, requestCategory.Categories)
	var category Category
	category.ID = uint(requestCategory.Id)
	category.Name = requestCategory.Name
	category.Approved = false
	category.Parent = uint(requestCategory.Parent)
	categories = append(categories, category)
	for i := len(categories) - 1; i >= 0; i-- {
		category := categories[i]
		if category.Parent == 0 {
			var dbCategory Category
			err = db.Last(&dbCategory).Error
			if err != nil {
				log.Errorf("failed to retrieve category: %v", err)
			} else {
				category.Parent = dbCategory.ID
			}
		}
		category.Level = GetLevel(&category, db)
		if category.ID == 0 {
			err = db.Create(&category).Error
		} else {
			err = db.Find(&category).Error
		}
		if err != nil {
			log.Errorf("failed to create question category: %v", err)
			categoriesFailed = append(categoriesFailed, category.Name)
		} else {
			question.CategoryID = uint(category.ID)
			err = db.Save(&question).Error
			if err != nil {
				log.Errorf("failed to save question: %v", err)
			}
		}
	}
}

// |
// |
// |

func createAnswerCategories(ctx context.Context, requestCategories []*qpb.CreateQuestionRequest_Answer_Category, db *gorm.DB, answer Answer) {
	categories = categories[:0]
	assembleAnswerCategories(ctx, requestCategories)

	for i := len(categories) - 1; i >= 0; i-- {
		category := categories[i]
		if category.Parent == 0 {
			var dbCategory Category
			err := db.Last(&dbCategory).Error
			if err != nil {
				log.Errorf("failed to retrieve category: %v", err)
			} else {
				category.Parent = dbCategory.ID
			}
		}
		category.Level = GetLevel(&category, db)
		err := db.Create(&category).Error
		if err != nil {
			log.Errorf("failed to create answer category: %v", err)
			categoriesFailed = append(categoriesFailed, category.Name)
		} else {
			if category.Option != 0 {
				answer.Categories = make([]byte, category.Option)
				answer.Categories[category.Option-1] = byte(category.ID)
			}
		}
	}
}

// |
// |
// |

func assembleQuestionCategories(ctx context.Context, requestCategories []*qpb.CreateQuestionRequest_Category) {
	for _, requestCategory := range requestCategories {
		if len(requestCategory.Categories) > 0 {
			assembleQuestionCategories(ctx, requestCategory.Categories)
		}
		var category Category
		category.ID = uint(requestCategory.Id)
		category.Name = requestCategory.Name
		category.Approved = false
		category.Parent = uint(requestCategory.Parent)
		categories = append(categories, category)
	}
}

// |
// |
// |

func assembleAnswerCategories(ctx context.Context, requestCategories []*qpb.CreateQuestionRequest_Answer_Category) {
	for _, requestCategory := range requestCategories {
		if len(requestCategory.Categories) > 0 {
			assembleAnswerCategories(ctx, requestCategory.Categories)
		}
		if requestCategory.Id == 0 {
			var category Category
			category.Name = requestCategory.Name
			category.Approved = false
			category.Parent = uint(requestCategory.Parent)
			category.Option = requestCategory.Option
			categories = append(categories, category)
		}
	}
}

// |
// |
// |

func getQuestionFromSubCategory(category Category, question Question, user User, db *gorm.DB) Question {
	var categories []Category
	var nextQuestion Question
	fmt.Println("Question--ID:", nextQuestion.ID)
	result := db.Where("category_id = ? AND id != ?", category.ID, question.ID).First(&nextQuestion).RecordNotFound()
	if result == false {
		fmt.Println("Next Question--ID:", nextQuestion.ID)
		var answer Answer
		err := db.Model(&nextQuestion).Related(&answer, "Answer").Error
		if err == nil {
			fmt.Println("Answer11: ->", answer)
			var userAnswer UsersAnswer
			result = db.Where("user_id = ? AND answer_id = ?", user.ID, answer.ID).First(&userAnswer).RecordNotFound()
			if result == false {
				db.Where("parent = ?", category.ID).Find(&categories)
				for _, nextCategory := range categories {
					nextQuestion = getQuestionFromSubCategory(nextCategory, nextQuestion, user, db)
				}
			} else {
				nextQuestion.Answer = answer
				return nextQuestion
			}
		}
	} else {
		db.Where("parent = ?", category.ID).Find(&categories)
		for _, nextCategory := range categories {
			nextQuestion = getQuestionFromSubCategory(nextCategory, nextQuestion, user, db)
		}
	}
	return nextQuestion
}
