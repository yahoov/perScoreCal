package server

import (
	"fmt"

	"golang.org/x/net/context"

	"perScoreCal/models"
<<<<<<< HEAD

	qpb "perScoreCal/perScoreProto/question"
	upb "perScoreCal/perScoreProto/user"
=======
	pb "perScoreCal/perScoreProto/question"
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

<<<<<<< HEAD
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
=======
// Server implements the QuestionServer interface
type Server struct {
	question models.Question
}

// CreateQuestion creates a new question
func (s *Server) CreateQuestion(ctx context.Context, in *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	fmt.Println("request: ", in)
	var result *pb.CreateQuestionResponse
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2
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
<<<<<<< HEAD
func (s *QuestionServer) GetQuestion(ctx context.Context, in *qpb.GetQuestionRequest) (*qpb.GetQuestionResponse, error) {
	fmt.Println("request: ", in)
	var result *qpb.GetQuestionResponse
=======
func (s *Server) GetQuestion(ctx context.Context, in *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	var result *pb.GetQuestionResponse
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2
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
