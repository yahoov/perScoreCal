package server

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"perScoreCal/models"

	qpb "perScoreCal/perScoreProto/question"
	upb "perScoreCal/perScoreProto/user"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// UserServer implements the UserServer interface
type UserServer struct {
	user models.User
}

// QuestionServer implements the QuestionServer interface
type QuestionServer struct {
	question models.Question
}

func (s *UserServer) GetInterests(ctx context.Context, in *upb.GetInterestRequest) (*upb.GetInterestResponse, error) {
	fmt.Println("Request: ", in)
	var result *upb.GetInterestResponse
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DEV_HOST"), os.Getenv("DEV_DBNAME"), os.Getenv("DEV_USERNAME"), os.Getenv("DEV_PASSWORD"), os.Getenv("DEV_SSLMODE"))
	db, err := gorm.Open(os.Getenv("DEV_DB_DRIVER"), dbString)
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.GetInterests(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetInterests: %+v", err)
		}
	}
	fmt.Println("Result: ", in)
	return result, nil
}

func (s *UserServer) GetEntries(ctx context.Context, in *upb.GetEntriesRequest) (*upb.GetEntriesResponse, error) {
	fmt.Println("Request: ", in)
	var result *upb.GetEntriesResponse
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DEV_HOST"), os.Getenv("DEV_DBNAME"), os.Getenv("DEV_USERNAME"), os.Getenv("DEV_PASSWORD"), os.Getenv("DEV_SSLMODE"))
	db, err := gorm.Open(os.Getenv("DEV_DB_DRIVER"), dbString)
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.GetEntries(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetEntries: %+v", err)
		}
	}
	fmt.Println("Result: ", in)
	return result, nil
}

func (s *UserServer) ApproveEntries(ctx context.Context, in *upb.ApproveEntriesRequest) (*upb.ApproveEntriesResponse, error) {
	fmt.Println("Request: ", in)
	var result *upb.ApproveEntriesResponse
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DEV_HOST"), os.Getenv("DEV_DBNAME"), os.Getenv("DEV_USERNAME"), os.Getenv("DEV_PASSWORD"), os.Getenv("DEV_SSLMODE"))
	db, err := gorm.Open(os.Getenv("DEV_DB_DRIVER"), dbString)
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.ApproveEntries(ctx, in, db)
		if err != nil {
			log.Errorf("Error in ApproveEntries: %+v", err)
		}
	}
	fmt.Println("Result: ", in)
	return result, nil
}

// CreateQuestion creates a new question
func (s *QuestionServer) CreateQuestion(ctx context.Context, in *qpb.CreateQuestionRequest) (*qpb.CreateQuestionResponse, error) {
	fmt.Println("Request: ", in)
	var result *qpb.CreateQuestionResponse
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DEV_HOST"), os.Getenv("DEV_DBNAME"), os.Getenv("DEV_USERNAME"), os.Getenv("DEV_PASSWORD"), os.Getenv("DEV_SSLMODE"))
	db, err := gorm.Open(os.Getenv("DEV_DB_DRIVER"), dbString)
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.question.CreateInDB(ctx, in, db)
		if err != nil {
			log.Errorf("Error in CreateInDB: %+v", err)
		}
	}
	fmt.Println("Result: ", in)
	return result, nil
}

// GetQuestion fetches a new question for the given answer
func (s *QuestionServer) GetQuestion(ctx context.Context, in *qpb.GetQuestionRequest) (*qpb.GetQuestionResponse, error) {
	fmt.Println("Request: ", in)
	var result *qpb.GetQuestionResponse
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DEV_HOST"), os.Getenv("DEV_DBNAME"), os.Getenv("DEV_USERNAME"), os.Getenv("DEV_PASSWORD"), os.Getenv("DEV_SSLMODE"))
	db, err := gorm.Open(os.Getenv("DEV_DB_DRIVER"), dbString)
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.question.GetFromDB(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetFromDB: %+v", err)
		}
	}
	fmt.Println("Result: ", in)
	return result, nil
}
