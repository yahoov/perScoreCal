package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect for gorm
)

// SetupDatabase - Creates the tables in the database
func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Question{},
		&Answer{},
		&UsersAnswer{},
		&Weight{},
		&Category{},
	).Error

	return err
}
