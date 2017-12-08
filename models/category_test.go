package models_test

import (
	"perScoreCal/models"
	"testing"
)

func TestWeightRange(t *testing.T) {

	var questions = []models.Question{
		models.Question{
			Title:    "daya",
			Body:     "12",
			Approved: true,
			Weight: models.Weight{
				Value: 1,
			},
		},
		models.Question{
			Title:    "daya",
			Body:     "12",
			Approved: true,
			Weight: models.Weight{
				Value: 2,
			},
		},
		models.Question{
			Title:    "daya",
			Body:     "12",
			Approved: true,
			Weight: models.Weight{
				Value: 3,
			},
		},
	}

	// weighRange := models.WeightRange(questions)
	//
	// fmt.Println("weighRange", weighRange)
	if weighRange := models.WeightRange(questions); weighRange != "1..3" {
		t.Errorf("Expected email to be %s, but it was %s", "1..3", weighRange)
	}
	// dbconn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("TEST_HOST"), os.Getenv("TEST_DBNAME"), os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_SSLMODE"))
	// db, err := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbconn)
	// if err != nil {
	// 	t.Errorf("Error in testdbsetup: %+v", err)
	// }
	// weighRange := models.GetWeightRange(&test_Category, db)
	// fmt.Println("weithrange data has printed:::", weighRange)
}
