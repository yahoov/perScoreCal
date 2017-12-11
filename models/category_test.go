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
	if weighRange := models.WeightRange(questions); weighRange != "1..3" {
		t.Errorf("Expected email to be %s, but it was %s", "1..3", weighRange)
	}
}
