package server

import (
	"fmt"

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
	fmt.Println("request: ", in)
	var result *upb.GetInterestResponse
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.GetInterests(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetInterests: %+v", err)
		}
	}
	return result, err
}

func (s *UserServer) GetEntries(ctx context.Context, in *upb.GetEntriesRequest) (*upb.GetEntriesResponse, error) {
	fmt.Println("request: ", in)
	var result *upb.GetEntriesResponse
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.GetEntries(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetEntries: %+v", err)
		}
	}
	return result, err
}

func (s *UserServer) ApproveEntries(ctx context.Context, in *upb.ApproveEntriesRequest) (*upb.ApproveEntriesResponse, error) {
	fmt.Println("request: ", in)
	var result *upb.ApproveEntriesResponse
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.user.ApproveEntries(ctx, in, db)
		if err != nil {
			log.Errorf("Error in ApproveEntries: %+v", err)
		}
	}
	return result, err
}

// CreateQuestion creates a new question
func (s *QuestionServer) CreateQuestion(ctx context.Context, in *qpb.CreateQuestionRequest) (*qpb.CreateQuestionResponse, error) {
	fmt.Println("request: ", in)
	var result *qpb.CreateQuestionResponse
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.question.CreateInDB(ctx, in, db)
		if err != nil {
			log.Errorf("Error in CreateInDB: %+v", err)
		}
	}
	return result, err
}

// GetQuestion fetches a new question for the given answer
func (s *QuestionServer) GetQuestion(ctx context.Context, in *qpb.GetQuestionRequest) (*qpb.GetQuestionResponse, error) {
	fmt.Println("request: ", in)
	var result *qpb.GetQuestionResponse
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=perscorecal dbname=per_score_cal sslmode=disable password=perscorecal-dm")
	defer db.Close()
	if err != nil {
		log.Errorf("Error opening DB connection: %+v", err)
	} else {
		result, err = s.question.GetFromDB(ctx, in, db)
		if err != nil {
			log.Errorf("Error in GetFromDB: %+v", err)
		}
	}
	return result, err
}
