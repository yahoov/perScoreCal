package models

import (
	"context"
	"database/sql"
	"fmt"

	qpb "perScoreCal/perScoreProto/question"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
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

// CreateWeight ...
func CreateWeight(ctx context.Context, in *qpb.CreateQuestionRequest, db *gorm.DB, answer Answer) ([5]byte, error) {
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
