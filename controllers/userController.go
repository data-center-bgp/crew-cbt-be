package controllers

import (
	"quiz-crew/config"
	"quiz-crew/models"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input!",
		})
	}

	matched, _ := regexp.MatchString(`^\d{16}$`, user.Nik)
	if !matched {
		return c.Status(400).JSON(fiber.Map{
			"error": "NIK must be 16 digits!",
		})
	}

	config.DB.Create(&user)
	return c.JSON(user)
}

func GetUserResult(c *fiber.Ctx) error {
	nik := c.Params("nik")
	categoryID := c.QueryInt("category_id")

	var user models.User
	if err := config.DB.Where("nik = ?", nik).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found!",
		})
	}

	var attempt models.QuizAttempt
	result := config.DB.Where("user_id = ? AND quiz_category_id = ?", user.ID, categoryID).Order("created_at DESC").First(&attempt)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Quiz attempt not found",
		})
	}

	return c.JSON(fiber.Map{
		"nama":           user.Nama,
		"nik":            user.Nik,
		"jabatan":        user.Jabatan,
		"perusahaan":     user.Perusahaan,
		"category_id":    attempt.QuizCategoryID,
		"score":          attempt.Score,
		"passing_status": attempt.PassingStatus,
		"attempted_at":   attempt.CreatedAt,
	})

}
