package main

import (
	"log"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/fakhri-rasyad/wpu_goreact/controller"
	"github.com/fakhri-rasyad/wpu_goreact/database/seed"
	"github.com/fakhri-rasyad/wpu_goreact/repositories"
	"github.com/fakhri-rasyad/wpu_goreact/routes"
	"github.com/fakhri-rasyad/wpu_goreact/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectToDB()

	seed.SeedAdmin()

	app := fiber.New()
	userRepo := repositories.NewUserRepostiry()
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	routes.Setup(app, userController)

	port := config.APPConfig.APPPort

	log.Println("Server is running on port: ", port)
	log.Fatal(app.Listen(":" + port))

}