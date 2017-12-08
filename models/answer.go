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
	Option1     string `json:"Option1"`
	Option2     string `json:"Option2"`
	Option3     string `json:"Option3"`
	Option4     string `json:"Option4"`
	Option5     string `json:"Option5"`
	UsersAnswer UsersAnswer
}
