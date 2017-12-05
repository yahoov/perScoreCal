package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// CreateInitialCategories ...
func CreateInitialCategories(db *gorm.DB) error {
	var err error
	categories := [6]string{"Science", "Maths", "History", "Entertainment", "Sports", "Spirituality"}

<<<<<<< HEAD
	for _, categoryTile := range categories {
		category := Category{
			Name:     categoryTile,
=======
	for _, categoryTitle := range categories {
		category := Category{
			Name:     categoryTitle,
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2
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
