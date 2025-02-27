package routes

import (
	"quiz-crew/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/crew_cbt")

	// User registration
	users := api.Group("/user")
	users.Post("/register", controllers.RegisterUser)
	users.Get("/:nik", controllers.GetUserResult)

	// Quiz routes
	quiz := api.Group("/quiz")
	quiz.Get("/categories", controllers.GetCategories)
	quiz.Get("/categories/:id", controllers.GetQuestions)
	quiz.Post("/submit", controllers.SubmitAnswers)

	// Create quiz and its categories
	quiz.Post("/createCategory", controllers.CreateCategory)
	quiz.Post("/createQuestion", controllers.CreateQuestion)
}
