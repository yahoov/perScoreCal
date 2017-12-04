package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

// Weight is a gorm model
type Weight struct {
	gorm.Model
	QuestionID sql.NullInt64
	AnswerID   sql.NullInt64
	Value      int32
	Option     int32
}

// |
// |
// |

func CreateWeight(ctx context.Context, in *pb.CreateQuestionRequest, db *gorm.DB, answer Answer) ([5]byte, error) {
	var answerWeights [5]byte
	var err error
	for _, weight := range in.Answer.Weights {
		var createdWeight Weight
		createdWeight.QuestionID = sql.NullInt64{Int64: 0, Valid: false}
		createdWeight.AnswerID = sql.NullInt64{Int64: int64(answer.ID), Valid: true}
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
