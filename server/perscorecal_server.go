package server

import (
	"fmt"

	"golang.org/x/net/context"

	"perScoreCal/models"
	pb "perScoreCal/perScoreProto/question"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Server implements the QuestionServer interface
type Server struct {
	question models.Question
}

// CreateQuestion creates a new question
func (s *Server) CreateQuestion(ctx context.Context, in *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	fmt.Println("request: ", in)
	var result *pb.CreateQuestionResponse
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
func (s *Server) GetQuestion(ctx context.Context, in *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	var result *pb.GetQuestionResponse
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
