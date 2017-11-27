package models

import (
	"sort"
	"strconv"

	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// Category is a gorm model
type Category struct {
	gorm.Model
	Questions []Question
	Name      string
	Parent    uint
	Level     int32
	Approved  bool
}

// GetLevel ...
func GetLevel(category *Category, db *gorm.DB) int32 {
	var parentCategory Category
	db.First(&parentCategory, category.Parent)
	return parentCategory.Level + 1
}

// GetWeightRange ...
func GetWeightRange(category *Category, db *gorm.DB) string {
	var questions []Question
	err := db.Model(category).Association("Questions").Find(&questions).Error
	// err := db.Preload("Category", "ID = (?)", category.ID).Find(&questions).Order("Weight asc").Error
	// db.Model(&category).Related(&questions).Order("Weight asc")
	if err != nil {
		log.Errorf("failed to retrieve associated questions: %v", err)
		return ""
	}
	var weightValues []int
	for _, question := range questions {
		weightValues = append(weightValues, int(question.Weight.Value))
	}
	sort.Ints(weightValues)
	min := strconv.Itoa(weightValues[0])
	max := strconv.Itoa(weightValues[len(weightValues)-1])
	return min + ".." + max
}