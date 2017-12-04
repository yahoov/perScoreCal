package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect for gorm
)

// Answer is a gorm model
type Answer struct {
	gorm.Model
	UserID      uint
	QuestionID  uint
	Weights     []byte `gorm:"type=bytea"`
	Categories  []byte `gorm:"type=bytea"`
	Option1     string
	Option2     string
	Option3     string
	Option4     string
	Option5     string
	UsersAnswer UsersAnswer
}
