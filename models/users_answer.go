package models

import (
	"github.com/jinzhu/gorm"

	pb "perScoreCal/perScoreProto/question"
)

// UsersAnswer is a gorm model
type UsersAnswer struct {
	gorm.Model
	UserID   uint
	AnswerID uint
	Option   int32
}

// RegisterAnswer ...
func RegisterAnswer(answer Answer, user User, in *pb.GetQuestionRequest, db *gorm.DB) (int32, error) {
	usersAnswer := UsersAnswer{
		UserID:   user.ID,
		AnswerID: answer.ID,
		Option:   getOption(in.Answer),
	}

	err := db.Create(&usersAnswer).Error

	return usersAnswer.Option, err
}

func getOption(answer *pb.GetQuestionRequest_Answer) int32 {
	var result int32
	if answer.Option1 == true {
		result = 0
	} else if answer.Option2 == true {
		result = 1
	} else if answer.Option3 == true {
		result = 2
	} else if answer.Option4 == true {
		result = 3
	} else if answer.Option5 == true {
		result = 4
	}

	return result
}
