package models

import "github.com/jinzhu/gorm"

// Weight is a gorm model
type Weight struct {
	gorm.Model
	QuestionID uint
	AnswerID   uint
	Value      int32
	Option     int32
}
