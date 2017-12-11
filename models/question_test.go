package models_test

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "perScoreCal/models"
	qpb "perScoreCal/perScoreProto/question"

	"github.com/jinzhu/gorm"
)

var _ = Describe("Question", func() {
	var (
		// user     User
		question Question
		in       *qpb.CreateQuestionRequest
		ctx      context.Context
	)

	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("TEST_HOST"), os.Getenv("TEST_DBNAME"), os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_SSLMODE"))

	Describe("CreateInDB", func() {
		db, err := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbString)
		defer db.Close()

		if err != nil {
			fmt.Println("Error opening DB connection:", err)
		} else {

			Context("with correct input", func() {
				BeforeEach(func() {
					in = createRequestData(in)
				})
				It("should create the question in DB", func() {
					result, _ := question.CreateInDB(ctx, in, db)
					Expect(result.Status).To(Equal("SUCCESS"))
				})
			})

			Context("with incorrect input", func() {
				It("should not create the question in DB", func() {
					result, _ := question.CreateInDB(ctx, in, db)
					Expect(result.Status).To(Equal("FAILURE"))
				})
			})

		}
	})

	Describe("GetFromDB", func() {
		var (
			question Question
			in       *qpb.GetQuestionRequest
			ctx      context.Context
		)
		db, err := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbString)
		defer db.Close()
		if err != nil {
			fmt.Println("Error opening DB connection:", err)
		} else {
			Context("with correct input", func() {
				BeforeEach(func() {
					in = getRequestData(in)
				})
				It("result status should be success", func() {
					result, _ := question.GetFromDB(ctx, in, db)
					Expect(result.Status).To(Equal("SUCCESS"))
				})
			})

			Context("with incorrect input", func() {
				It("result status should be failed", func() {
					result, _ := question.GetFromDB(ctx, in, db)
					Expect(result.Status).To(Equal("FAILURE"))
				})
			})
		}
	})
})

var _ = BeforeSuite(func() {
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("TEST_HOST"), os.Getenv("TEST_DBNAME"), os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_SSLMODE"))
	db, _ := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbString)
	CreateInitialCategories(db)
})

var _ = AfterSuite(func() {
	dbString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("TEST_HOST"), os.Getenv("TEST_DBNAME"), os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_SSLMODE"))
	db, _ := gorm.Open(os.Getenv("TEST_DB_DRIVER"), dbString)
	db.Exec("TRUNCATE answers, categories, questions, users, users_answers, weights CASCADE;")
	db.Close()
})

func createRequestData(in *qpb.CreateQuestionRequest) *qpb.CreateQuestionRequest {
	in.AuthToken = Encrypt("test@email.com")
	in.Title = "example title"
	in.Body = "example body"
	in.Answer.Option1 = "option1"
	in.Answer.Option2 = "option2"
	in.Answer.Option3 = "option3"
	in.Answer.Option4 = "option4"
	in.Answer.Option5 = "option5"
	in.Answer.Weights = []*qpb.CreateQuestionRequest_Answer_Weight{
		&qpb.CreateQuestionRequest_Answer_Weight{
			Value:  1,
			Option: 1,
		},
		&qpb.CreateQuestionRequest_Answer_Weight{
			Value:  2,
			Option: 2,
		},
		&qpb.CreateQuestionRequest_Answer_Weight{
			Value:  3,
			Option: 3,
		},
		&qpb.CreateQuestionRequest_Answer_Weight{
			Value:  4,
			Option: 4,
		},
		&qpb.CreateQuestionRequest_Answer_Weight{
			Value:  5,
			Option: 5,
		},
	}
	in.Answer.Categories = []*qpb.CreateQuestionRequest_Answer_Category{
		&qpb.CreateQuestionRequest_Answer_Category{
			Name:   "answer_cat_1",
			Option: 1,
			Parent: 1,
		},
		&qpb.CreateQuestionRequest_Answer_Category{
			Name:   "answer_cat_2",
			Option: 2,
			Parent: 2,
		},
		&qpb.CreateQuestionRequest_Answer_Category{
			Name:   "answer_cat_3",
			Option: 3,
			Parent: 3,
		},
		&qpb.CreateQuestionRequest_Answer_Category{
			Name:   "answer_cat_4",
			Option: 4,
			Parent: 1,
		},
		&qpb.CreateQuestionRequest_Answer_Category{
			Name:   "answer_cat_5",
			Option: 5,
			Parent: 2,
		},
	}
	in.Weight = &qpb.CreateQuestionRequest_Weight{
		Value: 1,
	}
	in.Categories = []*qpb.CreateQuestionRequest_Category{
		&qpb.CreateQuestionRequest_Category{
			Name:   "question_cat_1",
			Parent: 1,
		},
	}
	return in
}
func getRequestData(in *qpb.GetQuestionRequest) *qpb.GetQuestionRequest {
	in.AuthToken = Encrypt("test@email.com")
	in.GetQuestionId()
	in.Answer.Option1 = true
	in.Answer.Option2 = true
	in.Answer.Option3 = true
	in.Answer.Option4 = true
	in.Answer.Option5 = true
	return in
}
