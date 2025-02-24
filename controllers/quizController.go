package controllers

import (
	"fmt"
	"quiz-crew/config"
	"quiz-crew/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetCategories retrieves all quiz categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.QuizCategory
	result := config.DB.Find(&categories)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve quiz categories!",
		})
	}

	return c.JSON(categories)
}

// CreateCategory creates a new quiz category
func CreateCategory(c *fiber.Ctx) error {
	type CategoryInput struct {
		Nama string `json:"nama"`
	}

	input := new(CategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input!",
		})
	}

	category := models.QuizCategory{
		Nama: input.Nama,
	}

	if result := config.DB.Create(&category); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create quiz category!",
		})
	}

	return c.Status(201).JSON(category)
}

// CreateQuestion creates a new question for a category
func CreateQuestion(c *fiber.Ctx) error {
	type AnswerInput struct {
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}

	type QuestionInput struct {
		QuizCategoryID uint          `json:"quiz_category_id"`
		Text           string        `json:"text"`
		ImageUrl       string        `json:"image_url"`
		Answers        []AnswerInput `json:"answers"`
	}

	input := new(QuestionInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input!",
		})
	}

	// Validate category exists
	var category models.QuizCategory
	if result := config.DB.First(&category, input.QuizCategoryID); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Quiz category not found!",
		})
	}

	// Create question
	question := models.Question{
		QuizCategoryID: input.QuizCategoryID,
		Text:           input.Text,
		ImageUrl:       input.ImageUrl,
	}
	tx := config.DB.Begin()

	if result := tx.Create(&question); result.Error != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create question!",
		})
	}

	// Create answers
	for _, answerInput := range input.Answers {
		answer := models.Answer{
			QuestionID: question.ID,
			Text:       answerInput.Text,
			IsCorrect:  answerInput.IsCorrect,
		}

		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create answer!",
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save changes!",
		})
	}

	var completeQuestion models.Question
	config.DB.Preload("Answers").First(&completeQuestion, question.ID)

	return c.Status(201).JSON(completeQuestion)
}

// GetQuestions retrieves shuffled questions for a category
func GetQuestions(c *fiber.Ctx) error {
	var questions []models.Question
	categoryID := c.Params("id")

	if err := config.DB.Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC") // Order answers by ID
	}).Where("quiz_category_id = ?", categoryID).
		Order("id ASC"). // Order questions by ID
		Find(&questions).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Quiz category not found!",
		})
	}

	fmt.Println("\nQuestion order in response:")
	for i, q := range questions {
		fmt.Printf("Position %d: Question ID %d\n", i+1, q.ID)
	}

	return c.JSON(fiber.Map{
		"questions": questions,
	})
}

func SubmitAnswers(c *fiber.Ctx) error {
	type Submission struct {
		Nik        string `json:"nik"`
		CategoryID uint   `json:"category_id"`
		Answers    []int  `json:"answers"`
		TimeTaken  int    `json:"time_taken"`
	}

	submission := new(Submission)
	if err := c.BodyParser(submission); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input!",
		})
	}

	// Find user by NIK first
	var user models.User
	if err := config.DB.Where("nik = ?", submission.Nik).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found!",
		})
	}

	var questions []models.Question
	if err := config.DB.Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC") // Order answers by ID
	}).Where("quiz_category_id = ?", submission.CategoryID).
		Order("id ASC"). // Order questions by ID
		Find(&questions).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Quiz category not found!",
		})
	}

	// Debug print question order in submission
	fmt.Println("\nQuestion order in submission:")
	for i, q := range questions {
		fmt.Printf("Position %d: Question ID %d\n", i+1, q.ID)
	}

	fmt.Printf("\nSubmitted answers array: %v\n", submission.Answers)

	// Calculate score
	score := 0
	fmt.Printf("\nStarting score calculation:\n")
	fmt.Printf("Total questions: %d\n", len(questions))
	fmt.Printf("Total answers submitted: %d\n", len(submission.Answers))
	for i, question := range questions {
		if i >= len(submission.Answers) {
			fmt.Printf("No answer submitted for question %d\n", i+1)
			continue
		}

		submittedAnswerID := submission.Answers[i]

		// Debug print for each question's answers
		fmt.Printf("\nQuestion %d:\n", i+1)
		fmt.Printf("Submitted answer ID: %d\n", submittedAnswerID)
		fmt.Printf("Available answers:\n")

		// Print all answers for debugging
		for _, ans := range question.Answers {
			fmt.Printf("- Answer ID: %d, Text: %s, IsCorrect: %v\n",
				ans.ID, ans.Text, ans.IsCorrect)
		}

		foundCorrect := false
		correctAnswerID := 0
		for _, ans := range question.Answers {
			if ans.IsCorrect {
				correctAnswerID = int(ans.ID)
			}
			if ans.IsCorrect && int(ans.ID) == submittedAnswerID {
				score++
				foundCorrect = true
				fmt.Printf("✓ Correct answer found! Score increased to %d\n", score)
				break
			}
		}
		if !foundCorrect {
			fmt.Printf("✗ Wrong answer for question %d. Expected: %d, Got: %d\n", i+1, correctAnswerID, submittedAnswerID)
		}
	}

	totalQuestions := len(questions)
	passingScore := int(0.8 * float64(totalQuestions))
	passed := score >= passingScore

	fmt.Printf("Score: %d/%d (Passing: %d)\n", score, totalQuestions, passingScore)

	attempt := models.QuizAttempt{
		UserID:         user.ID,
		QuizCategoryID: submission.CategoryID,
		Score:          score,
		PassingStatus:  passed,
	}

	if err := config.DB.Create(&attempt).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save quiz attempt!",
		})
	}

	return c.JSON(fiber.Map{
		"user_id":       user.ID,
		"category_id":   submission.CategoryID,
		"score":         score,
		"total":         totalQuestions,
		"passing_score": passingScore,
		"passed":        passed,
		"time_taken":    submission.TimeTaken,
	})
}
