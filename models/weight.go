package models

import (
	"database/sql"

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
