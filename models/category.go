package models

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect for gorm
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
	Option    int32 `gorm:"-"` // ignore this field
}

// GetLevel ...
func GetLevel(category *Category, db *gorm.DB) int32 {
	var parentCategory Category
	db.First(&parentCategory, category.Parent)
	return parentCategory.Level + 1
}

// GetWeightRange ...
func GetWeightRange(category *Category, db *gorm.DB) string {
	fmt.Println("category", category.Name)
	fmt.Println("db", db)
	var questions []Question
	err := db.Model(category).Association("Questions").Find(&questions).Error
	if err != nil {
		log.Errorf("failed to retrieve associated questions: %v", err)
		return ""
	}
	if len(questions) == 0 {
		fmt.Println("No associated questions with categoryID:", category.ID)
		return ""
	}
	return WeightRange(questions)
}

//WeightRange ...
func WeightRange(questions []Question) string {
	var weightValues []int
	for _, question := range questions {
		weightValues = append(weightValues, int(question.Weight.Value))
	}
	sort.Ints(weightValues)
	min := strconv.Itoa(weightValues[0])
	max := strconv.Itoa(weightValues[len(weightValues)-1])
	return min + " .. " + max
}
