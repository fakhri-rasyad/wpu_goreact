package routes

import (
	"log"

	"github.com/fakhri-rasyad/wpu_goreact/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controller.UserController){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Post("v1/auth/register", uc.Register)
}