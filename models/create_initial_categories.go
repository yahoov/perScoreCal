package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// CreateInitialCategories ...
func CreateInitialCategories(db *gorm.DB) error {
	var err error
	categories := [6]string{"Science", "Maths", "History", "Entertainment", "Sports", "Spirituality"}

	for _, categoryTitle := range categories {
		category := Category{
			Name:     categoryTitle,
			Parent:   0,
			Level:    1,
			Approved: true,
		}
		err = db.Create(&category).Error
		if err != nil {
			log.Errorf("failed to create initial category: %v", err)
		}
	}

	return err
}
