package main

import (
	"log"
	"quiz-crew/config"
	"quiz-crew/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()
	defer config.CloseDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://crew-cbt.barokahperkasagroup.id, https://images.barokahperkasagroup.id",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origins, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))

	routes.SetupRoutes(app)

	log.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
